package recipescraper

import (
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"recipescraper",
	fx.Provide(
		fx.Annotate(
			NewDefaultClient,
			fx.As(new(recipescraper.DefaultClient)),
		),
	),
)
