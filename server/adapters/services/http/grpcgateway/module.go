package grpcgateway

import (
	server "github.com/jcfug8/daylear/server/adapters/servers/http"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpcGateway",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
