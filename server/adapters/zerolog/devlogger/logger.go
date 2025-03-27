package devlogger

import (
	"io"
	"os"
	"time"

	"github.com/jcfug8/daylear/server/adapters/zerolog/logger"

	"github.com/rs/zerolog"
)

// NewWriter -
func NewWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
}

// New creates a new logger Plugin
func New() logger.Plugin {
	return &plugin{}
}

type plugin struct{}

// Logger -
func (*plugin) Logger(log zerolog.Logger) zerolog.Logger {
	return log.With().Caller().Logger()
}
