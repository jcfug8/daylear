package v1alpha1

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	domain "github.com/jcfug8/daylear/server/ports/domain"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

// NewCalendarServiceParams defines the dependencies for the CalendarService.
type NewCalendarServiceParams struct {
	fx.In

	Domain                    domain.Domain
	Log                       zerolog.Logger
	CalendarNamer             namer.ReflectNamer    `name:"v1alpha1CalendarNamer"`
	CalendarAccessNamer       namer.ReflectNamer    `name:"v1alpha1CalendarAccessNamer"`
	CalendarFieldMasker       fieldmask.FieldMasker `name:"v1alpha1CalendarFieldMasker"`
	CalendarAccessFieldMasker fieldmask.FieldMasker `name:"v1alpha1CalendarAccessFieldMasker"`
}

// NewCalendarService creates a new CalendarService.
func NewCalendarService(params NewCalendarServiceParams) (*CalendarService, error) {
	return &CalendarService{
		domain:                    params.Domain,
		log:                       params.Log,
		calendarNamer:             params.CalendarNamer,
		calendarAccessNamer:       params.CalendarAccessNamer,
		calendarFieldMasker:       params.CalendarFieldMasker,
		calendarAccessFieldMasker: params.CalendarAccessFieldMasker,
	}, nil
}

// CalendarService defines the grpc handlers for the CalendarService.
type CalendarService struct {
	pb.UnimplementedCalendarServiceServer
	domain                    domain.Domain
	log                       zerolog.Logger
	calendarNamer             namer.ReflectNamer
	calendarAccessNamer       namer.ReflectNamer
	calendarFieldMasker       fieldmask.FieldMasker
	calendarAccessFieldMasker fieldmask.FieldMasker
}

// Register registers s to the grpc implementation of the service.
func (s *CalendarService) Register(server *grpc.Server) error {
	pb.RegisterCalendarServiceServer(server, s)
	return nil
}
