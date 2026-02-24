package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/linggaaskaedo/go-kill/common/correlation"
	"github.com/linggaaskaedo/go-kill/common/preference"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Config struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type ServiceRegistrar func(*grpc.Server)

type GRPCServerComponent struct {
	log        zerolog.Logger
	cfg        Config
	registrars []ServiceRegistrar
	ready      chan struct{}
	server     *grpc.Server
	lis        net.Listener
}

// NewGRPCServerComponent creates a new server component with the given service registrars.
func NewGRPCServerComponent(log zerolog.Logger, cfg Config, registrars ...ServiceRegistrar) *GRPCServerComponent {
	return &GRPCServerComponent{
		log:        log,
		cfg:        cfg,
		registrars: registrars,
		ready:      make(chan struct{}),
	}
}

// Start creates the listener, registers services, and begins serving.
// It blocks until the context is cancelled or the server fails.
func (s *GRPCServerComponent) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.Port, err)
	}
	s.lis = lis

	s.log.Info().Str("port", s.cfg.Port).Msg("gRPC server listening")

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(s.correlationIDServerInterceptor))
	for _, registrar := range s.registrars {
		registrar(grpcServer)
	}

	s.server = grpcServer

	serveErr := make(chan error, 1)
	go func() {
		s.log.Debug().Str("port", s.cfg.Port).Msg("gRPC server starting")
		if err := grpcServer.Serve(lis); err != nil {
			serveErr <- err
		}
	}()

	close(s.ready) // signal readiness
	s.log.Debug().Msg("gRPC server connected")

	// Wait for shutdown signal or serve error
	select {
	case <-ctx.Done():
		s.log.Debug().Msg("gRPC server context cancelled – stopping")

		return nil
	case err := <-serveErr:
		return fmt.Errorf("gRPC server serve error: %w", err)
	}
}

// Stop gracefully stops the server, waiting for ongoing requests to finish.
func (s *GRPCServerComponent) Stop(ctx context.Context) error {
	s.log.Debug().Msg("Starting shut down gRPC server")
	if s.server == nil {
		return nil
	}

	// Use a separate context for the graceful stop deadline
	stopCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	s.log.Debug().Msg("Shutting down gRPC server gracefully")
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.log.Debug().Msg("gRPC server stopped gracefully")
		return nil
	case <-stopCtx.Done():
		s.log.Debug().Msg("Timeout shutting down gRPC server, forcing Stop")
		s.server.Stop() // force stop

		return fmt.Errorf("gRPC server shutdown timed out, forced stop")
	}
}

// Ready returns a channel that is closed when the connection is established.
func (s *GRPCServerComponent) Ready() <-chan struct{} {
	return s.ready
}

func (s *GRPCServerComponent) correlationIDServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
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
	ctx = correlation.WithReqID(ctx, preference.CONTEXT_KEY_REQ_ID, corrID)

	// Enhance logger with correlation ID
	ctx = s.log.With().Str(preference.REQ_ID, corrID).Logger().WithContext(ctx)

	return handler(ctx, req)
}
