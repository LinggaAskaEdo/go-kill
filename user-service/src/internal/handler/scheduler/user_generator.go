package scheduler

import (
	"context"

	"github.com/linggaaskaedo/go-kill/common/component/scheduler"

	"github.com/rs/zerolog"
)

type UserGeneratorJob struct {
	log zerolog.Logger
	cfg scheduler.Config
}

func NewUserGeneratorJob(log zerolog.Logger, cfg scheduler.Config) *UserGeneratorJob {
	return &UserGeneratorJob{
		log: log,
		cfg: cfg,
	}
}

func (j *UserGeneratorJob) Name() string {
	return "user_generator_job"
}

func (j *UserGeneratorJob) Schedule() string {
	return j.cfg.Cron
}

func (j *UserGeneratorJob) Run(ctx context.Context) error {
	if !j.cfg.Enabled {
		zerolog.Ctx(ctx).Debug().Msg(j.cfg.Name + " is disabled")
		return nil
	}

	j.log.Info().
		Int("batch_size", j.cfg.BatchSize).
		Msg("Generating random users")

	return nil
}
