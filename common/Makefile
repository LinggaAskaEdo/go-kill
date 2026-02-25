clean:
	@echo "Cleaning..."
	@find ./pkg/proto -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -exec rm -f {} \;
	@find ./pkg/proto -mindepth 1 -type d -empty -delete
	@echo "Clean complete"

install:
	@echo "Installing dependencies..."
	@go get github.com/rs/zerolog \
		golang.org/x/sync/errgroup \
		github.com/jmoiron/sqlx \
		github.com/gin-gonic/gin \
		github.com/google/uuid \
		github.com/lib/pq \
		google.golang.org/protobuf/cmd/protoc-gen-go@latest \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Dependencies installed"

update:
	@echo "Updating dependencies..."
	@go mod tidy
	@go get -u ./...	
	@echo "Dependencies updated"

proto:
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

build: clean install update proto
	@echo "Building application..."
	@go mod tidy
	@echo "Build complete"