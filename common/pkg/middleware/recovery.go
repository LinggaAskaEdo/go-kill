package middleware

import (
	"net/http"

	"github.com/linggaaskaedo/go-kill/common/pkg/correlation"
	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/gin-gonic/gin"
)

func (mw *middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				mw.log.Error().
					Str(preference.REQ_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_REQ_ID)).
					Interface("panic", err).
					Msg("Recovered from panic")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
