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

func LoggingUnaryServerInterceptor(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		reqID := correlation.GetReqID(ctx, preference.CONTEXT_KEY_REQ_ID)
		if reqID == "" {
			reqID = xid.New().String()
			ctx = correlation.WithReqID(ctx, preference.CONTEXT_KEY_REQ_ID, reqID)
		}

		start := time.Now()
		resp, err := handler(ctx, req)
		latency := time.Since(start)

		event := logger.Info()
		if err != nil {
			event = logger.Error().Err(err)
		}
		event.
			Str("req_id", reqID).
			Str("method", info.FullMethod).
			Dur("latency", latency).
			Int("status_code", int(status.Code(err))).
			Msg("gRPC unary request completed")

		return resp, err
	}
}

func LoggingStreamServerInterceptor(logger zerolog.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		reqID := correlation.GetReqID(ss.Context(), preference.CONTEXT_KEY_REQ_ID)
		if reqID == "" {
			reqID = xid.New().String()
		}

		start := time.Now()
		err := handler(srv, ss)
		latency := time.Since(start)

		event := logger.Info()
		if err != nil {
			event = logger.Error().Err(err)
		}
		event.
			Str("req_id", reqID).
			Str("method", info.FullMethod).
			Str("stream_type", "server_stream").
			Dur("latency", latency).
			Int("status_code", int(status.Code(err))).
			Msg("gRPC stream request completed")

		return err
	}
}
