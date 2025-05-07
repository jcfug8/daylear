package config

import (
	"github.com/jcfug8/daylear/server/ports/config"
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"config",
	fx.Provide(
		fx.Annotate(
			NewConfig,
			fx.As(new(config.Client)),
		),
	),
)
