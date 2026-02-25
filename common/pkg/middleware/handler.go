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

		if !strings.HasPrefix(path, "/swagger/") { // skip logging swagger request
			start := time.Now()

			// Get trace context from OpenTelemetry
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

			mw.log.Info().
				Str(preference.EVENT, "START").
				// Str(preference.TRACE_ID, traceID).
				// Str(preference.SPAN_ID, spanID).
				Str(preference.REQ_ID, reqID).
				Str(preference.METHOD, c.Request.Method).
				Str(preference.URL, path).
				Str(preference.USER_AGENT, c.Request.UserAgent()).
				Str(preference.ADDR, c.Request.Host).
				Send()

			// Process request
			c.Request = c.Request.WithContext(ctx)
			c.Next()

			// Fill the params
			param := gin.LogFormatterParams{}

			param.TimeStamp = time.Now() // Stop timer
			param.Latency = param.TimeStamp.Sub(start)
			if param.Latency > time.Minute {
				param.Latency = param.Latency.Truncate(time.Second)
			}

			param.StatusCode = c.Writer.Status()

			mw.log.Info().
				Str(preference.EVENT, "END").
				// Str(preference.TRACE_ID, traceID).
				// Str(preference.SPAN_ID, spanID).
				Str(preference.REQ_ID, reqID).
				Str(preference.LATENCY, param.Latency.String()).
				Int(preference.STATUS, param.StatusCode).
				Send()
		}
	}
}

func (mw *middleware) attachLogger(ctx context.Context) context.Context {
	return mw.log.With().
		// Str(string(preference.CONTEXT_KEY_TRACE_ID), ctx.Value(preference.CONTEXT_KEY_TRACE_ID).(string)).
		// Str(string(preference.CONTEXT_KEY_SPAN_ID), ctx.Value(preference.CONTEXT_KEY_SPAN_ID).(string)).
		Str(string(preference.CONTEXT_KEY_REQ_ID), ctx.Value(preference.CONTEXT_KEY_REQ_ID).(string)).
		Logger().
		WithContext(ctx)
}
