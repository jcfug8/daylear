package start

import (
	"flag"

	config "github.com/jcfug8/daylear/server/adapters/clients/config"
	"github.com/jcfug8/daylear/server/adapters/clients/gemini"
	gorm "github.com/jcfug8/daylear/server/adapters/clients/gorm"
	gorm_dialer "github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/dialer/dialects/postgres"
	"github.com/jcfug8/daylear/server/adapters/clients/http/fileretriever"
	"github.com/jcfug8/daylear/server/adapters/clients/imagemagick"
	tokenClient "github.com/jcfug8/daylear/server/adapters/clients/jwt/token"
	s3 "github.com/jcfug8/daylear/server/adapters/clients/s3"
	grpcCalendarsV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/calendars/calendar/v1alpha1"
	grpcCirclesV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1"
	grpcListsV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/lists/list/v1alpha1"
	grpcRecipesV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/meals/recipes/v1alpha1"
	grpcUsersV1alpha1 "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1"
	oauth2 "github.com/jcfug8/daylear/server/adapters/services/http/auth/oauth2"
	tokenService "github.com/jcfug8/daylear/server/adapters/services/http/auth/token"
	caldav "github.com/jcfug8/daylear/server/adapters/services/http/caldav"
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
		caldav.Module,
		files.Module,
		grpcgateway.Module,
		openapi.Module,
		// users
		grpcUsersV1alpha1.Module,
		// recipes
		grpcRecipesV1alpha1.Module,
		// circles
		grpcCirclesV1alpha1.Module,
		// calendars
		grpcCalendarsV1alpha1.Module,
		// lists
		grpcListsV1alpha1.Module,

		// driven/secondary adapters
		gorm.Module,
		gorm_dialer.Module,
		postgres.Module,
		tokenClient.Module,
		s3.Module,
		imagemagick.Module,
		fileretriever.Module,
		gemini.Module,

		// domain
		domain.Module,

		fx.WithLogger(logger.NewFxLogger),
	}, opts...)

	app := fx.New(opts...)
	app.Run()

	return app.Err()
}
