package database

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/linggaaskaedo/go-kill/common/preference"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Enabled         bool
	Driver          string
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         bool
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type DB struct {
	*sqlx.DB
}

func New(cfg Config) (*DB, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	driver, host, err := getURI(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("%s status: FAILED", strings.ToUpper(cfg.Driver)))
		return nil, err
	}

	db, err := sqlx.Connect(driver, host)
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("%s status: FAILED", strings.ToUpper(cfg.Driver)))
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return &DB{db}, nil
}

func getURI(cfg Config) (string, string, error) {
	switch cfg.Driver {
	case preference.POSTGRES:
		ssl := `disable`
		if cfg.SSLMode {
			ssl = `require`
		}

		return cfg.Driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, ssl), nil

	case preference.MYSQL:
		ssl := `false`
		if cfg.SSLMode {
			ssl = `true`
		}

		return cfg.Driver, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s&parseTime=%t", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, ssl, true), nil

	default:
		return "", "", errors.New("DB Driver is not supported ")
	}
}

func (d *DB) Close() {
	if err := d.DB.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database")
	}
}
