package grpcserver

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/common/pkg/correlation"
	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggingUnaryServerInterceptor returns a new unary server interceptor that logs request/response.
func LoggingUnaryServerInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// Get or generate correlation ID (it should have been set by your ReqIDServerInterceptor)
		reqID := correlation.GetReqID(ctx, preference.CONTEXT_KEY_REQ_ID)
		if reqID == "" {
			// Fallback: generate one if missing (should not happen if ReqIDServerInterceptor runs first)
			reqID = xid.New().String()
			ctx = correlation.WithReqID(ctx, preference.CONTEXT_KEY_REQ_ID, reqID)
		}

		logger.Info().
			Str(preference.EVENT, "START").
			Str("req_id", reqID).
			Str("method", info.FullMethod).
			Msg("gRPC request started")

		start := time.Now()

		// Call the handler
		resp, err := handler(ctx, req)

		// Log END event with duration and result
		if err != nil {
			// For errors, you might want to log at error level
			logger.Error().
				Str("req_id", reqID).
				Err(err).
				Send()
		}

		logger.Info().
			Str("event", "END").
			Str("req_id", reqID).
			Dur("latency", time.Since(start)).
			Int("status_code", int(status.Code(err))). // gRPC status code as integer
			Msg("gRPC request completed")

		return resp, err
	}
}
