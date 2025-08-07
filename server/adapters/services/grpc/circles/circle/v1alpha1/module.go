package v1alpha1

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
)

var Module = fx.Module(
	"circleGrpcAdapter",
	fx.Provide(
		NewCircleService,
		func(s *CircleService) pb.CircleServiceServer { return s },
		func(s *CircleService) pb.CircleAccessServiceServer { return s },
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name:"v1alpha1CircleAccessNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Circle]() },
			fx.ResultTags(`name:"v1alpha1CircleNamer"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Circle{}, circleFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CircleFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, circleAccessFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CircleAccessFieldMasker"`),
		),
	),
)
