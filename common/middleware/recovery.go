package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (mw *middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Str("trace_id", mw.getTraceID(c)).
					Str("span_id", mw.getSpanID(c)).
					Str("req_id", mw.getReqID(c)).
					Interface("panic", err).
					Msg("Recovered from panic")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
