package namer

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"domain",
	fx.Provide(
		NewUserNamer,
	),
)
