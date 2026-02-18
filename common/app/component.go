package app

import (
	"context"
	"time"
)

type Component interface {
	// Start begins the component. It should block until the component is ready
	// and running, and return nil when stopped normally (ctx cancelled) or error.
	Start(ctx context.Context) error

	// Stop initiates graceful shutdown. The context has a timeout.
	Stop(ctx context.Context) error
}

// ComponentWithShutdownTimeout can be implemented to override global timeout.
type ComponentWithShutdownTimeout interface {
	Component
	ShutdownTimeout() time.Duration
}
