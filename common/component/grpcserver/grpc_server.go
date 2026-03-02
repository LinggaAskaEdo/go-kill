package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Config struct {
	Port            string        `yaml:"port"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type GRPCServerComponent struct {
	log       zerolog.Logger
	cfg       Config
	registrar func(context.Context, *grpc.Server) error
	ready     chan struct{}
	server    *grpc.Server
	lis       net.Listener
}

// NewGRPCServerComponent creates a new server component with the given service registrars.
func NewGRPCServerComponent(log zerolog.Logger, cfg Config, registrar func(context.Context, *grpc.Server) error) *GRPCServerComponent {
	return &GRPCServerComponent{
		log:       log,
		cfg:       cfg,
		registrar: registrar,
		ready:     make(chan struct{}),
	}
}

// Start creates the listener, registers services, and begins serving.
// It blocks until the context is cancelled or the server fails.
func (s *GRPCServerComponent) Start(ctx context.Context) error {
	// 1. Create listener (non‑blocking, but we'll close it on cancellation).
	lis, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.Port, err)
	}
	s.lis = lis

	// 2. Ensure listener is closed if context is cancelled before we finish.
	go func() {
		<-ctx.Done()
		if err := lis.Close(); err != nil {
			panic(err)
		}
	}()

	s.log.Info().Str("port", s.cfg.Port).Msg("gRPC server listening")

	// 3. Create server with interceptors.
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.ReqIDServerInterceptor,
			LoggingUnaryServerInterceptor(s.log),
		))

	// 4. Run the (possibly blocking) registrar with context.
	if err := s.registrar(ctx, grpcServer); err != nil {
		// Registration failed – close listener and return error.
		lis.Close()
		return fmt.Errorf("registrar failed: %w", err)
	}

	s.server = grpcServer

	// 5. Channel for Serve errors.
	serveErr := make(chan error, 1)
	go func() {
		s.log.Debug().Str("port", s.cfg.Port).Msg("gRPC server starting")
		if err := grpcServer.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			serveErr <- err
		}
	}()

	// 6. Signal readiness.
	close(s.ready)

	// 7. Wait for shutdown or serve error.
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
	s.log.Debug().Msg("GRPCServerComponent.Stop: starting")
	if s.server == nil {
		s.log.Debug().Msg("GRPCServerComponent.Stop: server nil, returning")
		return nil
	}

	stopCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	s.log.Debug().Msg("GRPCServerComponent.Stop: calling GracefulStop")
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		s.log.Debug().Msg("GRPCServerComponent.Stop: GracefulStop completed")
		return nil
	case <-stopCtx.Done():
		s.log.Debug().Msg("GRPCServerComponent.Stop: timeout, forcing Stop")
		s.server.Stop()
		return fmt.Errorf("shutdown timed out")
	}
}

// Ready returns a channel that is closed when the connection is established.
func (s *GRPCServerComponent) Ready() <-chan struct{} {
	return s.ready
}
