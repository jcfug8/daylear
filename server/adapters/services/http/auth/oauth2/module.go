package oauth2

import (
	"go.uber.org/fx"

	server "github.com/jcfug8/daylear/server/adapters/servers/http"
)

var Module = fx.Module(
	"oauth2",
	fx.Provide(
		// NewAmazonService,
		// NewFacebookService,
		NewGoogleService,
	),

	// server.ProvideAsService[*AmazonService](),
	// server.ProvideAsService[*FacebookService](),
	server.ProvideAsService[*GoogleService](),
)
