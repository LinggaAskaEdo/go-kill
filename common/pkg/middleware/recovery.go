package middleware

import (
	"fmt"
	"net/http"

	"github.com/linggaaskaedo/go-kill/common/pkg/correlation"
	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/gin-gonic/gin"
)

func (mw *middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// log.Error().
				// 	// Str(preference.TRACE_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_TRACE_ID)).
				// 	// Str(preference.SPAN_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_SPAN_ID)).
				// 	Str(preference.REQ_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_REQ_ID)).
				// 	Interface("panic", err).
				// 	Msg("Recovered from panic")

				// Print raw value for debugging
				fmt.Printf("PANIC RAW: %#v\n", err)

				// // Get or generate request ID
				// reqID := c.GetHeader(preference.REQ_ID)
				// if reqID == "" {
				// 	reqID = xid.New().String()
				// }

				// Log the panic using the component's logger (includes caller info)
				mw.log.Error().
					// Str(preference.REQ_ID, reqID).
					Str(preference.REQ_ID, correlation.GetCtxKeyVal(c, preference.CONTEXT_KEY_REQ_ID)).
					Interface("panic", err).
					Msg("Recovered from panic")

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
