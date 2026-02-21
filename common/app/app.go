package app

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type componentWrapper struct {
	Component
	shutdownTimeout time.Duration
}

type App struct {
	components    []componentWrapper
	globalTimeout time.Duration
	logger        zerolog.Logger
}

type Option func(*App)

func WithShutdownTimeout(timeout time.Duration) Option {
	return func(a *App) { a.globalTimeout = timeout }
}

func WithLogger(logger zerolog.Logger) Option {
	return func(a *App) { a.logger = logger }
}

func New(opts ...Option) *App {
	a := &App{
		globalTimeout: 30 * time.Second,
		logger:        zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// Add registers components with optional per-component shutdown timeout.
func (a *App) Add(comp Component, opts ...any) {
	wrapper := componentWrapper{Component: comp}
	// Options can be used to set per-component timeout, etc.
	for _, opt := range opts {
		if to, ok := opt.(time.Duration); ok {
			wrapper.shutdownTimeout = to
		}
	}
	a.components = append(a.components, wrapper)
}

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	a.startComponents(g, ctx)

	if err := a.waitForShutdown(ctx, g); err != nil {
		a.logger.Error().Err(err).Msg("Component error before shutdown")
	}

	a.stopComponents()

	// Wait for all components to fully exit (they should have due to ctx cancellation)
	if err := g.Wait(); err != nil && err != context.Canceled {
		return err
	}

	a.logger.Info().Msg("Application stopped")

	return nil
}

// startComponents launches all component goroutines.
func (a *App) startComponents(g *errgroup.Group, ctx context.Context) {
	for _, comp := range a.components {
		comp := comp // capture loop variable
		g.Go(func() error {
			a.logger.Info().Type("component", comp.Component).Msg("Starting component")
			err := comp.Start(ctx)
			if err != nil {
				a.logger.Error().Err(err).Type("component", comp.Component).Msg("Component failed")
			}

			return err
		})
	}
}

// waitForShutdown blocks until either a shutdown signal is received or a component error occurs.
func (a *App) waitForShutdown(ctx context.Context, g *errgroup.Group) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- g.Wait()
	}()

	select {
	case <-ctx.Done():
		a.logger.Info().Msg("Shutdown signal received")
		return nil
	case err := <-errCh:
		return err
	}
}

// stopComponents shuts down all components in reverse order with appropriate timeouts.
func (a *App) stopComponents() {
	for i := len(a.components) - 1; i >= 0; i-- {
		comp := a.components[i]
		timeout := a.getComponentTimeout(comp)
		stopCtx, cancel := context.WithTimeout(context.Background(), timeout)
		a.logger.Info().Type("component", comp.Component).Dur("timeout", timeout).Msg("Stopping component")
		if err := comp.Stop(stopCtx); err != nil {
			a.logger.Error().Err(err).Type("component", comp.Component).Msg("Stop error")
		}

		cancel()
	}
}

// getComponentTimeout determines the shutdown timeout for a component.
func (a *App) getComponentTimeout(comp componentWrapper) time.Duration {
	if comp.shutdownTimeout != 0 {
		return comp.shutdownTimeout
	}

	if c, ok := comp.Component.(ComponentWithShutdownTimeout); ok {
		return c.ShutdownTimeout()
	}

	return a.globalTimeout
}
