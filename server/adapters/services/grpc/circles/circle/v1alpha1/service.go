package v1alpha1

import (
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/fieldmasker"
	namer "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/namer"
	fieldValidator "github.com/jcfug8/daylear/server/adapters/services/grpc/fieldbehaviorvalidator"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewCircleServiceParams defines the dependencies for the CircleService.
type NewCircleServiceParams struct {
	fx.In

	Domain                  domain.Domain
	FieldValidator          fieldValidator.FieldBehaviorValidator
	Log                     zerolog.Logger
	CircleFieldMasker       fieldMasker.CircleFieldMasker
	CircleNamer             namer.CircleNamer
	PublicCircleFieldMasker fieldMasker.PublicCircleFieldMasker
	PublicCircleNamer       namer.PublicCircleNamer
}

// NewCircleService creates a new CircleService.
func NewCircleService(params NewCircleServiceParams) *CircleService {
	return &CircleService{
		domain:                  params.Domain,
		fieldBehaviorValidator:  params.FieldValidator,
		log:                     params.Log,
		circleFieldMasker:       params.CircleFieldMasker,
		circleNamer:             params.CircleNamer,
		publicCircleFieldMasker: params.PublicCircleFieldMasker,
		publicCircleNamer:       params.PublicCircleNamer,
	}
}

// CircleService defines the grpc handlers for the CircleService.
type CircleService struct {
	pb.UnimplementedCircleServiceServer
	pb.UnimplementedPublicCircleServiceServer
	domain                  domain.Domain
	fieldBehaviorValidator  fieldValidator.FieldBehaviorValidator
	log                     zerolog.Logger
	circleFieldMasker       fieldMasker.CircleFieldMasker
	circleNamer             namer.CircleNamer
	publicCircleFieldMasker fieldMasker.PublicCircleFieldMasker
	publicCircleNamer       namer.PublicCircleNamer
}

// Register registers s to the grpc implementation of the service.
func (s *CircleService) Register(server *grpc.Server) error {
	pb.RegisterCircleServiceServer(server, s)
	pb.RegisterPublicCircleServiceServer(server, s)
	return nil
}
