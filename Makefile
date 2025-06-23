# Makefile for URL Shortener

# Variables
BINARY_NAME=url-shortener
BINARY_PATH=bin/$(BINARY_NAME)
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")
CONFIG_FILE=configs/config.yaml
DB_FILE=url_shortener.db

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build clean run migrate create stats test test-simple dev install deps fmt vet lint check

# Default target
help: ## Show this help message
	@echo "$(BLUE)URL Shortener - Available Make Commands$(NC)"
	@echo "======================================="
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Build commands
build: ## Build the application
	@echo "$(BLUE)🔨 Building $(BINARY_NAME)...$(NC)"
	@mkdir -p bin
	@go build -o $(BINARY_PATH) .
	@echo "$(GREEN)✅ Build complete: $(BINARY_PATH)$(NC)"

clean: ## Clean build artifacts and database
	@echo "$(YELLOW)🧹 Cleaning up...$(NC)"
	@rm -rf bin/
	@rm -f $(DB_FILE)
	@rm -f server.log
	@rm -f coverage.out coverage.html
	@rm -f *.test
	@echo "$(GREEN)✅ Cleanup complete$(NC)"

clean-test: ## Clean test artifacts
	@echo "$(YELLOW)🧹 Cleaning test artifacts...$(NC)"
	@rm -f coverage.out coverage.html
	@rm -f *.test
	@find . -name "*.test" -delete 2>/dev/null || true
	@echo "$(GREEN)✅ Test cleanup complete$(NC)"

install: deps build ## Install dependencies and build

deps: ## Download and tidy dependencies
	@echo "$(BLUE)📦 Installing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies installed$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)🔧 Installing development tools...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(GREEN)✅ Development tools installed$(NC)"

# Development commands
dev: build migrate ## Build and setup for development
	@echo "$(GREEN)🚀 Development setup complete!$(NC)"
	@echo "Run 'make run' to start the server"

fmt: ## Format Go code
	@echo "$(BLUE)📝 Formatting code...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(BLUE)🔍 Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet complete$(NC)"

lint: ## Run golangci-lint
	@echo "$(BLUE)🔍 Running golangci-lint...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)⚠️ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi
	@echo "$(GREEN)✅ Linting complete$(NC)"


check: fmt vet lint ## Run code quality checks
	@echo "$(GREEN)✅ All checks passed$(NC)"

check-all: fmt vet lint security vuln-check ## Run all quality and security checks
	@echo "$(GREEN)✅ All quality and security checks passed$(NC)"

pre-commit: clean fmt vet lint test-coverage ## Run pre-commit checks (format, lint, test)
	@echo "$(GREEN)🚀 Pre-commit checks passed! Ready to commit.$(NC)"

ci: clean check-all test-all ## Simulate CI pipeline locally
	@echo "$(GREEN)🎉 CI simulation completed successfully!$(NC)"

# Database commands
migrate: build ## Run database migrations
	@echo "$(BLUE)🗃️  Running database migrations...$(NC)"
	@$(BINARY_PATH) migrate
	@echo "$(GREEN)✅ Migrations complete$(NC)"

# Application commands
run: build ## Start the server
	@echo "$(BLUE)🚀 Starting URL shortener server...$(NC)"
	@$(BINARY_PATH) run-server

create: build ## Create a short URL (usage: make create URL=https://example.com)
	@if [ -z "$(URL)" ]; then \
		echo "$(RED)❌ Usage: make create URL=https://example.com$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)🔗 Creating short URL for: $(URL)$(NC)"
	@$(BINARY_PATH) create --url="$(URL)"

stats: build ## Get statistics for a short code (usage: make stats CODE=abc123)
	@if [ -z "$(CODE)" ]; then \
		echo "$(RED)❌ Usage: make stats CODE=abc123$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)📊 Getting statistics for: $(CODE)$(NC)"
	@$(BINARY_PATH) stats --code="$(CODE)"

# Testing commands
test: ## Run all unit tests
	@echo "$(BLUE)🧪 Running unit tests...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✅ Unit tests complete$(NC)"

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)🧪 Running tests with coverage...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

test-race: ## Run tests with race detection
	@echo "$(BLUE)🧪 Running tests with race detection...$(NC)"
	@go test -v -race ./...
	@echo "$(GREEN)✅ Race condition tests complete$(NC)"

test-bench: ## Run benchmark tests
	@echo "$(BLUE)🧪 Running benchmark tests...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo "$(GREEN)✅ Benchmark tests complete$(NC)"

test-integration: build ## Run integration tests
	@echo "$(BLUE)🧪 Running integration tests...$(NC)"
	@if [ -f test_workflow.sh ]; then \
		chmod +x test_workflow.sh && ./test_workflow.sh; \
	else \
		echo "$(YELLOW)⚠️ test_workflow.sh not found, skipping integration tests$(NC)"; \
	fi
	@echo "$(GREEN)✅ Integration tests complete$(NC)"
	
test-all: test test-race test-coverage test-bench ## Run all types of tests
	@echo "$(GREEN)✅ All tests completed successfully!$(NC)"

# Release commands
build-all: ## Build for multiple platforms
	@echo "$(BLUE)🔨 Building for multiple platforms...$(NC)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 .
	@GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)✅ Multi-platform build complete$(NC)"

# Quick development workflow
quick-test: clean build migrate ## Quick development test
	@echo "$(BLUE)⚡ Quick test workflow...$(NC)"
	@$(BINARY_PATH) create --url="https://www.google.com"
	@echo "$(GREEN)✅ Quick test complete$(NC)"

# Show project status
status: ## Show project status and info
	@echo "$(BLUE)📋 Project Status$(NC)"
	@echo "=================="
	@echo "Binary: $(BINARY_PATH)"
	@echo "Config: $(CONFIG_FILE)"
	@echo "Database: $(DB_FILE)"
	@echo "Go version: $$(go version)"
	@echo "Files: $$(echo $(GO_FILES) | wc -w) Go files"
	@if [ -f $(BINARY_PATH) ]; then echo "$(GREEN)✅ Binary exists$(NC)"; else echo "$(YELLOW)⚠️  Binary not built$(NC)"; fi
	@if [ -f $(CONFIG_FILE) ]; then echo "$(GREEN)✅ Config exists$(NC)"; else echo "$(RED)❌ Config missing$(NC)"; fi
	@if [ -f $(DB_FILE) ]; then echo "$(GREEN)✅ Database exists$(NC)"; else echo "$(YELLOW)⚠️  Database not initialized$(NC)"; fi

# Development helpers
logs: ## Show server logs
	@if [ -f server.log ]; then tail -f server.log; else echo "$(RED)❌ No server.log found$(NC)"; fi

kill: ## Kill any running url-shortener processes
	@echo "$(YELLOW)🔪 Killing url-shortener processes...$(NC)"
	@pkill -f url-shortener || echo "No processes found"
	@echo "$(GREEN)✅ Processes killed$(NC)"

# Example workflows
demo: clean build migrate ## Run a complete demo
	@echo "$(BLUE)🎬 Running demo...$(NC)"
	@echo "$(YELLOW)Starting server in background...$(NC)"
	@$(BINARY_PATH) run-server > server.log 2>&1 &
	@sleep 3
	@echo "$(YELLOW)Creating test URLs...$(NC)"
	@$(BINARY_PATH) create --url="https://www.google.com"
	@$(BINARY_PATH) create --url="https://github.com"
	@echo "$(YELLOW)Stopping server...$(NC)"
	@pkill -f url-shortener || true
	@echo "$(GREEN)✅ Demo complete$(NC)"