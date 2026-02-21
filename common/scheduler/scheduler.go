package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type Config struct {
	Enabled       bool                 `yaml:"enabled"`
	SchedulerJobs SchedulerJobsOptions `yaml:"jobs"`
}

type SchedulerJobsOptions struct {
	UserGeneratorJob UserGeneratorJobOptions `yaml:"user_generator"`
}

type UserGeneratorJobOptions struct {
	Enabled   bool   `yaml:"enabled"`
	Cron      string `yaml:"cron"`
	BatchSize int    `yaml:"batch_size"`
	MinAge    int    `yaml:"min_age"`
	MaxAge    int    `yaml:"max_age"`
}

type Job interface {
	Name() string
	Schedule() string
	Run(ctx context.Context) error
}

type SchedulerComponent struct {
	log  zerolog.Logger
	cron *cron.Cron
	jobs []Job
	mu   sync.RWMutex
}

// NewSchedulerComponent creates a new scheduler component, registers all provided jobs,
// and returns nil if the scheduler is disabled.
func NewSchedulerComponent(log zerolog.Logger, cfg Config, jobs []Job) *SchedulerComponent {
	if !cfg.Enabled {
		return nil
	}

	sc := &SchedulerComponent{
		log:  log,
		cron: cron.New(cron.WithSeconds()),
		jobs: make([]Job, 0),
	}

	// Register each job
	for _, job := range jobs {
		if err := sc.addJob(job); err != nil {
			log.Warn().Err(err).Str("job", job.Name()).Msg("Failed to register job, skipping")
		}
	}

	return sc
}

// addJob registers a job with the cron scheduler. It is called during construction.
func (sc *SchedulerComponent) addJob(job Job) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	_, err := sc.cron.AddFunc(job.Schedule(), func() {
		sc.log.Info().Str("job", job.Name()).Msg("Job started")

		ctx := context.Background() // You may want to use a derived context with timeout
		if err := job.Run(ctx); err != nil {
			sc.log.Error().Err(err).Str("job", job.Name()).Msg("Job execution failed")
			return
		}

		sc.log.Info().Str("job", job.Name()).Msg("Job completed successfully")
	})
	if err != nil {
		return fmt.Errorf("add job %s: %w", job.Name(), err)
	}

	sc.jobs = append(sc.jobs, job)
	sc.log.Info().Str("job", job.Name()).Str("schedule", job.Schedule()).Msg("Job registered")
	return nil
}

// Start begins the cron scheduler and blocks until the context is cancelled.
func (sc *SchedulerComponent) Start(ctx context.Context) error {
	sc.mu.RLock()
	sc.cron.Start()
	sc.mu.RUnlock()

	sc.log.Debug().Msgf("Scheduler started, jobs registered: %d", len(sc.jobs))

	// Block until shutdown signal
	<-ctx.Done()

	sc.log.Debug().Msg("Scheduler component context cancelled – stopping")
	return nil
}

// Stop gracefully shuts down the cron scheduler, waiting for running jobs to finish
// (up to the context timeout).
func (sc *SchedulerComponent) Stop(ctx context.Context) error {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// cron.Stop returns a channel that is closed when all jobs have finished.
	stopCtx := sc.cron.Stop()

	select {
	case <-stopCtx.Done():
		sc.log.Debug().Msg("Scheduler stopped")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("scheduler stop timed out")
	}
}

// ListJobs returns the names of all registered jobs.
func (sc *SchedulerComponent) ListJobs() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	jobNames := make([]string, len(sc.jobs))
	for i, job := range sc.jobs {
		jobNames[i] = job.Name()
	}

	return jobNames
}
