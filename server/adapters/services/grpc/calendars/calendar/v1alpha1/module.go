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
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Calendar]() },
			fx.ResultTags(`name:"v1alpha1CalendarNamer"`),
		),
		fx.Annotate(
			func() (namer.ReflectNamer, error) { return namer.NewReflectNamer[*pb.Access]() },
			fx.ResultTags(`name:"v1alpha1CalendarAccessNamer"`),
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
	),
)
