package namer

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"userNamer",
	fx.Provide(
		NewUserNamer,
	),
)
