package logger

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// LoggerParams -
type LoggerParams struct {
	fx.In

	Hooks   []zerolog.Hook `group:"loggerHooks"`
	Options []Option       `group:"loggerOptions"`
	Output  *Output
	Parser  *FlagParser `optional:"true"`
	Plugins []Plugin    `group:"loggerPlugins"`
}

// NewLogger -
func NewLogger(params LoggerParams) zerolog.Logger {
	opts := NewOptions(params.Options...)

	log := zerolog.New(params.Output).
		Level(zerolog.TraceLevel).
		With().Caller().Timestamp().Logger()

	if len(params.Hooks) > 0 {
		log = log.Hook(params.Hooks...)
	}

	log = ApplyPlugins(log, params.Plugins)

	lvl := params.Parser.Level()
	if opts.level != nil {
		lvl = *opts.level
	}

	log.Info().Msgf("setting log level to %s", lvl.String())
	params.Output.offset.Level(lvl)

	// Relay USR1 and USR2 signals to the offset
	signals := make(chan os.Signal, 100)
	signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2)
	go params.Output.offset.Signal(signals)

	log.Info().Msgf("initialized logger with %d hooks, %d plugins, and %d options at level %s", len(params.Hooks), len(params.Plugins), len(params.Options), params.Output.offset.GetLevel().String())

	return log
}
