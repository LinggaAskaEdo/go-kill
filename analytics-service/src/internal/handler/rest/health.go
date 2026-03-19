package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/linggaaskaedo/go-kill/common/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (e *rest) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "alive"})
}

func (e *rest) readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if e.redis != nil {
		if err := e.redis.Ping(ctx).Err(); err != nil {
			e.log.Error().Err(err).Msg("Redis health check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"redis":  "down",
			})

			return
		}
	}

	if e.mongo != nil {
		if err := e.mongo.Client().Ping(ctx, nil); err != nil {
			e.log.Error().Err(err).Msg("MongoDB health check failed")
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"mongo":  "down",
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"redis":  "up",
		"mongo":  "up",
	})
}

func (e *rest) metricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (e *rest) initMetrics() {
	metrics.MessagesReceived.WithLabelValues("")
	metrics.MessagesProcessed.WithLabelValues("", "")
	metrics.MessageProcessingDuration.WithLabelValues("")
	metrics.DLQMessagesSent.WithLabelValues("")
	metrics.RetryAttempts.WithLabelValues("")
	metrics.MongoOperations.WithLabelValues("", "", "")
	metrics.MongoOperationDuration.WithLabelValues("", "")
	metrics.RedisOperations.WithLabelValues("", "")
}
