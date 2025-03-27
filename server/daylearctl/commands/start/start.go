package start

import (
	"flag"

	"github.com/jcfug8/daylear/server/adapters/config"
	"github.com/jcfug8/daylear/server/adapters/gorm"
	gorm_dialer "github.com/jcfug8/daylear/server/adapters/gorm/dialer"
	"github.com/jcfug8/daylear/server/adapters/gorm/dialer/dialects/postgres"
	"github.com/jcfug8/daylear/server/adapters/grpc/fieldbehaviorvalidator"
	grpcRecipesV1alpha1 "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1"
	recipesV1alpha1Masker "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/fieldmasker"
	recipesV1alpha1Namer "github.com/jcfug8/daylear/server/adapters/grpc/meals/recipes/v1alpha1/namer"
	grpcUsersV1alpha1 "github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1"
	usersV1alpha1Masker "github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1/fieldmasker"
	usersV1alpha1Namer "github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1/namer"
	"github.com/jcfug8/daylear/server/adapters/http/auth"
	"github.com/jcfug8/daylear/server/adapters/http/files"
	"github.com/jcfug8/daylear/server/adapters/http/grpcgateway"
	"github.com/jcfug8/daylear/server/adapters/http/oauth2"
	"github.com/jcfug8/daylear/server/adapters/http/ping"
	"github.com/jcfug8/daylear/server/adapters/jwt/token"
	"github.com/jcfug8/daylear/server/adapters/mimetype"
	"github.com/jcfug8/daylear/server/adapters/s3"

	domain "github.com/jcfug8/daylear/server/domain"
	"go.uber.org/fx"

	// IRIOMO:CUSTOM_CODE_SLOT_START startImports

	httpServer "github.com/jcfug8/daylear/server/adapters/http/server"
	logger "github.com/jcfug8/daylear/server/adapters/zerolog/logger"
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

const (
	// StartCmd defines the command for starting the service.
	StartCmd = "start"
)

var logLevel string

// Start starts the service
func Start() (err error) {
	flag.StringVar(&logLevel, "log", "info", "verbose output")

	flag.Parse()

	return start()
}

func start(opts ...fx.Option) error {
	opts = append([]fx.Option{
		httpServer.Module,
		logger.Module,
		config.Module,
		fieldbehaviorvalidator.Module,

		// driving/primary adapters
		ping.Module,
		oauth2.Module,
		auth.Module,
		grpcgateway.Module,
		files.Module,
		// users
		grpcUsersV1alpha1.Module,
		usersV1alpha1Namer.Module,
		usersV1alpha1Masker.Module,
		// recipes
		grpcRecipesV1alpha1.Module,
		recipesV1alpha1Namer.Module,
		recipesV1alpha1Masker.Module,

		// driven/secondary adapters
		gorm.Module,
		gorm_dialer.Module,
		postgres.Module,
		token.Module,
		s3.Module,
		mimetype.Module,

		// domain
		domain.Module,

		fx.WithLogger(logger.NewFxLogger),
	}, opts...)

	app := fx.New(opts...)
	app.Run()

	return app.Err()
}
