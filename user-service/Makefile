help: ## Show this help message
	@printf "\033[36m%-30s\033[0m %s\n" "Target" "Description"
	@printf "\033[36m%-30s\033[0m %s\n" "------" "-----------"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: build run ## Execute all steps `clean update check swagger build run`

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf ./bin/
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

build: clean update check swagger ## Build the application
	@echo "Building application..."
	@go mod tidy
	@go generate ./src/cmd
	@go build -o ./bin/app ./src/cmd
	@echo "Build complete: bin/app"

run: ## Run the application
	@echo "Starting application..."
	@./bin/app

deps: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Tools installed"

sql-postgres-create: ## Create SQL migration files for postgres
	@echo "Creating postgres SQL migration files..."
	@read -p "Enter migration name (use underscores): " name; \
		goose -dir ./etc/migrations/postgres create postgres_$${name} sql

sql-postgres-up: ## Apply up migrations for postgres
	@echo "Applying up migrations for postgres..."; \
		{ \
			stty -echo ; \
			trap 'stty echo' EXIT ; \
			read -p "Enter postgres password: " pass ; \
			stty echo ; \
			echo ; \
			goose -dir ./etc/migrations/postgres postgres "host=localhost user=postgres password=$$pass dbname=go-kill sslmode=disable" up ; \
		}
