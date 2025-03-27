package ping

import (
	"github.com/jcfug8/daylear/server/adapters/http/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"http-ping-service",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
