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

type ServiceRegistrar func(*grpc.Server)

type GRPCServerComponent struct {
	log        zerolog.Logger
	cfg        Config
	registrars []ServiceRegistrar
	server     *grpc.Server
	lis        net.Listener
}

// NewGRPCServerComponent creates a new server component with the given service registrars.
func NewGRPCServerComponent(log zerolog.Logger, cfg Config, registrars ...ServiceRegistrar) *GRPCServerComponent {
	return &GRPCServerComponent{
		log:        log,
		cfg:        cfg,
		registrars: registrars,
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

	grpcServer := grpc.NewServer()
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
	if s.server == nil {
		return nil
	}

	// Use a separate context for the graceful stop deadline
	stopCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

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
		s.server.Stop() // force stop
		return fmt.Errorf("gRPC server shutdown timed out, forced stop")
	}
}
