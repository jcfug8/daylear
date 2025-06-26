package v1alpha1

import (
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/fieldmasker"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewUserServiceParams defines the dependencies for the UserService.
type NewUserServiceParams struct {
	fx.In

	Domain          domain.Domain
	Log             zerolog.Logger
	UserFieldMasker fieldMasker.UserFieldMasker
}

// NewUserService creates a new UserService.
func NewUserService(params NewUserServiceParams) (*UserService, error) {
	userNamer, err := namer.NewReflectNamer[*pb.User]()
	if err != nil {
		return nil, err
	}

	return &UserService{
		domain:          params.Domain,
		log:             params.Log,
		userFieldMasker: params.UserFieldMasker,
		userNamer:       userNamer,
	}, nil
}

// UserService defines the grpc handlers for the UserService.
type UserService struct {
	pb.UnimplementedUserServiceServer
	domain          domain.Domain
	log             zerolog.Logger
	userFieldMasker fieldMasker.UserFieldMasker
	userNamer       namer.ReflectNamer
}

// Register registers s to the grpc implementation of the service.
func (s *UserService) Register(server *grpc.Server) error {
	pb.RegisterUserServiceServer(server, s)
	return nil
}
