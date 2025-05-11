package v1alpha1

import (
	"go.uber.org/fx"

	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
)

var Module = fx.Module(
	"circleGrpcAdapter",
	fx.Provide(
		NewCircleService,
		func(s *CircleService) pb.CircleServiceServer { return s },
		func(s *CircleService) pb.PublicCircleServiceServer { return s },
	),
)
