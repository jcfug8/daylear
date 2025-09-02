package v1alpha1

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/lists/list/v1alpha1"
)

var Module = fx.Module(
	"listGrpcAdapter",
	fx.Provide(
		NewListService,
		func(s *ListService) pb.ListServiceServer { return s },
		func(s *ListService) pb.ListAccessServiceServer { return s },
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name:"v1alpha1ListAccessNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.List]() },
			fx.ResultTags(`name:"v1alpha1ListNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) {
				return namer.NewReflectNamer[*pb.List_ListSection](namer.WithExtraPatterns([]string{"lists/{list}/listSections/{list_section}"}))
			},
			fx.ResultTags(`name:"v1alpha1ListSectionNamer"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.List{}, listFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1ListFieldMasker"`),
		),

		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, listAccessFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1ListAccessFieldMasker"`),
		),
	),
)
