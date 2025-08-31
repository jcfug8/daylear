package v1alpha1

import (
	"go.uber.org/fx"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
)

var Module = fx.Module(
	"calendarGrpcAdapter",
	fx.Provide(
		NewCalendarService,
		func(s *CalendarService) pb.CalendarServiceServer { return s },
		func(s *CalendarService) pb.CalendarAccessServiceServer { return s },
		func(s *CalendarService) pb.EventServiceServer { return s },
		func(s *CalendarService) pb.EventRecipeServiceServer { return s },
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Calendar]() },
			fx.ResultTags(`name:"v1alpha1CalendarNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name:"v1alpha1CalendarAccessNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Event]() },
			fx.ResultTags(`name:"v1alpha1EventNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.EventRecipe]() },
			fx.ResultTags(`name:"v1alpha1EventRecipeNamer"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Calendar{}, calendarFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CalendarFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, calendarAccessFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CalendarAccessFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Event{}, eventFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1EventFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.EventRecipe{}, eventRecipeFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1EventRecipeFieldMasker"`),
		),
	),
)
