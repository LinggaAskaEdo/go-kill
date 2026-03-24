package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linggaaskaedo/go-kill/user-service/src/internal/service"
)

const (
	pathUsersRegister     = "/api/v1/users/register"
	pathUsersMe           = "/api/v1/users/me"
	pathUsersMeActivities = "/api/v1/users/me/activities"
	pathUsersMeAddresses  = "/api/v1/users/me/addresses"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestInitRestHandler(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}

	InitRestHandler(router, svc)
}

func TestServeRoutesRegistered(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}

	handler := &rest{
		gin: router,
		svc: svc,
	}

	handler.Serve()

	tests := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodPost, pathUsersRegister, http.StatusBadRequest},
		{http.MethodGet, pathUsersMe, http.StatusUnauthorized},
		{http.MethodGet, pathUsersMeActivities, http.StatusUnauthorized},
		{http.MethodGet, pathUsersMeAddresses, http.StatusUnauthorized},
		{http.MethodPost, pathUsersMeAddresses, http.StatusUnauthorized},
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

func TestAuthMiddlewareNoHeader(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}

	handler := &rest{
		gin: router,
		svc: svc,
	}

	router.GET(pathUsersMe, handler.authMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}

	handler := &rest{
		gin: router,
		svc: svc,
	}

	router.GET(pathUsersMe, handler.authMiddleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, pathUsersMe, nil)
	req.Header.Set("Authorization", "InvalidToken")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
