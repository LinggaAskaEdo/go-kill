package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/linggaaskaedo/go-kill/common/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Config struct {
	AppName         string        `yaml:"app_name"`
	Port            int           `yaml:"port"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HTTPServerComponent struct {
	log        zerolog.Logger
	cfg        Config
	middleware middleware.Middleware
	engine     *gin.Engine
	httpServer *http.Server
}

// NewHTTPServerComponent creates a new HTTP server component.
// The engine and server are created during Start.
func NewHTTPServerComponent(log zerolog.Logger, cfg Config, mw middleware.Middleware, gin *gin.Engine) *HTTPServerComponent {
	return &HTTPServerComponent{
		log:        log,
		cfg:        cfg,
		middleware: mw,
		engine:     gin,
	}
}

// Start builds the Gin engine, applies middleware, and begins listening.
// It blocks until ctx is done or the server fails to start.
func (h *HTTPServerComponent) Start(ctx context.Context) error {
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
	if h.httpServer == nil {
		return nil
	}

	// Use the configured shutdown timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), h.cfg.ShutdownTimeout)
	defer cancel()

	h.log.Debug().Msg("Shutting down HTTP server gracefully")
	if err := h.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server shutdown error: %w", err)
	}

	h.log.Debug().Msg("HTTP server stopped")
	return nil
}

// Engine returns the Gin engine (useful for tests or if other components need to add routes after start).
func (h *HTTPServerComponent) Engine() *gin.Engine {
	return h.engine
}
