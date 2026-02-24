package grpcserver

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/correlation"
	"github.com/linggaaskaedo/go-kill/common/preference"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *GRPCServerComponent) ReqIDServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	// Extract metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Get correlation ID from header
	var corrID string
	if vals := md.Get(preference.REQ_ID); len(vals) > 0 && vals[0] != "" {
		corrID = vals[0]
	} else {
		corrID = xid.New().String()
	}

	// Store in context
	// ctx = correlation.WithReqID(ctx, preference.CONTEXT_KEY_REQ_ID, corrID)
	ctx = correlation.WithCorrelationID(ctx, preference.CONTEXT_KEY_REQ_ID, corrID)

	// Enhance logger with correlation ID
	// logger := zerolog.Ctx(ctx).With().Str(preference.REQ_ID, corrID).Logger()
	// ctx = logger.WithContext(ctx)
	ctx = s.log.With().Str(preference.REQ_ID, corrID).Logger().WithContext(ctx)

	// return mw.log.With().
	// 	Str(string(preference.CONTEXT_KEY_REQ_ID), ctx.Value(preference.CONTEXT_KEY_REQ_ID).(string)).
	// 	Logger().
	// 	WithContext(ctx)

	return handler(ctx, req)
}
