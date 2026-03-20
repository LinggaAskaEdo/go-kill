package grpcserver

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/pkg/correlation"
	"github.com/linggaaskaedo/go-kill/common/pkg/preference"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *GRPCServerComponent) ReqIDServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	var corrID string
	if vals := md.Get(preference.REQ_ID); len(vals) > 0 && vals[0] != "" {
		corrID = vals[0]
	} else {
		corrID = xid.New().String()
	}

	ctx = correlation.WithReqID(ctx, preference.CONTEXT_KEY_REQ_ID, corrID)
	ctx = s.log.With().Str(preference.REQ_ID, corrID).Logger().WithContext(ctx)

	return handler(ctx, req)
}

func (s *GRPCServerComponent) StreamServerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		md = metadata.New(nil)
	}

	var corrID string
	if vals := md.Get(preference.REQ_ID); len(vals) > 0 && vals[0] != "" {
		corrID = vals[0]
	} else {
		corrID = xid.New().String()
	}

	ctx := correlation.WithReqID(ss.Context(), preference.CONTEXT_KEY_REQ_ID, corrID)
	ctx = s.log.With().Str(preference.REQ_ID, corrID).Logger().WithContext(ctx)

	wrapped := &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	}

	return handler(srv, wrapped)
}

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStream) Context() context.Context {
	return s.ctx
}
