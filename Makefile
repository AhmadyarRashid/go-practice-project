# Go Enterprise API Makefile

# Variables
APP_NAME=go-enterprise-api
MAIN_PATH=./cmd/api
BUILD_DIR=./build
GO=go
GOFLAGS=-ldflags="-s -w"

# Colors for terminal output
GREEN=\033[0;32m
NC=\033[0m # No Color

.PHONY: all build run test clean deps lint fmt vet coverage help docker-build docker-run migrate

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## all: Run all checks and build
all: deps fmt vet lint test build

## build: Build the application
build:
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

## run: Run the application
run:
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	$(GO) run $(MAIN_PATH)/main.go

## dev: Run the application with hot reload (requires air)
dev:
	@echo "$(GREEN)Running $(APP_NAME) in development mode...$(NC)"
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

## test: Run all tests
test:
	@echo "$(GREEN)Running tests...$(NC)"
	$(GO) test -v -race ./...

## test-short: Run short tests only
test-short:
	@echo "$(GREEN)Running short tests...$(NC)"
	$(GO) test -v -short ./...

## coverage: Run tests with coverage
coverage:
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f *.db
	@echo "$(GREEN)Clean complete$(NC)"

## deps: Download and tidy dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)Dependencies ready$(NC)"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## fmt: Format code
fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

## vet: Run go vet
vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	$(GO) vet ./...

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(APP_NAME):latest .

## docker-run: Run Docker container
docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

## docker-compose-up: Start all services with docker-compose
docker-compose-up:
	@echo "$(GREEN)Starting services...$(NC)"
	docker-compose up -d

## docker-compose-down: Stop all services
docker-compose-down:
	@echo "$(GREEN)Stopping services...$(NC)"
	docker-compose down

## swagger: Generate Swagger documentation (requires swag)
swagger:
	@echo "$(GREEN)Generating Swagger documentation...$(NC)"
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g $(MAIN_PATH)/main.go -o ./docs

## migrate: Run database migrations
migrate:
	@echo "$(GREEN)Running migrations...$(NC)"
	$(GO) run $(MAIN_PATH)/main.go migrate

## seed: Seed the database with sample data
seed:
	@echo "$(GREEN)Seeding database...$(NC)"
	$(GO) run $(MAIN_PATH)/main.go seed

## install-tools: Install development tools
install-tools:
	@echo "$(GREEN)Installing development tools...$(NC)"
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "$(GREEN)Tools installed$(NC)"

## check: Run all checks without building
check: fmt vet lint test
	@echo "$(GREEN)All checks passed$(NC)"
