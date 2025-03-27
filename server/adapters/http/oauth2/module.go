package oauth2

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/adapters/http/server"
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
