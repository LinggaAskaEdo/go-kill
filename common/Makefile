install:
	@echo "Installing dependencies..."
	@go get github.com/rs/zerolog \
		golang.org/x/sync/errgroup \
		github.com/jmoiron/sqlx \
		github.com/gin-gonic/gin \
		github.com/google/uuid \
		github.com/lib/pq
	@echo "Dependencies installed"

update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated"

build: install update
	@echo "Building application..."
	@go mod tidy
	@echo "Build complete"