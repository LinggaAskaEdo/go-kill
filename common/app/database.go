package app

import (
	"context"
	"time"

	"github.com/linggaaskaedo/go-kill/common/database"

	"github.com/rs/zerolog"
)

type DatabaseComponent struct {
	db     *database.DB
	logger zerolog.Logger
}

func NewDatabaseComponent(db *database.DB, logger zerolog.Logger) *DatabaseComponent {
	return &DatabaseComponent{db: db, logger: logger}
}

func (d *DatabaseComponent) Start(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := d.db.PingContext(pingCtx); err != nil {
		return err
	}

	d.logger.Info().Msg("Database ready")
	<-ctx.Done()

	return nil
}

func (d *DatabaseComponent) Stop(ctx context.Context) error {
	d.logger.Info().Msg("Closing database")
	d.db.Close()

	return nil
}
