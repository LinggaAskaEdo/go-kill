package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/linggaaskaedo/go-kill/common/preference"
	"github.com/rs/zerolog"
)

type Config struct {
	Enabled         bool          `yaml:"enabled"`
	Driver          string        `yaml:"driver"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         bool          `yaml:"sslmode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

type DatabaseComponent struct {
	log zerolog.Logger
	cfg Config
	db  *sqlx.DB
}

// NewDatabaseComponent creates a new database component but does not start it.
// If the component is disabled, it returns nil (so you can skip adding it to the slice).
func NewDatabaseComponent(log zerolog.Logger, cfg Config) *DatabaseComponent {
	if !cfg.Enabled {
		return nil
	}

	return &DatabaseComponent{
		log: log,
		cfg: cfg,
	}
}

// Start establishes the database connection, verifies it with a ping,
// and then blocks until the context is cancelled.
// It returns an error if the connection cannot be established.
func (d *DatabaseComponent) Start(ctx context.Context) error {
	driver, uri, err := getURI(d.cfg)
	if err != nil {
		return fmt.Errorf("build database URI: %w", err)
	}

	db, err := sqlx.ConnectContext(ctx, driver, uri)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	d.db = db

	// Configure connection pool
	d.db.SetMaxOpenConns(d.cfg.MaxOpenConns)
	d.db.SetMaxIdleConns(d.cfg.MaxIdleConns)
	d.db.SetConnMaxLifetime(d.cfg.ConnMaxLifetime)
	d.db.SetConnMaxIdleTime(d.cfg.ConnMaxIdleTime)

	d.log.Debug().Msgf("%s database connected and ping OK", strings.ToUpper(d.cfg.Driver))

	// Block until shutdown signal
	<-ctx.Done()

	d.log.Debug().Msgf("%s database context cancelled – stopping", strings.ToUpper(d.cfg.Driver))
	return nil
}

// Stop closes the database connection pool.
// It is called after Start has returned.
func (d *DatabaseComponent) Stop(ctx context.Context) error {
	if d.db == nil {
		return nil
	}

	// Close waits for all connections to be returned to the pool before closing.
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	
	d.log.Debug().Msgf("%s database stopped", strings.ToUpper(d.cfg.Driver))
	return nil
}

// Client returns the underlying *sqlx.DB for use by other components.
// It is safe to call only after Start has completed successfully.
func (d *DatabaseComponent) Client() *sqlx.DB {
	return d.db
}

// getURI constructs the driver name and connection string based on config.
// It is a slightly modified version of your existing getURI.
func getURI(cfg Config) (string, string, error) {
	switch cfg.Driver {
	case preference.POSTGRES:
		ssl := "disable"
		if cfg.SSLMode {
			ssl = "require"
		}

		uri := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, ssl)
		return cfg.Driver, uri, nil

	case preference.MYSQL:
		tls := "false"
		if cfg.SSLMode {
			tls = "true"
		}

		uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=%s&parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, tls)
		return cfg.Driver, uri, nil

	default:
		return "", "", errors.New("unsupported database driver: " + cfg.Driver)
	}
}
