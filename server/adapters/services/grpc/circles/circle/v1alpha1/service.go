package v1alpha1

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewCircleServiceParams defines the dependencies for the CircleService.
type NewCircleServiceParams struct {
	fx.In

	Domain            domain.Domain
	Log               zerolog.Logger
	CircleFieldMasker fieldmask.FieldMasker `name:"v1alpha1CircleFieldMasker"`
	AccessFieldMasker fieldmask.FieldMasker `name:"v1alpha1CircleAccessFieldMasker"`
	CircleNamer       namer.ReflectNamer    `name:"v1alpha1CircleNamer"`
	AccessNamer       namer.ReflectNamer    `name:"v1alpha1CircleAccessNamer"`
	UserNamer         namer.ReflectNamer    `name:"v1alpha1UserNamer"`
}

// NewCircleService creates a new CircleService.
func NewCircleService(params NewCircleServiceParams) (*CircleService, error) {
	return &CircleService{
		domain:            params.Domain,
		log:               params.Log,
		circleFieldMasker: params.CircleFieldMasker,
		accessFieldMasker: params.AccessFieldMasker,
		circleNamer:       params.CircleNamer,
		accessNamer:       params.AccessNamer,
		userNamer:         params.UserNamer,
	}, nil
}

// CircleService defines the grpc handlers for the CircleService.
type CircleService struct {
	pb.UnimplementedCircleServiceServer
	pb.UnimplementedCircleAccessServiceServer
	domain            domain.Domain
	log               zerolog.Logger
	circleFieldMasker fieldmask.FieldMasker
	accessFieldMasker fieldmask.FieldMasker
	circleNamer       namer.ReflectNamer
	accessNamer       namer.ReflectNamer
	userNamer         namer.ReflectNamer
}

// Register registers s to the grpc implementation of the service.
func (s *CircleService) Register(server *grpc.Server) error {
	pb.RegisterCircleServiceServer(server, s)
	pb.RegisterCircleAccessServiceServer(server, s)
	return nil
}
