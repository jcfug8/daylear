package s3

import (
	"github.com/jcfug8/daylear/server/ports/filestorage"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"s3",
	fx.Provide(
		fx.Annotate(
			NewClient,
			fx.As(new(filestorage.Client)),
		),
	),
)
