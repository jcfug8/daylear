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
	EventFieldMasker          fieldmask.FieldMasker `name:"v1alpha1EventFieldMasker"`
	UserNamer                 namer.ReflectNamer    `name:"v1alpha1UserNamer"`
	CircleNamer               namer.ReflectNamer    `name:"v1alpha1CircleNamer"`
	EventNamer                namer.ReflectNamer    `name:"v1alpha1EventNamer"`
	EventRecipeNamer          namer.ReflectNamer    `name:"v1alpha1EventRecipeNamer"`
	RecipeNamer               namer.ReflectNamer    `name:"v1alpha1RecipeNamer"`
	EventRecipeFieldMasker    fieldmask.FieldMasker `name:"v1alpha1EventRecipeFieldMasker"`
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
		eventFieldMasker:          params.EventFieldMasker,
		userNamer:                 params.UserNamer,
		circleNamer:               params.CircleNamer,
		eventNamer:                params.EventNamer,
		eventRecipeNamer:          params.EventRecipeNamer,
		recipeNamer:               params.RecipeNamer,
		eventRecipeFieldMasker:    params.EventRecipeFieldMasker,
	}, nil
}

// CalendarService defines the grpc handlers for the CalendarService.
type CalendarService struct {
	pb.UnimplementedCalendarServiceServer
	pb.UnimplementedCalendarAccessServiceServer
	pb.UnimplementedEventServiceServer
	pb.UnimplementedEventRecipeServiceServer
	domain                    domain.Domain
	log                       zerolog.Logger
	calendarNamer             namer.ReflectNamer
	calendarAccessNamer       namer.ReflectNamer
	calendarFieldMasker       fieldmask.FieldMasker
	calendarAccessFieldMasker fieldmask.FieldMasker
	eventFieldMasker          fieldmask.FieldMasker
	userNamer                 namer.ReflectNamer
	circleNamer               namer.ReflectNamer
	eventNamer                namer.ReflectNamer
	eventRecipeNamer          namer.ReflectNamer
	recipeNamer               namer.ReflectNamer
	eventRecipeFieldMasker    fieldmask.FieldMasker
}

// Register registers s to the grpc implementation of the service.
func (s *CalendarService) Register(server *grpc.Server) error {
	pb.RegisterCalendarServiceServer(server, s)
	pb.RegisterCalendarAccessServiceServer(server, s)
	pb.RegisterEventServiceServer(server, s)
	pb.RegisterEventRecipeServiceServer(server, s)
	return nil
}
