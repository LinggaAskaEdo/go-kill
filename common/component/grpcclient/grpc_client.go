package grpcclient

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Target   string        `yaml:"target"`
	Timeout  time.Duration `yaml:"timeout"`
	Insecure bool          `yaml:"insecure"`
}

type GRPCClientComponent struct {
	log   zerolog.Logger
	cfg   Config
	ready chan struct{}
	conn  *grpc.ClientConn
}

// NewGRPCClientComponent creates a new component.
func NewGRPCClientComponent(log zerolog.Logger, cfg Config) *GRPCClientComponent {
	return &GRPCClientComponent{
		log:   log,
		cfg:   cfg,
		ready: make(chan struct{}),
	}
}

// Start dials the target and blocks until the context is cancelled.
func (c *GRPCClientComponent) Start(ctx context.Context) error {
	opts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(c.ReqIDClientInterceptor),
	}

	if c.cfg.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// In production must configure TLS here.
		return fmt.Errorf("secure mode not implemented; set insecure: true for development")
	}

	client, err := grpc.NewClient(c.cfg.Target, opts...)
	if err != nil {
		return fmt.Errorf("dial %s: %w", c.cfg.Target, err)
	}

	c.conn = client

	waitCtx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	for {
		state := client.GetState()
		if state == connectivity.Ready {
			break
		}

		if !client.WaitForStateChange(waitCtx, state) {
			client.Close()
			return fmt.Errorf("timeout waiting for connection to %s", c.cfg.Target)
		}
	}

	close(c.ready) // signal readiness
	c.log.Debug().Str("target", c.cfg.Target).Msg("gRPC client connected")
	<-ctx.Done() // Block until shutdown signal
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

// Ready returns a channel that is closed when the connection is established.
func (c *GRPCClientComponent) Ready() <-chan struct{} {
	return c.ready
}
