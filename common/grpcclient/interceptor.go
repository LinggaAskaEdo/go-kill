package grpcclient

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/correlation"
	"github.com/linggaaskaedo/go-kill/common/preference"

	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (c *GRPCClientComponent) ReqIDClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Get correlation ID from context
	corrID := correlation.GetCorrelationID(ctx, preference.CONTEXT_KEY_REQ_ID)
	if corrID == "" {
		// If no correlation ID in context, generate one (but ideally it should be present from HTTP)
		// You might want to log a warning, but we can generate to maintain traceability.
		corrID = xid.New().String()
		// Optionally store back? Not necessary.
	}

	// Attach to outgoing metadata
	ctx = metadata.AppendToOutgoingContext(ctx, preference.REQ_ID, corrID)

	return invoker(ctx, method, req, reply, cc, opts...)
}
