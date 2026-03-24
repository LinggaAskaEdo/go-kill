package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/linggaaskaedo/go-kill/product-service/src/internal/service"
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

	router.GET("/api/v1/products", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/api/v1/products/:id", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/api/v1/categories", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/api/v1/products/:id/categories", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/api/v1/categories/:id/products", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	tests := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/api/v1/products", http.StatusOK},
		{http.MethodGet, "/api/v1/products/123", http.StatusOK},
		{http.MethodGet, "/api/v1/categories", http.StatusOK},
		{http.MethodGet, "/api/v1/products/123/categories", http.StatusOK},
		{http.MethodGet, "/api/v1/categories/123/products", http.StatusOK},
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
