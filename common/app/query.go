package app

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/query"
	"github.com/rs/zerolog"
)

type QueryComponent struct {
	query  *query.QueryLoader
	logger zerolog.Logger
}

func NewQueryComponent(query *query.QueryLoader, logger zerolog.Logger) *QueryComponent {
	return &QueryComponent{query: query, logger: logger}
}

func (d *QueryComponent) Start(ctx context.Context) error {
	if err := d.query.Load(); err != nil {
		d.logger.Panic().Err(err).Msg("Failed to load queries")
		return err
	}

	return nil
}

func (d *QueryComponent) Stop(ctx context.Context) error {
	d.logger.Info().Msg("Closing query loader")
	if err := d.query.Clear(); err != nil {
		d.logger.Error().Err(err).Msg("Failed to clear queries")
		return err
	}

	return nil
}
