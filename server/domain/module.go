package domain

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"domain",
	fx.Provide(
		NewDomain,
	),
)
