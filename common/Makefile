help: ## Show this help message
	@printf "\033[36m%-30s\033[0m %s\n" "Target" "Description"
	@printf "\033[36m%-30s\033[0m %s\n" "------" "-----------"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[33m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@find ./pkg/proto -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -exec rm -f {} \;
	@find ./pkg/proto -mindepth 1 -type d -empty -delete
	@echo "Clean complete"

update: ## Update dependencies
	@echo "Updating dependencies..."
	@go mod tidy
	@go get -u ./...	
	@echo "Dependencies updated"

install-tools: ## Install development tools
	@echo "Installing tools..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Tools installed"

proto: ## Compile proto file
	@echo "Generating protobuf files..."
	@for dir in pkg/proto/*/; do \
		if [ -d "$$dir" ]; then \
			echo "Generating in $$dir"; \
			(cd "$$dir" && protoc --go_out=. \
				--go_opt=paths=source_relative \
				--go-grpc_out=. \
				--go-grpc_opt=paths=source_relative \
				*.proto); \
		fi \
	done
	@echo "Protobuf files generated successfully"

build: clean install update proto ## Build the application
	@echo "Building application..."
	@go mod tidy
	@echo "Build complete"

deps: ## Installing dependencies
	@echo "Installing dependencies..."
	@go mod download
	@echo "Dependencies installed"
