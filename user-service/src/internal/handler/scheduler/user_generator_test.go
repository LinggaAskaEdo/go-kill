package scheduler

import (
	"context"
	"testing"

	"github.com/linggaaskaedo/go-kill/common/component/scheduler"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const (
	testJobName      = "UserGenerator"
	testJobCron      = "0 * * * *"
	testJobBatchSize = 100
)

func TestName(t *testing.T) {
	cfg := scheduler.Config{
		Name:      testJobName,
		Enabled:   true,
		Cron:      testJobCron,
		BatchSize: testJobBatchSize,
	}

	job := NewUserGeneratorJob(zerolog.Logger{}, cfg)

	assert.Equal(t, "user_generator_job", job.Name())
}

func TestSchedule(t *testing.T) {
	cfg := scheduler.Config{
		Name:      testJobName,
		Enabled:   true,
		Cron:      testJobCron,
		BatchSize: testJobBatchSize,
	}

	job := NewUserGeneratorJob(zerolog.Logger{}, cfg)

	assert.Equal(t, testJobCron, job.Schedule())
}

func TestRunDisabled(t *testing.T) {
	cfg := scheduler.Config{
		Name:      testJobName,
		Enabled:   false,
		Cron:      testJobCron,
		BatchSize: testJobBatchSize,
	}

	job := NewUserGeneratorJob(zerolog.Logger{}, cfg)
	ctx := context.Background()

	err := job.Run(ctx)

	assert.NoError(t, err)
}

func TestRunEnabled(t *testing.T) {
	cfg := scheduler.Config{
		Name:      testJobName,
		Enabled:   true,
		Cron:      testJobCron,
		BatchSize: testJobBatchSize,
	}

	job := NewUserGeneratorJob(zerolog.Logger{}, cfg)
	ctx := context.Background()

	err := job.Run(ctx)

	assert.NoError(t, err)
}
