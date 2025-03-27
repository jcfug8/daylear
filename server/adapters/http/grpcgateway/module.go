package grpcgateway

import (
	"github.com/jcfug8/daylear/server/adapters/http/server"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"userGrpcAdapter",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
