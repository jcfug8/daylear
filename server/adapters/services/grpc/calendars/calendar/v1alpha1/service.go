package v1alpha1

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
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
	CalendarFieldMasker       fieldmask.FieldMasker `name:"v1alpha1CalendarFieldMasker"`
	CalendarAccessFieldMasker fieldmask.FieldMasker `name:"v1alpha1CalendarAccessFieldMasker"`
}

// NewCalendarService creates a new CalendarService.
func NewCalendarService(params NewCalendarServiceParams) (*CalendarService, error) {
	return &CalendarService{
		domain:                    params.Domain,
		log:                       params.Log,
		calendarNamer:             params.CalendarNamer,
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
	calendarFieldMasker       fieldmask.FieldMasker
	calendarAccessFieldMasker fieldmask.FieldMasker
}

// Register registers s to the grpc implementation of the service.
func (s *CalendarService) Register(server *grpc.Server) error {
	pb.RegisterCalendarServiceServer(server, s)
	return nil
}

// ProtoToCalendar converts a proto Calendar to a model Calendar
func (s *CalendarService) ProtoToCalendar(proto *pb.Calendar) (model.Calendar, error) {
	calendar := model.Calendar{
		CalendarId:      0, // Will be set by the database
		Title:           proto.GetTitle(),
		Description:     proto.GetDescription(),
		VisibilityLevel: proto.GetVisibility(),
	}

	// Parse parent from name if provided
	if proto.GetName() != "" {
		_, err := s.calendarNamer.Parse(proto.GetName(), &calendar)
		if err != nil {
			return model.Calendar{}, err
		}
	}

	return calendar, nil
}

// CalendarToProto converts a model Calendar to a proto Calendar
func (s *CalendarService) CalendarToProto(calendar model.Calendar, nameIndex ...namer.FormatReflectNamerOption) (*pb.Calendar, error) {
	proto := &pb.Calendar{
		Title:       calendar.Title,
		Description: calendar.Description,
		Visibility:  calendar.VisibilityLevel,
	}

	// Generate name
	name, err := s.calendarNamer.Format(calendar, nameIndex...)
	if err != nil {
		return nil, err
	}
	proto.Name = name

	// Add calendar access data if available
	if calendar.CalendarAccess.CalendarAccessId.CalendarAccessId != 0 {
		proto.CalendarAccess = &pb.Calendar_CalendarAccess{
			Name:            "", // TODO: Generate access name
			PermissionLevel: calendar.CalendarAccess.PermissionLevel,
			State:           calendar.CalendarAccess.State,
		}
	}

	return proto, nil
}
