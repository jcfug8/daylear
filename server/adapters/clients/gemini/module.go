package gemini

import (
	"github.com/jcfug8/daylear/server/ports/imagegenerator"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"gemini",
	fx.Provide(
		fx.Annotate(
			NewRecipeGeminiClient,
			fx.As(new(recipescraper.Client)),
			fx.As(new(imagegenerator.Client)),
		),
	),
)
