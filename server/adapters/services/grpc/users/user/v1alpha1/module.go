package v1alpha1

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

var Module = fx.Module(
	"userGrpcAdapter",
	fx.Provide(
		fx.Annotate(
			NewUserService,
			fx.As(new(pb.UserServiceServer)),
			fx.As(new(pb.UserSettingsServiceServer)),
			fx.As(new(pb.UserAccessServiceServer)),
			fx.As(new(pb.AccessKeyServiceServer)),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.User]() },
			fx.ResultTags(`name:"v1alpha1UserNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name:"v1alpha1UserAccessNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.UserSettings]() },
			fx.ResultTags(`name:"v1alpha1UserSettingsNamer"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.User{}, userFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1UserFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, userAccessFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1UserAccessFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.UserSettings{}, userSettingsFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1UserSettingsFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.AccessKey{}, accessKeyFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1AccessKeyFieldMasker"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.AccessKey]() },
			fx.ResultTags(`name:"v1alpha1AccessKeyNamer"`),
		),
	),
)
