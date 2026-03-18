package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins []string
var allowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}

func SetAllowedOrigins(origins []string) {
	allowedOrigins = origins
}

func (mw *middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(allowedOrigins) > 0 {
			if !slices.Contains(allowedOrigins, "*") && !slices.Contains(allowedOrigins, origin) {
				c.Next()
				return
			}
		}

		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem 'self'; img-src * data:; style-src * 'self'")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")

		if len(allowedOrigins) > 0 && !slices.Contains(allowedOrigins, "*") {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		if !slices.Contains(allowedMethods, c.Request.Method) {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		c.Next()
	}
}
