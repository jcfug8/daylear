package caldav

import (
	"go.uber.org/fx"

	server "github.com/jcfug8/daylear/server/adapters/servers/http"
)

var Module = fx.Module(
	"caldavAdapter",
	fx.Provide(
		NewService,
	),

	server.ProvideAsService[*Service](),
)
