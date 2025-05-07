package postgres

import (
	gormdialer "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"

	"go.uber.org/fx"
)

// Fx tags -
const (
	ConfigsTag  = `group:"postgresConfigs"`
	DialectsTag = `group:"postgresDialects"`
)

// Module -
var Module = fx.Module(
	"postgres",
	fx.Provide(
		fx.Annotate(
			NewDialect,
			fx.As(new(gormdialer.Dialect)),
		),

		// func(config config.Config) (*config.PostgresConfigerror) {
		// 	return config.GetPostgresConfig("default")
		// },

		fx.Annotate(
			NewDialects,
			// fx.ParamTags(`optional:"true"`, ConfigsTag),
			fx.ResultTags(DialectsTag),
		),
	),
)
