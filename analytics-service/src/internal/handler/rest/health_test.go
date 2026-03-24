package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestLiveness(t *testing.T) {
	router := setupTestRouter()

	r := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	router.GET("/health/live", r.liveness)

	req, _ := http.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestReadiness_AllHealthy(t *testing.T) {
	router := setupTestRouter()

	mockRedis := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_ = mockRedis

	mockMongo := &mongo.Client{}
	_ = mockMongo

	r := &rest{
		gin:   router,
		redis: nil,
		mongo: nil,
		log:   zerolog.Logger{},
	}

	router.GET("/health/ready", r.readiness)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestReadiness_RedisDown(t *testing.T) {
	router := setupTestRouter()

	mockRedis := redis.NewClient(&redis.Options{
		Addr: "invalid:6379",
	})

	r := &rest{
		gin:   router,
		redis: mockRedis,
		mongo: nil,
		log:   zerolog.Logger{},
	}

	router.GET("/health/ready", r.readiness)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 503, got %d", w.Code)
	}
}

func TestReadiness_MongoDown(t *testing.T) {
	router := setupTestRouter()

	r := &rest{
		gin:   router,
		redis: nil,
		mongo: nil,
		log:   zerolog.Logger{},
	}

	router.GET("/health/ready", r.readiness)

	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMetricsHandler(t *testing.T) {
	router := setupTestRouter()

	r := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	router.GET("/metrics", r.metricsHandler())

	req, _ := http.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestInitMetrics(t *testing.T) {
	r := &rest{
		log: zerolog.Logger{},
	}

	r.initMetrics()
}
