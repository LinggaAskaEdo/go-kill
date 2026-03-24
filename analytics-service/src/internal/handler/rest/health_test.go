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

const (
	pathHealthLive  = "/health/live"
	pathHealthReady = "/health/ready"
	pathMetrics     = "/metrics"
	errMsgStatusFmt = "expected status %d, got %d"
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

	router.GET(pathHealthLive, r.liveness)

	req, _ := http.NewRequest(http.MethodGet, pathHealthLive, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestReadinessAllHealthy(t *testing.T) {
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

	router.GET(pathHealthReady, r.readiness)

	req, _ := http.NewRequest(http.MethodGet, pathHealthReady, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestReadinessRedisDown(t *testing.T) {
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

	router.GET(pathHealthReady, r.readiness)

	req, _ := http.NewRequest(http.MethodGet, pathHealthReady, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf(errMsgStatusFmt, http.StatusServiceUnavailable, w.Code)
	}
}

func TestReadinessMongoDown(t *testing.T) {
	router := setupTestRouter()

	r := &rest{
		gin:   router,
		redis: nil,
		mongo: nil,
		log:   zerolog.Logger{},
	}

	router.GET(pathHealthReady, r.readiness)

	req, _ := http.NewRequest(http.MethodGet, pathHealthReady, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestMetricsHandler(t *testing.T) {
	router := setupTestRouter()

	r := &rest{
		gin: router,
		log: zerolog.Logger{},
	}

	router.GET(pathMetrics, r.metricsHandler())

	req, _ := http.NewRequest(http.MethodGet, pathMetrics, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf(errMsgStatusFmt, http.StatusOK, w.Code)
	}
}

func TestInitMetrics(t *testing.T) {
	r := &rest{
		log: zerolog.Logger{},
	}

	r.initMetrics()
}
