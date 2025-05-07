package files

import (
	"go.uber.org/fx"

	server "github.com/jcfug8/daylear/server/adapters/servers/http"
)

var Module = fx.Module(
	"imagesAdapter",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
