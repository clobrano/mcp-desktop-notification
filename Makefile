.PHONY: help build test clean install lint run

# Variables
BINARY_NAME=mcp-poke
BUILD_DIR=./bin
CMD_DIR=./cmd/mcp-poke
VERSION?=dev
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary for current platform
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run tests
	@echo "Running tests..."
	@go test -short -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -short -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(CMD_DIR)
	@echo "Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)"

lint: ## Run linters
	@echo "Running linters..."
	@go vet ./...
	@go fmt ./...
	@echo "Linting complete"

run: build ## Build and run the binary
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Cross-compilation targets
build-linux: ## Build for Linux AMD64
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@echo "Linux binary built: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-macos: ## Build for macOS ARM64
	@echo "Building for macOS ARM64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64 $(CMD_DIR)
	@echo "macOS ARM64 binary built: $(BUILD_DIR)/$(BINARY_NAME)-macos-arm64"

build-macos-amd64: ## Build for macOS AMD64
	@echo "Building for macOS AMD64..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64 $(CMD_DIR)
	@echo "macOS AMD64 binary built: $(BUILD_DIR)/$(BINARY_NAME)-macos-amd64"

build-windows: ## Build for Windows AMD64
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "Windows binary built: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-all: build-linux build-macos build-macos-amd64 build-windows ## Build for all platforms
	@echo "All platform binaries built"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"
