# Variables
OUTPUT		:= app
BIN_DIR     := ./bin
PLATFORMS   := linux/amd64 linux/386 linux/arm linux/arm64 windows/amd64 windows/386 darwin/amd64 darwin/arm64

help: ## Show this help message
	@printf "\033[36m%-30s\033[0m %s\n" "Target" "Description"
	@printf "\033[36m%-30s\033[0m %s\n" "------" "-----------"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: build run ## Execute all steps `clean update check swagger build run`

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)/
	@rm -rf ./logs/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

update: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "Vet complete"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run
	@echo "Linting complete"

check: fmt vet lint ## Run all checks
	@echo "All checks passed"

swagger: ## Generate swagger documentation
	@echo "Generating Swagger docs..."
	@(swag fmt -d ./src 2>&1 | grep -v "warning: failed to get package name in dir") || true
	@(swag init -g ./src/cmd/app.go -o ./docs 2>&1 | grep -v "warning: failed to get package name in dir") || true
	@echo "Fixing generated docs (removing LeftDelim/RightDelim)..."
	@sed -i.bak '/LeftDelim/d' ./docs/docs.go 2>/dev/null || sed -i '/LeftDelim/d' ./docs/docs.go 2>/dev/null
	@sed -i.bak '/RightDelim/d' ./docs/docs.go 2>/dev/null || sed -i '/RightDelim/d' ./docs/docs.go 2>/dev/null
	@rm -f ./docs/docs.go.bak 2>/dev/null || true
	@echo "Swagger docs generated and fixed successfully"

build: clean check swagger ## Build the application
	@echo "Building application..."
	@go mod tidy
	@go generate ./src/cmd
	@go build -o $(BIN_DIR)/$(OUTPUT) ./src/cmd
	@echo "Build complete: bin/app"

build-all: ## Cross-compile for all operating system (e.g., linux, windows, darwin, freebsd) & architecture (e.g., amd64, 386, arm, arm64)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output_name=$(OUTPUT)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name="$$output_name.exe"; \
		fi; \
		echo "Building for $$GOOS/$$GOARCH..."; \
		env GOOS=$$GOOS GOARCH=$$GOARCH go build -o $(BIN_DIR)/$$output_name ./src/cmd; \
		if [ $$? -ne 0 ]; then \
			echo "Error building for $$GOOS/$$GOARCH"; \
			exit 1; \
		fi; \
	done
	@echo "All builds completed successfully. Binaries are in $(BIN_DIR)/"

run: ## Run the application
	@echo "Starting application..."
	@$(BIN_DIR)/$(OUTPUT)

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"