package logger

import (
	"go.uber.org/fx"
)

// Fx tags -
const (
	HooksTag   = `group:"loggerHooks"`
	OptionsTag = `group:"loggerOptions"`
	PluginsTag = `group:"loggerPlugins"`
	WriterTag  = `name:"loggerWriter"`
)

// Module -
var Module = fx.Module(
	"logger",
	fx.Provide(
		NewFlagParser,
		NewLogger,
		NewOffset,
		NewOutput,
	),
)

// LevelFlag -
func LevelFlag(flag string) fx.Option {
	return fx.Provide(
		func() *Flags { return NewFlags(flag) },
	)
}

// ProvideOptions -
func ProvideOptions(opts ...Option) fx.Option {
	return fx.Provide(
		fx.Annotate(
			func() []Option { return opts },
			fx.ResultTags(`group:"loggerOptions,flatten"`),
		),
	)
}
