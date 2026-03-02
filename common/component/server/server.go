package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/linggaaskaedo/go-kill/common/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Engine = gin.Engine

type Config struct {
	AppName         string        `yaml:"app_name"`
	Port            int           `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HTTPServerComponent struct {
	log            zerolog.Logger
	cfg            Config
	middleware     middleware.Middleware
	engine         *gin.Engine
	routeRegistrar func(*Engine)
	ready          chan struct{}
	httpServer     *http.Server
}

// NewHTTPServerComponent creates a new HTTP server component.
// The engine and server are created during Start.
func NewHTTPServerComponent(log zerolog.Logger, cfg Config, mw middleware.Middleware, gin *gin.Engine, routeRegistrar func(*Engine)) *HTTPServerComponent {
	return &HTTPServerComponent{
		log:            log,
		cfg:            cfg,
		middleware:     mw,
		engine:         gin,
		routeRegistrar: routeRegistrar,
		ready:          make(chan struct{}),
	}
}

// Start builds the Gin engine, applies middleware, and begins listening.
// It blocks until ctx is done or the server fails to start.
func (h *HTTPServerComponent) Start(ctx context.Context) error {
	// Register service‑specific routes
	if h.routeRegistrar != nil {
		h.routeRegistrar(h.engine)
	}

	// Create HTTP server
	addr := fmt.Sprintf(":%d", h.cfg.Port)
	h.httpServer = &http.Server{
		Addr:         addr,
		Handler:      h.engine,
		ReadTimeout:  h.cfg.ReadTimeout,
		WriteTimeout: h.cfg.WriteTimeout,
		IdleTimeout:  h.cfg.IdleTimeout,
	}

	// Channel to capture ListenAndServe errors
	serveErr := make(chan error, 1)
	go func() {
		h.log.Debug().Str("addr", addr).Msg("HTTP server starting")
		if err := h.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serveErr <- err
		}
	}()

	close(h.ready) // signal readiness
	h.log.Debug().Msg("HTTP server started")

	// Wait for shutdown signal or startup error
	select {
	case <-ctx.Done():
		h.log.Debug().Msg("HTTP server context cancelled – stopping")
		return nil
	case err := <-serveErr:
		return fmt.Errorf("HTTP server listen error: %w", err)
	}
}

// Stop gracefully shuts down the server with a timeout.
func (h *HTTPServerComponent) Stop(ctx context.Context) error {
	h.log.Debug().Msg("HTTPServerComponent.Stop: starting")
	if h.httpServer == nil {
		h.log.Debug().Msg("HTTPServerComponent.Stop: server nil, returning")
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), h.cfg.ShutdownTimeout)
	defer cancel()

	h.log.Debug().Msg("HTTPServerComponent.Stop: calling Shutdown")
	if err := h.httpServer.Shutdown(shutdownCtx); err != nil {
		h.log.Error().Err(err).Msg("HTTPServerComponent.Stop: Shutdown error")
		return err
	}

	h.log.Debug().Msg("HTTPServerComponent.Stop: Shutdown completed successfully")
	return nil
}

// Engine returns the Gin engine (useful for tests or if other components need to add routes after start).
func (h *HTTPServerComponent) Engine() *gin.Engine {
	return h.engine
}

// Ready returns a channel that is closed when the connection is established.
func (h *HTTPServerComponent) Ready() <-chan struct{} {
	return h.ready
}
