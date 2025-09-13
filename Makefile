# Makefile for Employee Management System

.PHONY: help build run test test-contract test-integration test-unit test-all lint clean docker-up docker-down migrate-up migrate-down

# Default target
help:
	@echo "Available commands:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run all tests"
	@echo "  test-contract - Run contract tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-unit  - Run unit tests"
	@echo "  test-all   - Run all tests with coverage"
	@echo "  lint       - Run linter"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker-up  - Start Docker services"
	@echo "  docker-down - Stop Docker services"
	@echo "  migrate-up - Run database migrations"
	@echo "  migrate-down - Rollback database migrations"

# Tidy & Vendor
tidy:
	@echo "Tidy & Vendor"
	go mod tidy
	go mod vendor

# Build targets
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

# Development targets
run:
	@echo "Starting server..."
	go run cmd/server/main.go

# Test targets
test:
	@echo "Running all tests..."
	go test -v ./...

test-contract:
	@echo "Running contract tests..."
	go test -v ./tests/contract/...

test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

test-all:
	@echo "Running all tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality targets
lint:
	@echo "Running linter..."
	golangci-lint run

fmt:
	@echo "Formatting code..."
	go fmt ./...

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker targets
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Database targets
migrate-up:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	go run cmd/migrate/main.go down

migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)..."
	migrate create -ext sql -dir migrations -seq $(name)

# Development workflow
dev-setup: docker-up migrate-up
	@echo "Development environment setup complete!"

dev-test: test-contract
	@echo "Contract tests completed!"

dev-build: fmt lint test-all build
	@echo "Build completed successfully!"