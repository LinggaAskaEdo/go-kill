package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Config struct {
	Enabled         bool          `yaml:"enabled"`
	Network         string        `yaml:"network"`
	Address         string        `yaml:"address"`
	Password        string        `yaml:"password"`
	CacheTTL        time.Duration `yaml:"cache_ttl"`
	MaxRetries      int           `yaml:"max_retries"`
	MinRetryBackoff time.Duration `yaml:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `yaml:"max_retry_backoff"`
	DialTimeout     time.Duration `yaml:"dial_timeout"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	PoolSize        int           `yaml:"pool_size"`
	MinIdleConns    int           `yaml:"min_idle_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxActiveConns  int           `yaml:"max_active_conns"`
	PoolTimeout     time.Duration `yaml:"pool_timeout"`
}

type RedisComponent struct {
	log    zerolog.Logger
	cfg    Config
	db     int
	client *redis.Client
}

func NewRedisComponent(log zerolog.Logger, cfg Config, redisType string) *RedisComponent {
	var db int

	switch redisType {
	case "apps":
		db = 0
	case "auth":
		db = 11
	default:
		db = 13
	}

	return &RedisComponent{
		log: log,
		cfg: cfg,
		db:  db,
	}
}

// Start initialises the Redis client, verifies the connection with a Ping,
// and then blocks until the context is cancelled.
// It returns an error if the client cannot be created or the ping fails.
func (r *RedisComponent) Start(ctx context.Context) error {
	// Create the client
	r.client = redis.NewClient(&redis.Options{
		Network:         r.cfg.Network,
		Addr:            r.cfg.Address,
		Password:        r.cfg.Password,
		DB:              r.db,
		MaxRetries:      r.cfg.MaxRetries,
		MinRetryBackoff: r.cfg.MinRetryBackoff,
		MaxRetryBackoff: r.cfg.MaxRetryBackoff,
		DialTimeout:     r.cfg.DialTimeout,
		ReadTimeout:     r.cfg.ReadTimeout,
		WriteTimeout:    r.cfg.WriteTimeout,
		PoolSize:        r.cfg.PoolSize,
		MinIdleConns:    r.cfg.MinIdleConns,
		MaxIdleConns:    r.cfg.MaxIdleConns,
		MaxActiveConns:  r.cfg.MaxActiveConns,
		PoolTimeout:     r.cfg.PoolTimeout,
	})

	// Verify connectivity (Ping uses the provided context)
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}

	r.log.Debug().Msg("Redis component started and ping successful")

	// Block until the context is cancelled.
	// This keeps the component "alive" from the errgroup's perspective.
	<-ctx.Done()

	r.log.Debug().Msg("Redis component context cancelled – stopping")
	return nil
}

// Stop performs final cleanup. It closes the Redis client.
// It is called after Start has returned (due to context cancellation).
func (r *RedisComponent) Stop(ctx context.Context) error {
	if r.client == nil {
		return nil
	}

	// Close the client – it will wait for pending commands to finish.
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("redis close error: %w", err)
	}

	r.log.Debug().Msg("Redis component stopped")

	return nil
}

// Client returns the underlying Redis client for use by other components.
// It is safe to call only after Start has completed successfully.
func (r *RedisComponent) Client() *redis.Client {
	return r.client
}
