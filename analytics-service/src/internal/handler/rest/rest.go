package rest

import (
	"sync"

	"github.com/linggaaskaedo/go-kill/analytics-service/src/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var onceRestHandler = &sync.Once{}

type rest struct {
	gin   *gin.Engine
	svc   *service.Service
	redis *redis.Client
	mongo *mongo.Database
	log   zerolog.Logger
}

func InitRestHandler(gin *gin.Engine, svc *service.Service, redis *redis.Client, mongo *mongo.Database) {
	var e *rest

	onceRestHandler.Do(func() {
		e = &rest{
			gin:   gin,
			svc:   svc,
			redis: redis,
			mongo: mongo,
			log:   zerolog.Logger{},
		}

		e.Serve()
	})
}

func (e *rest) Serve() {
	e.gin.GET("/health/live", e.liveness)
	e.gin.GET("/health/ready", e.readiness)
	e.gin.GET("/metrics", e.metricsHandler())

	e.initMetrics()
}
