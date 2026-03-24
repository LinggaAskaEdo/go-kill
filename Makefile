MODULES := analytics-service auth-service common notification-service order-service product-service user-service

build-all: 
	@echo "🔨 Building modules..."
	@for service in $(MODULES); do \
		if [ -d "$$service" ]; then \
			echo "  Building $$service..."; \
			(cd $$service && go mod tidy && go get -u ./...); \
		fi \
	done
	@echo "✅ Build modules finished"

test-all:
	@echo "🧪 Running tests..."
	@for service in $(MODULES); do \
		if [ -d "$$service/src" ]; then \
			echo "  Testing $$service..."; \
			(cd $$service && go test -v ./src/internal/handler/... 2>&1 || true); \
		fi \
	done
	@echo "✅ Tests finished"