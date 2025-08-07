package openapi

import (
	"go.uber.org/fx"

	server "github.com/jcfug8/daylear/server/adapters/servers/http"
)

var Module = fx.Module(
	"openapiAdapter",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
