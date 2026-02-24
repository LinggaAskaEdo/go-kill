package middleware

import (
	"net/http"

	"github.com/linggaaskaedo/go-kill/common/correlation"
	"github.com/linggaaskaedo/go-kill/common/preference"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (mw *middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					// Str(preference.TRACE_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_TRACE_ID)).
					// Str(preference.SPAN_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_SPAN_ID)).
					Str(preference.REQ_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_REQ_ID)).
					Interface("panic", err).
					Msg("Recovered from panic")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
