package v1alpha1

import (
	"go.uber.org/fx"

	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

var Module = fx.Module(
	"userGrpcAdapter",
	fx.Provide(
		NewUserService,
		func(s *UserService) pb.UserServiceServer { return s },
		func(s *UserService) pb.PublicUserServiceServer { return s },
	),
)
