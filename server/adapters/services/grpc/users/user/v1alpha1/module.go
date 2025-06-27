package v1alpha1

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

var Module = fx.Module(
	"userGrpcAdapter",
	fx.Provide(
		NewUserService,
		func(s *UserService) pb.UserServiceServer { return s },
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.User]() },
			fx.ResultTags(`name="v1alpha1UserNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name="v1alpha1UserAccessNamer"`),
		),
	),
)
