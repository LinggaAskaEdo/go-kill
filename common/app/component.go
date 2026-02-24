package app

import (
	"context"
	"time"
)

type Component interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Ready() <-chan struct{}
}

type ComponentWithShutdownTimeout interface {
	Component
	ShutdownTimeout() time.Duration
}
