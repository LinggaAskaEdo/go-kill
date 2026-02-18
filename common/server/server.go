package server

import (
	"context"
	"net/http"
	"time"

	"github.com/linggaaskaedo/go-kill/common/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Server struct {
	engine *gin.Engine
	srv    *http.Server
}

func New(cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(middleware.RequestIDMiddleware(), middleware.LoggingMiddleware(), middleware.RecoveryMiddleware())

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{engine: engine, srv: srv}
}

func (s *Server) Engine() *gin.Engine { return s.engine }

func (s *Server) Start() <-chan error {
	errCh := make(chan error, 1)
	go func() {
		log.Info().Str("port", s.srv.Addr).Msg("HTTP server starting")
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()
	
	return errCh
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server")
	return s.srv.Shutdown(ctx)
}
