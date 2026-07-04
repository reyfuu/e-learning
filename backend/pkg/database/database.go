// Package database mengelola koneksi PostgreSQL dan migrasi skema GORM.
package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"learn-platform/backend/internal/model"
)

// New membuka koneksi ke PostgreSQL, memverifikasi dengan Ping,
// lalu menjalankan AutoMigrate untuk model yang terdaftar (saat ini: User).
func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database handle: %w", err)
	}
	// Pastikan database benar-benar dapat dijangkau sebelum migrasi
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Category{}, &model.Course{}, &model.Section{}, &model.Lesson{}, &model.Enrollment{}, &model.LessonProgress{}, &model.Transaction{}, &model.Review{}, &model.Certificate{}); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ PostgreSQL connected and migrations completed")
	return db, nil
}
