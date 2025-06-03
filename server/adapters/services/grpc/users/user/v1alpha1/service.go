package v1alpha1

import (
	fieldValidator "github.com/jcfug8/daylear/server/adapters/services/grpc/fieldbehaviorvalidator"
	fieldMasker "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/fieldmasker"
	"github.com/jcfug8/daylear/server/core/model"
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
	FieldValidator  fieldValidator.FieldBehaviorValidator
	Log             zerolog.Logger
	UserFieldMasker fieldMasker.UserFieldMasker
}

// NewUserService creates a new UserService.
func NewUserService(params NewUserServiceParams) (*UserService, error) {
	userNamer, err := namer.NewReflectNamer[model.User, *pb.User]()
	if err != nil {
		return nil, err
	}

	publicUserNamer, err := namer.NewReflectNamer[model.User, *pb.PublicUser]()
	if err != nil {
		return nil, err
	}

	return &UserService{
		domain:                 params.Domain,
		fieldBehaviorValidator: params.FieldValidator,
		log:                    params.Log,
		userFieldMasker:        params.UserFieldMasker,
		userNamer:              userNamer,
		publicUserNamer:        publicUserNamer,
	}, nil
}

// UserService defines the grpc handlers for the UserService.
type UserService struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedPublicUserServiceServer
	domain                 domain.Domain
	fieldBehaviorValidator fieldValidator.FieldBehaviorValidator
	log                    zerolog.Logger
	userFieldMasker        fieldMasker.UserFieldMasker
	userNamer              namer.ReflectNamer[model.User]
	publicUserNamer        namer.ReflectNamer[model.User]
}

// Register registers s to the grpc implementation of the service.
func (s *UserService) Register(server *grpc.Server) error {
	pb.RegisterUserServiceServer(server, s)
	pb.RegisterPublicUserServiceServer(server, s)
	return nil
}
