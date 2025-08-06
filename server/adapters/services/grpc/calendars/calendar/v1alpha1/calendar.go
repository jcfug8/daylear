package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	calendarMaxPageSize     int32 = 1000
	calendarDefaultPageSize int32 = 100
)

var calendarFieldMap = map[string][]string{
	"name":        {model.CalendarField_Parent, model.CalendarField_CalendarId},
	"title":       {model.CalendarField_Title},
	"description": {model.CalendarField_Description},
	"visibility":  {model.CalendarField_Visibility},

	"calendar_access": {model.CalendarField_CalendarAccess},
}

// CreateCalendar creates a new calendar
func (s *CalendarService) CreateCalendar(ctx context.Context, request *pb.CreateCalendarRequest) (response *pb.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateCalendar called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	calendarProto := request.GetCalendar()
	calendarProto.Name = ""
	_, mCalendar, err := s.ProtoToCalendar(calendarProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	nameIndex, err := s.calendarNamer.ParseParent(request.GetParent(), &mCalendar.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create calendar
	mCalendar, err = s.domain.CreateCalendar(ctx, authAccount, mCalendar)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateCalendar failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	calendarProto, err = s.CalendarToProto(mCalendar, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(calendarProto)
	log.Info().Msg("gRPC CreateCalendar returning successfully")
	return calendarProto, nil
}

// DeleteCalendar deletes a calendar
func (s *CalendarService) DeleteCalendar(ctx context.Context, request *pb.DeleteCalendarRequest) (*pb.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteCalendar called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mCalendar model.Calendar
	nameIndex, err := s.calendarNamer.Parse(request.GetName(), &mCalendar)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCalendar, err = s.domain.DeleteCalendar(ctx, authAccount, mCalendar.Parent, mCalendar.CalendarId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteCalendar failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	calendarProto, err := s.CalendarToProto(mCalendar, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(calendarProto)
	log.Info().Msg("gRPC DeleteCalendar returning successfully")
	return calendarProto, nil
}

// GetCalendar retrieves a calendar
func (s *CalendarService) GetCalendar(ctx context.Context, request *pb.GetCalendarRequest) (*pb.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetCalendar called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mCalendar model.Calendar
	nameIndex, err := s.calendarNamer.Parse(request.GetName(), &mCalendar)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCalendar, err = s.domain.GetCalendar(ctx, authAccount, mCalendar.Parent, mCalendar.CalendarId, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetCalendar failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	calendarProto, err := s.CalendarToProto(mCalendar, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(calendarProto)
	log.Info().Msg("gRPC GetCalendar returning successfully")
	return calendarProto, nil
}

// UpdateCalendar updates a calendar
func (s *CalendarService) UpdateCalendar(ctx context.Context, request *pb.UpdateCalendarRequest) (*pb.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateCalendar called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	calendarProto := request.GetCalendar()
	nameIndex, mCalendar, err := s.ProtoToCalendar(calendarProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// get update mask
	fieldMask := request.GetUpdateMask()
	updateMask := s.calendarFieldMasker.Convert(fieldMask.GetPaths())

	// update calendar
	mCalendar, err = s.domain.UpdateCalendar(ctx, authAccount, mCalendar, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateCalendar failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	calendarProto, err = s.CalendarToProto(mCalendar, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(calendarProto)
	log.Info().Msg("gRPC UpdateCalendar returning successfully")
	return calendarProto, nil
}

// ListCalendars lists calendars
func (s *CalendarService) ListCalendars(ctx context.Context, request *pb.ListCalendarsRequest) (*pb.ListCalendarsResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListCalendars called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent
	var mCalendar model.Calendar
	nameIndex, err := s.calendarNamer.ParseParent(request.GetParent(), &mCalendar)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.GetParent()).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: calendarDefaultPageSize,
		MaxPageSize:     calendarMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}

	// list calendars
	mCalendars, err := s.domain.ListCalendars(ctx, authAccount, mCalendar.Parent, pageSize, pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListCalendars failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	calendarProtos := make([]*pb.Calendar, len(mCalendars))
	for i, mCalendar := range mCalendars {
		calendarProto, err := s.CalendarToProto(mCalendar, namer.AsPatternIndex(nameIndex))
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		calendarProtos[i] = calendarProto
	}

	// check field behavior
	for _, calendarProto := range calendarProtos {
		grpc.ProcessResponseFieldBehavior(calendarProto)
	}

	// create response
	response := &pb.ListCalendarsResponse{
		Calendars: calendarProtos,
	}

	// add next page token if there are more results
	if len(mCalendars) == int(pageSize) {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListCalendars returning successfully")
	return response, nil
}

// ProtoToCalendar converts a proto Calendar to a model Calendar
func (s *CalendarService) ProtoToCalendar(proto *pb.Calendar) (nameIndex int, calendar model.Calendar, err error) {
	calendar = model.Calendar{
		Title:           proto.GetTitle(),
		Description:     proto.GetDescription(),
		VisibilityLevel: proto.GetVisibility(),
	}

	// Parse parent from name if provided
	if proto.GetName() != "" {
		nameIndex, err = s.calendarNamer.Parse(proto.GetName(), &calendar)
		if err != nil {
			return 0, model.Calendar{}, err
		}
	}

	return nameIndex, calendar, nil
}

// CalendarToProto converts a model Calendar to a proto Calendar
func (s *CalendarService) CalendarToProto(calendar model.Calendar, options ...namer.FormatReflectNamerOption) (*pb.Calendar, error) {
	proto := &pb.Calendar{
		Title:       calendar.Title,
		Description: calendar.Description,
		Visibility:  calendar.VisibilityLevel,
	}

	// Generate name
	if calendar.CalendarId.CalendarId != 0 {
		name, err := s.calendarNamer.Format(calendar, options...)
		if err != nil {
			return nil, err
		}
		proto.Name = name
	}

	// Add calendar access data if available
	if calendar.CalendarAccess.CalendarAccessId.CalendarAccessId != 0 {
		name, err := s.calendarAccessNamer.Format(calendar.CalendarAccess)
		if err == nil {
			proto.CalendarAccess = &pb.Calendar_CalendarAccess{
				Name:            name,
				PermissionLevel: calendar.CalendarAccess.PermissionLevel,
				State:           calendar.CalendarAccess.State,
			}
		}
	}

	return proto, nil
}
