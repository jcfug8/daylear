package dialer

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Fx tags -
const (
	ConfigOptionTag = `group:"gormDialerConfigOptions"`
	OptionTag       = `group:"gormOptions"`
	PluginTag       = `group:"gormDialerPlugins"`
)

// Module -
var Module = fx.Module(
	"gormDialer",
	fx.Provide(
		NewConfig,

		fx.Annotate(
			NewDialer,
			fx.As(new(Dialer)),
		),

		Dial,
	),
)

// GormOption -
func GormOption(opt gorm.Option) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() gorm.Option { return opt },
			fx.ResultTags(OptionTag),
		),
	)
}

// ProvideConfigOptions -
func ProvideConfigOptions(opts ...ConfigOption) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() []ConfigOption { return opts },
			fx.ResultTags(`group:"gormDialerConfigOptions,flatten"`),
		),
	)
}
