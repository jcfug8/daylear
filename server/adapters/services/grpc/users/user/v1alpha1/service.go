package v1alpha1

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
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

	Domain                  domain.Domain
	Log                     zerolog.Logger
	UserFieldMasker         fieldmask.FieldMasker `name:"v1alpha1UserFieldMasker"`
	UserNamer               namer.ReflectNamer    `name:"v1alpha1UserNamer"`
	AccessNamer             namer.ReflectNamer    `name:"v1alpha1UserAccessNamer"`
	UserSettingsNamer       namer.ReflectNamer    `name:"v1alpha1UserSettingsNamer"`
	UserSettingsFieldMasker fieldmask.FieldMasker `name:"v1alpha1UserSettingsFieldMasker"`
}

// NewUserService creates a new UserService.
func NewUserService(params NewUserServiceParams) (*UserService, error) {
	return &UserService{
		domain:                  params.Domain,
		log:                     params.Log,
		userFieldMasker:         params.UserFieldMasker,
		userNamer:               params.UserNamer,
		accessNamer:             params.AccessNamer,
		userSettingsNamer:       params.UserSettingsNamer,
		userSettingsFieldMasker: params.UserSettingsFieldMasker,
	}, nil
}

// UserService defines the grpc handlers for the UserService.
type UserService struct {
	pb.UnimplementedUserServiceServer
	pb.UnimplementedUserAccessServiceServer
	pb.UnimplementedUserSettingsServiceServer
	domain                  domain.Domain
	log                     zerolog.Logger
	userFieldMasker         fieldmask.FieldMasker `name:"v1alpha1UserFieldMasker"`
	userNamer               namer.ReflectNamer
	accessNamer             namer.ReflectNamer
	userSettingsNamer       namer.ReflectNamer
	userSettingsFieldMasker fieldmask.FieldMasker
}

// Register registers s to the grpc implementation of the service.
func (s *UserService) Register(server *grpc.Server) error {
	pb.RegisterUserServiceServer(server, s)
	pb.RegisterUserAccessServiceServer(server, s)
	pb.RegisterUserSettingsServiceServer(server, s)
	return nil
}
