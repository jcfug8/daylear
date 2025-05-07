package fieldbehaviorvalidator

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"fieldbehaviorvalidator",
	fx.Provide(
		NewCreateFieldValidator,
	),
)
