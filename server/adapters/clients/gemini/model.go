package gemini

import (
	"github.com/jcfug8/daylear/server/ports/recipeocr"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gemini",
	fx.Provide(
		fx.Annotate(
			NewRecipeOCRClient,
			fx.As(new(recipeocr.Client)),
		),
	),
)
