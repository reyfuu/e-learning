package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"learn-platform/backend/internal/config"
	"learn-platform/backend/internal/handler"
	"learn-platform/backend/internal/repository"
	"learn-platform/backend/internal/service"
	"learn-platform/backend/pkg/database"
)

func init() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🚀 Starting LevelUp.dev Backend...")

	// Load configuration
	cfg := config.Load()
	log.Printf("📋 Config loaded: env=%s, port=%s\n", cfg.Environment, cfg.Port)

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	courseRepo := repository.NewCourseRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, userRepo, cfg.JWTSecret)
	courseService := service.NewCourseService(courseRepo)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// Security middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
	}))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	courseHandler := handler.NewCourseHandler(courseService)

	// API routes
	api := e.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.RefreshToken)

	// User routes (protected)
	users := api.Group("/users")
	users.Use(authHandler.Middleware)
	users.GET("", userHandler.GetAll)
	users.GET("/me", userHandler.GetCurrentUser) // must be before /:id
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	// Course routes (protected)
	courses := api.Group("/courses")
	courses.Use(authHandler.Middleware)
	courses.GET("", courseHandler.GetUserCourses)
	courses.GET("/catalog", courseHandler.GetAllCourses)
	courses.GET("/enrolled", courseHandler.GetEnrolledCourses)
	courses.POST("", courseHandler.CreateCourse)
	courses.GET("/:id", courseHandler.GetCourse)
	courses.POST("/:id/enroll", courseHandler.EnrollCourse)
	courses.PUT("/:id/status", courseHandler.UpdateCourseStatus)
	courses.POST("/:id/modules", courseHandler.CreateModule)
	courses.PUT("/modules/:id/toggle", courseHandler.ToggleModuleCompletion)

	// Start server
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("🌐 Server started on http://localhost%s\n", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏹️ Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown error: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}
