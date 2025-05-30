package token

import (
	"go.uber.org/fx"

	server "github.com/jcfug8/daylear/server/adapters/servers/http"
)

var Module = fx.Module(
	"tokenAuthAdapter",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
