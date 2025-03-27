package fileretriever

import (
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"fileretriever",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(fileretriever.Client)),
		),
	),
)
