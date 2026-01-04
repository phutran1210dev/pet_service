.PHONY: help build run dev docker-build docker-up docker-down docker-logs clean test swagger

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@command -v swag >/dev/null 2>&1 || { echo "Installing swag..."; go install github.com/swaggo/swag/cmd/swag@latest; }
	@if [ -f ~/go/bin/swag ]; then ~/go/bin/swag init -g main.go --parseDependency --parseInternal; \
	elif [ -f $(shell go env GOPATH)/bin/swag ]; then $(shell go env GOPATH)/bin/swag init -g main.go --parseDependency --parseInternal; \
	else swag init -g main.go --parseDependency --parseInternal; fi
	@echo "Swagger documentation generated successfully!"

build: ## Build the Go application
	go build -o main .

run: ## Run the application locally
	go run main.go

dev: ## Run with air hot reload
	air -c .air.toml

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start all services with docker-compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show docker-compose logs
	docker-compose logs -f

docker-dev-up: ## Start development environment with hot reload
	docker-compose -f docker-compose.dev.yml up

docker-dev-down: ## Stop development environment
	docker-compose -f docker-compose.dev.yml down

docker-restart: ## Restart all services
	docker-compose restart

docker-clean: ## Remove all containers, volumes, and images
	docker-compose down -v
	docker system prune -f

clean: ## Clean build artifacts
	rm -f main
	rm -rf tmp/
	rm -f build-errors.log

test: ## Run tests
	go test -v ./...

deps: ## Download dependencies
	go mod download

tidy: ## Tidy dependencies
	go mod tidy

install-air: ## Install air for hot reload
	go install github.com/air-verse/air@latest

migrate-up: ## Run database migrations (placeholder)
	@echo "Add migration command here"

migrate-down: ## Rollback database migrations (placeholder)
	@echo "Add migration rollback command here"

lint: ## Run golangci-lint
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Production targets
prod-build: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

prod-up: ## Start production services
	docker-compose up -d --build

prod-down: ## Stop production services
	docker-compose down

prod-logs: ## Show production logs
	docker-compose logs -f pet-service
