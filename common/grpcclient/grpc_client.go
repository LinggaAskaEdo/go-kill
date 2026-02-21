package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Target   string        `yaml:"target"`
	Timeout  time.Duration `yaml:"timeout"`
	Insecure bool          `yaml:"insecure"`
}

type GRPCClientComponent struct {
	log  zerolog.Logger
	cfg  Config
	conn *grpc.ClientConn
}

// NewGRPCClientComponent creates a new component.
func NewGRPCClientComponent(log zerolog.Logger, cfg Config) *GRPCClientComponent {
	return &GRPCClientComponent{
		log: log,
		cfg: cfg,
	}
}

// Start dials the target and blocks until the context is cancelled.
func (c *GRPCClientComponent) Start(ctx context.Context) error {
	opts := []grpc.DialOption{
		grpc.WithBlock(), // wait for connection to be ready (remove if you prefer async)
	}

	if c.cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// In production you would configure TLS here.
		return fmt.Errorf("secure mode not implemented; set insecure: true for development")
	}

	dialCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, c.cfg.Target, opts...)
	if err != nil {
		return fmt.Errorf("dial %s: %w", c.cfg.Target, err)
	}
	c.conn = conn

	c.log.Debug().Str("target", c.cfg.Target).Msg("gRPC client connected")

	// Wait for shutdown signal
	<-ctx.Done()

	c.log.Debug().Msg("gRPC client context cancelled – stopping")
	return nil
}

// Stop closes the connection.
func (c *GRPCClientComponent) Stop(ctx context.Context) error {
	if c.conn == nil {
		return nil
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close gRPC client connection: %w", err)
	}

	c.log.Debug().Msg("gRPC client stopped")
	return nil
}

// Conn returns the underlying *grpc.ClientConn. Use it to create stubs.
func (c *GRPCClientComponent) Conn() *grpc.ClientConn {
	return c.conn
}
