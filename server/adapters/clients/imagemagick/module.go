package imagemagick

import (
	"github.com/jcfug8/daylear/server/ports/image"
	"go.uber.org/fx"
)

var Module = fx.Module("imagemagick",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(image.Client)),
		),
	),
)
