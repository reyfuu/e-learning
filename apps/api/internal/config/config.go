// Package config memuat variabel lingkungan dari .env atau sistem.
package config

import (
	"os"
)

// Config berisi pengaturan runtime backend.
type Config struct {
	Environment     string // development | production
	Port            string // mis. 8080
	Host            string // alamat bind; loopback untuk penggunaan lokal
	DatabaseURL     string // DSN PostgreSQL (lihat .env.example)
	JWTSecret       string // penandatanganan access/refresh token
	CORSOrigin      string // origin frontend yang diizinkan
	CodeRunnerURL   string // URL instance Judge0; kosong berarti eksekusi dinonaktifkan
	CodeRunnerToken string
}

// Load membaca env; nilai default hanya untuk development lokal.
func Load() *Config {
	return &Config{
		Environment:     getEnv("ENVIRONMENT", "development"),
		Port:            getEnv("PORT", "8080"),
		Host:            getEnv("HOST", "127.0.0.1"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://audi:090393@localhost:5432/learn?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		CORSOrigin:      getEnv("CORS_ORIGIN", "http://localhost:3000"),
		CodeRunnerURL:   os.Getenv("CODE_RUNNER_URL"),
		CodeRunnerToken: os.Getenv("CODE_RUNNER_TOKEN"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
