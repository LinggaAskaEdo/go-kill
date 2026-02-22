package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (mw *middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				mw.log.Error().
					Interface("panic", err).
					Msg("Recovered from panic")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
