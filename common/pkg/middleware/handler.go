package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/linggaaskaedo/go-kill/common/pkg/correlation"
	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel/trace"
)

func (mw *middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if !strings.HasPrefix(path, "/swagger/") {
			start := time.Now()

			span := trace.SpanFromContext(ctx)
			spanContext := span.SpanContext()
			traceID := spanContext.TraceID().String()
			spanID := spanContext.SpanID().String()

			reqID := c.GetHeader(preference.REQUEST_ID)
			if reqID == "" {
				reqID = xid.New().String()
			}

			ctx = correlation.AttachKeyValCtx(ctx,
				preference.CONTEXT_KEY_TRACE_ID, traceID,
				preference.CONTEXT_KEY_SPAN_ID, spanID,
				preference.CONTEXT_KEY_REQ_ID, reqID,
				preference.CONTEXT_KEY_ADDR, c.Request.Host,
				preference.CONTEXT_KEY_USER_AGENT, c.Request.UserAgent())

			ctx = mw.attachLogger(ctx)

			c.Header(preference.REQUEST_ID, reqID)

			if raw != "" {
				path = path + "?" + raw
			}

			c.Request = c.Request.WithContext(ctx)
			c.Next()

			latency := time.Since(start)
			if latency > time.Minute {
				latency = latency.Truncate(time.Second)
			}

			event := mw.log.Info()
			if c.Writer.Status() >= 400 {
				event = mw.log.Error()
			}
			event.
				Str(preference.REQ_ID, reqID).
				Str(preference.METHOD, c.Request.Method).
				Str(preference.URL, path).
				Str(preference.USER_AGENT, c.Request.UserAgent()).
				Str(preference.ADDR, c.Request.Host).
				Str(preference.LATENCY, latency.String()).
				Int(preference.STATUS, c.Writer.Status()).
				Msg("request completed")
		}
	}
}

func (mw *middleware) attachLogger(ctx context.Context) context.Context {
	reqID := ""
	if id, ok := ctx.Value(preference.CONTEXT_KEY_REQ_ID).(string); ok {
		reqID = id
	}

	return mw.log.With().
		Str(string(preference.CONTEXT_KEY_REQ_ID), reqID).
		Logger().
		WithContext(ctx)
}
