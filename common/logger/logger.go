package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Level       string
	PrettyPrint bool
	Output      io.Writer // for testing
}

func New(cfg Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Nop(), err
	}
	zerolog.SetGlobalLevel(level) // still set global level for libraries

	out := cfg.Output
	if out == nil {
		out = os.Stderr
	}

	var logger zerolog.Logger
	if cfg.PrettyPrint {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.New(out)
	}

	logger = logger.With().Timestamp().Logger()

	return logger, nil
}
