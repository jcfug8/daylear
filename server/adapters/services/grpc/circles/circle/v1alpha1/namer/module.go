package namer

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"grpcCircleNamer",
	fx.Provide(
		NewCircleNamer,
	),
)
