package start

import (
	"flag"

	config "github.com/jcfug8/daylear/server/adapters/clients/config"
	"github.com/jcfug8/daylear/server/adapters/clients/gemini"
	gorm "github.com/jcfug8/daylear/server/adapters/clients/gorm"
	gorm_dialer "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer/dialects/postgres"
	"github.com/jcfug8/daylear/server/adapters/clients/http/fileretriever"
	"github.com/jcfug8/daylear/server/adapters/clients/http/recipescraper"
	tokenClient "github.com/jcfug8/daylear/server/adapters/clients/jwt/token"
	mimetype "github.com/jcfug8/daylear/server/adapters/clients/mimetype"
	s3 "github.com/jcfug8/daylear/server/adapters/clients/s3"
	grpcCirclesV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1"
	circlesV1alpha1Masker "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/fieldmasker"
	grpcRecipesV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1"
	recipesV1alpha1Masker "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1/fieldmasker"
	grpcUsersV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1"
	usersV1alpha1Masker "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/fieldmasker"
	oauth2 "github.com/jcfug8/daylear/server/adapters/services/http/auth/oauth2"
	tokenService "github.com/jcfug8/daylear/server/adapters/services/http/auth/token"
	files "github.com/jcfug8/daylear/server/adapters/services/http/files"
	grpcgateway "github.com/jcfug8/daylear/server/adapters/services/http/grpcgateway"
	openapi "github.com/jcfug8/daylear/server/adapters/services/http/openapi"
	domain "github.com/jcfug8/daylear/server/domain"
	"go.uber.org/fx"

	// IRIOMO:CUSTOM_CODE_SLOT_START startImports

	logger "github.com/jcfug8/daylear/server/adapters/clients/zerolog/logger"
	httpServer "github.com/jcfug8/daylear/server/adapters/servers/http"
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

		// driving/primary adapters
		oauth2.Module,
		tokenService.Module,
		files.Module,
		grpcgateway.Module,
		openapi.Module,
		// users
		grpcUsersV1alpha1.Module,
		usersV1alpha1Masker.Module,
		// recipes
		grpcRecipesV1alpha1.Module,
		recipesV1alpha1Masker.Module,
		// circles
		grpcCirclesV1alpha1.Module,
		circlesV1alpha1Masker.Module,

		// driven/secondary adapters
		gorm.Module,
		gorm_dialer.Module,
		postgres.Module,
		tokenClient.Module,
		s3.Module,
		mimetype.Module,
		fileretriever.Module,
		recipescraper.Module,
		gemini.Module,

		// domain
		domain.Module,

		fx.WithLogger(logger.NewFxLogger),
	}, opts...)

	app := fx.New(opts...)
	app.Run()

	return app.Err()
}
