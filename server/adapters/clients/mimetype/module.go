package mimetype

import (
	"github.com/jcfug8/daylear/server/ports/fileinspector"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"mimetypeInspector",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(fileinspector.Client)),
		),
	),
)
