package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/handler/grpc"
	"github.com/linggaaskaedo/go-kill/auth-service/src/internal/service"
)

const (
	pathHealth      = "/health"
	pathAuthLogin   = "/api/v1/auth/login"
	pathAuthRefresh = "/api/v1/auth/refresh"
	pathAuthLogout  = "/api/v1/auth/logout"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestInitRestHandler(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}
	jwtSecret := "test-secret"

	InitRestHandler(router, svc, grpcHandler, jwtSecret)
}

func TestServeRoutesRegistered(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}

	handler := &rest{
		gin:       router,
		svc:       svc,
		grpc:      grpcHandler,
		jwtSecret: []byte("test"),
	}

	handler.Serve()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodPost, pathAuthLogin, http.StatusBadRequest},
		{http.MethodPost, pathAuthRefresh, http.StatusBadRequest},
		{http.MethodPost, pathAuthLogout, http.StatusUnauthorized},
		{http.MethodGet, pathHealth, http.StatusOK},
	}

	for _, tt := range tests {
		req, _ := http.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.status {
			t.Errorf("route %s %s: expected %d, got %d", tt.method, tt.path, tt.status, w.Code)
		}
	}
}

func TestHandleHealth(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}

	handler := &rest{
		gin:       router,
		svc:       svc,
		grpc:      grpcHandler,
		jwtSecret: []byte("test"),
	}

	router.GET(pathHealth, handler.handleHealth)

	req, _ := http.NewRequest(http.MethodGet, pathHealth, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandleLoginInvalidBody(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}

	handler := &rest{
		gin:       router,
		svc:       svc,
		grpc:      grpcHandler,
		jwtSecret: []byte("test"),
	}

	router.POST(pathAuthLogin, handler.handleLogin)

	req, _ := http.NewRequest(http.MethodPost, pathAuthLogin, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusOK {
		t.Errorf("expected status %d or %d, got %d", http.StatusBadRequest, http.StatusOK, w.Code)
	}
}

func TestHandleRefreshInvalidBody(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}

	handler := &rest{
		gin:       router,
		svc:       svc,
		grpc:      grpcHandler,
		jwtSecret: []byte("test"),
	}

	router.POST(pathAuthRefresh, handler.handleRefresh)

	req, _ := http.NewRequest(http.MethodPost, pathAuthRefresh, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusOK {
		t.Errorf("expected status %d or %d, got %d", http.StatusBadRequest, http.StatusOK, w.Code)
	}
}

func TestHandleLogoutNoAuthorization(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	grpcHandler := &grpc.Grpc{}

	handler := &rest{
		gin:       router,
		svc:       svc,
		grpc:      grpcHandler,
		jwtSecret: []byte("test"),
	}

	router.POST(pathAuthLogout, handler.handleLogout)

	req, _ := http.NewRequest(http.MethodPost, pathAuthLogout, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
