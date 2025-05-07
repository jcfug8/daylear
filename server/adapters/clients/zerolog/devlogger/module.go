package devlogger

import (
	"github.com/jcfug8/daylear/server/adapters/clients/zerolog/logger"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// Fx tags -
const (
	LevelTag = `name:"devloggerLevel"`
)

// Module -
var Module = fx.Module(
	"logger",
	fx.Provide(
		fx.Annotate(
			NewWriter,
			fx.ResultTags(logger.WriterTag),
		),

		fx.Annotate(
			New,
			fx.ResultTags(logger.PluginsTag),
		),
	),

	fx.Provide(
		fx.Annotate(
			func(lvl zerolog.Level) logger.Option {
				return logger.WithLevel(lvl)
			},
			fx.ParamTags(LevelTag),
			fx.ResultTags(logger.OptionsTag),
		),
	),

	fx.Supply(
		fx.Annotate(
			zerolog.TraceLevel,
			fx.ResultTags(LevelTag),
		),
	),
)
