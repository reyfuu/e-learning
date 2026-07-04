.PHONY: help build run test clean migrate docker-build docker-up docker-down

help:
	@echo "Learn Platform Backend - Available Commands"
	@echo "==========================================="
	@echo "make build        - Build the application"
	@echo "make run          - Run the application"
	@echo "make dev          - Run with hot reload (requires air)"
	@echo "make test         - Run tests"
	@echo "make test-cover   - Run tests with coverage"
	@echo "make lint         - Run golangci-lint"
	@echo "make fmt          - Format code"
	@echo "make clean        - Clean build artifacts"
	@echo "make docker-build - Build Docker image"
	@echo "make docker-up    - Start Docker containers"
	@echo "make docker-down  - Stop Docker containers"

build:
	@echo "Building application..."
	go build -o bin/learn-platform ./

run: build
	@echo "Starting application..."
	./bin/learn-platform

dev:
	@echo "Starting with hot reload..."
	air -c .air.toml

test:
	@echo "Running tests..."
	go test -v -race -timeout 30s ./...

test-cover:
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

docker-build:
	@echo "Building Docker image..."
	docker build -t learn-platform-backend:latest .

docker-up:
	@echo "Starting Docker containers (from apps/)..."
	cd .. && docker compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	cd .. && docker compose down

docker-logs:
	@echo "Showing Docker logs..."
	cd .. && docker compose logs -f backend
