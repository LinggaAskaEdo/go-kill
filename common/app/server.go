package app

import (
	"context"
	"time"

	httpserver "github.com/linggaaskaedo/go-kill/common/server"

	"github.com/rs/zerolog"
)

type HTTPServerComponent struct {
	server *httpserver.Server
	logger zerolog.Logger
}

func NewHTTPServerComponent(server *httpserver.Server, logger zerolog.Logger) *HTTPServerComponent {
	return &HTTPServerComponent{server: server, logger: logger}
}

func (h *HTTPServerComponent) Start(ctx context.Context) error {
	errCh := h.server.Start()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (h *HTTPServerComponent) Stop(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}

// Optional: override shutdown timeout
func (h *HTTPServerComponent) ShutdownTimeout() time.Duration {
	return 10 * time.Second // HTTP server may need more time
}
