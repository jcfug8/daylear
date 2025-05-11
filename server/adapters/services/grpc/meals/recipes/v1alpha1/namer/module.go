package namer

import (
	"go.uber.org/fx"
)

// Module -
var Module = fx.Module(
	"recipeNamer",
	fx.Provide(
		NewRecipeNamer,
	),
)
