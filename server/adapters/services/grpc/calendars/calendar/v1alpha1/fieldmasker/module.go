package fieldmasker

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"go.uber.org/fx"
)

var updateCalendarAccessFieldMap = map[string][]string{
	"level": {"permission_level"},
	"state": {"state"},
}

var updateCalendarFieldMap = map[string][]string{}

// Module -
var Module = fx.Module(
	"fieldmasker",
	fx.Provide(
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, updateCalendarFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CalendarFieldMasker"`),
		),
		fx.Annotate(
			func() (fieldmask.FieldMasker, error) {
				return fieldmask.NewProtoFieldMasker(&pb.Access{}, updateCalendarAccessFieldMap)
			},
			fx.ResultTags(`name:"v1alpha1CalendarAccessFieldMasker"`),
		),
	),
)
