package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestInitRestHandler(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}
	_ = svc

	mockRedis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = mockRedis

	mockMongo := &mongo.Database{}
	_ = mockMongo

	InitRestHandler(router, svc, nil, nil)
}

func TestServe_RoutesRegistered(t *testing.T) {
	router := setupTestRouter()

	svc := &service.Service{}

	handler := &rest{
		gin:   router,
		svc:   svc,
		redis: nil,
		mongo: nil,
		log:   zerolog.Logger{},
	}

	handler.Serve()

	tests := []struct {
		method   string
		path     int
		expected int
	}{
		{http.MethodGet, http.StatusOK, http.StatusOK},
		{http.MethodGet, http.StatusOK, http.StatusOK},
		{http.MethodGet, http.StatusOK, http.StatusOK},
	}

	routes := []string{pathHealthLive, pathHealthReady, pathMetrics}

	for i, tt := range tests {
		req, _ := http.NewRequest(tt.method, routes[i], nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != tt.expected {
			t.Errorf("route %s %s: expected %d, got %d", tt.method, routes[i], tt.expected, w.Code)
		}
	}
}

func TestServe_LivenessRoute(t *testing.T) {
	router := setupTestRouter()

	handler := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	handler.Serve()

	req, _ := http.NewRequest(http.MethodGet, pathHealthLive, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestServe_ReadinessRoute(t *testing.T) {
	router := setupTestRouter()

	handler := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	handler.Serve()

	req, _ := http.NewRequest(http.MethodGet, pathHealthReady, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestServe_MetricsRoute(t *testing.T) {
	router := setupTestRouter()

	handler := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	handler.Serve()

	req, _ := http.NewRequest(http.MethodGet, pathMetrics, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}
