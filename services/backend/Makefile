.PHONY: help build run test clean docker-build migrate-up migrate-down lint fmt

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building..."
	@go build -o bin/server cmd/server/main.go

run: ## Run the application
	@echo "Running..."
	@go run cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

coverage: test ## Run tests with coverage
	@go tool cover -html=coverage.out

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t launchdate-backend:latest .

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	@go mod tidy

.DEFAULT_GOAL := help
