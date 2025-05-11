package v1alpha1

import (
	fieldValidator "github.com/jcfug8/daylear/server/adapters/services/grpc/fieldbehaviorvalidator"
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/fieldmasker"
	namer "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/namer"
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
	FieldValidator  fieldValidator.FieldBehaviorValidator
	Log             zerolog.Logger
	UserFieldMasker fieldMasker.UserFieldMasker
	UserNamer       namer.UserNamer
}

// NewUserService creates a new UserService.
func NewUserService(params NewUserServiceParams) *UserService {
	return &UserService{
		domain:                 params.Domain,
		fieldBehaviorValidator: params.FieldValidator,
		log:                    params.Log,
		userFieldMasker:        params.UserFieldMasker,
		userNamer:              params.UserNamer,
	}
}

// UserService defines the grpc handlers for the UserService.
type UserService struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedPublicUserServiceServer
	domain                 domain.Domain
	fieldBehaviorValidator fieldValidator.FieldBehaviorValidator
	log                    zerolog.Logger
	userFieldMasker        fieldMasker.UserFieldMasker
	userNamer              namer.UserNamer
}

// Register registers s to the grpc implementation of the service.
func (s *UserService) Register(server *grpc.Server) error {
	pb.RegisterUserServiceServer(server, s)
	pb.RegisterPublicUserServiceServer(server, s)
	return nil
}
