// Package config memuat variabel lingkungan dari .env atau sistem.
package config

import (
	"os"
)

// Config berisi pengaturan runtime backend.
type Config struct {
	Environment string // development | production
	Port        string // mis. 8080
	DatabaseURL string // DSN PostgreSQL (lihat .env.example)
	JWTSecret   string // penandatanganan access/refresh token
	CORSOrigin  string // origin frontend yang diizinkan
}

// Load membaca env; nilai default hanya untuk development lokal.
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://audi:090393@localhost:5432/learn?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
