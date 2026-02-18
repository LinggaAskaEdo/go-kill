package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logEvent := log.Info()
		if statusCode >= 400 {
			logEvent = log.Error()
		}

		logEvent.
			Str("request_id", GetRequestID(c)).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("ip", clientIP).
			Dur("latency", latency).
			Str("error", errorMessage).
			Msg("HTTP request handled")
	}
}

func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("requestID"); exists {
		return id.(string)
	}

	return ""
}
