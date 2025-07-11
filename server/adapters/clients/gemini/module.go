package gemini

import (
	"github.com/jcfug8/daylear/server/ports/imagegenerator"
	"github.com/jcfug8/daylear/server/ports/ingredientcleaner"
	"github.com/jcfug8/daylear/server/ports/recipeocr"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gemini",
	fx.Provide(
		fx.Annotate(
			NewRecipeGeminiClient,
			fx.As(new(recipeocr.Client)),
			fx.As(new(ingredientcleaner.Client)),
			fx.As(new(imagegenerator.Client)),
		),
	),
)
