package app

import (
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

// func (d *DatabaseComponent) Start(ctx context.Context) error {
// 	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()
// 	if err := d.db.PingContext(pingCtx); err != nil {
// 		return err
// 	}

// 	d.logger.Info().Msg("Database ready")
// 	<-ctx.Done()

// 	return nil
// }

// func (d *DatabaseComponent) Stop(ctx context.Context) error {
// 	d.logger.Info().Msg("Closing database")
// 	d.db.Close()

// 	return nil
// }
