package v1alpha1

import (
	"context"
	"time"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	latlng "google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	eventMaxPageSize     int32 = 1000
	eventDefaultPageSize int32 = 100
)

var eventFieldMap = map[string][]string{
	"name":                 {model.EventField_Parent, model.EventField_EventId},
	"title":                {model.EventField_Title},
	"start_time":           {model.EventField_StartTime},
	"end_time":             {model.EventField_EndTime},
	"description":          {model.EventField_Description},
	"location":             {model.EventField_Location},
	"uri":                  {model.EventField_URL},
	"recurrence_rule":      {model.EventField_RecurrenceRule, model.EventField_RecurrenceEndTime},
	"overriden_start_time": {model.EventField_OverridenStartTime},
	"excluded_times":       {model.EventField_ExcludedDates},
	"additional_times":     {model.EventField_AdditionalDates},
	"parent_event":         {model.EventField_ParentEventId},
	"alarms":               {model.EventField_Alarms},
	"recurrence_end_time":  {model.EventField_RecurrenceEndTime},
}

// CreateEvent creates a new event
func (s *CalendarService) CreateEvent(ctx context.Context, request *pb.CreateEventRequest) (response *pb.Event, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateEvent called")
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
	eventProto := request.GetEvent()
	eventProto.Name = ""
	_, mEvent, err := s.ProtoToEvent(eventProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	_, err = s.eventNamer.ParseParent(request.GetParent(), &mEvent.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create event
	mEvent, err = s.domain.CreateEvent(ctx, authAccount, mEvent)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateEvent failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	eventProto, err = s.EventToProto(mEvent)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventProto)
	log.Info().Msg("gRPC CreateEvent returning successfully")
	return eventProto, nil
}

// GetEvent gets an event
func (s *CalendarService) GetEvent(ctx context.Context, request *pb.GetEventRequest) (*pb.Event, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetEvent called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mEvent model.Event
	_, err = s.eventNamer.Parse(request.GetName(), &mEvent)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mEvent, err = s.domain.GetEvent(ctx, authAccount, mEvent.Parent, mEvent.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetEvent failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	eventProto, err := s.EventToProto(mEvent)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventProto)
	log.Info().Msg("gRPC GetEvent returning successfully")
	return eventProto, nil
}

// ListEvents lists events
func (s *CalendarService) ListEvents(ctx context.Context, request *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListEvents called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mEvent model.Event
	_, err = s.eventNamer.ParseParent(request.GetParent(), &mEvent.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: eventDefaultPageSize,
		MaxPageSize:     eventMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}

	// list events
	mEvents, err := s.domain.ListEvents(ctx, authAccount, mEvent.Parent, pageSize, pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListEvents failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	eventProtos := make([]*pb.Event, len(mEvents))
	for i, mEvent := range mEvents {
		eventProto, err := s.EventToProto(mEvent)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		eventProtos[i] = eventProto
	}

	// check field behavior
	for _, eventProto := range eventProtos {
		grpc.ProcessResponseFieldBehavior(eventProto)
	}

	// create response
	response := &pb.ListEventsResponse{
		Events: eventProtos,
	}

	// add next page token if there are more results
	if len(mEvents) == int(pageSize) {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListEvents returning successfully")
	return response, nil
}

// UpdateEvent updates an event
func (s *CalendarService) UpdateEvent(ctx context.Context, request *pb.UpdateEventRequest) (*pb.Event, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateEvent called")
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
	eventProto := request.GetEvent()
	_, mEvent, err := s.ProtoToEvent(eventProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// get fields to update
	var fields []string
	if request.GetUpdateMask() != nil {
		fields = s.eventFieldMasker.Convert(request.GetUpdateMask().GetPaths())
	}

	// update event
	mEvent, err = s.domain.UpdateEvent(ctx, authAccount, mEvent, fields)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateEvent failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	eventProto, err = s.EventToProto(mEvent)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventProto)
	log.Info().Msg("gRPC UpdateEvent returning successfully")
	return eventProto, nil
}

// DeleteEvent deletes an event
func (s *CalendarService) DeleteEvent(ctx context.Context, request *pb.DeleteEventRequest) (*pb.Event, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteEvent called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mEvent model.Event
	_, err = s.eventNamer.Parse(request.GetName(), &mEvent)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mEvent, err = s.domain.DeleteEvent(ctx, authAccount, mEvent.Parent, mEvent.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteEvent failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	eventProto, err := s.EventToProto(mEvent)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventProto)
	log.Info().Msg("gRPC DeleteEvent returning successfully")
	return eventProto, nil
}

// ProtoToEvent converts a proto Event to a model Event
func (s *CalendarService) ProtoToEvent(proto *pb.Event) (nameIndex int, event model.Event, err error) {
	title := proto.GetTitle()
	description := proto.GetDescription()
	uri := proto.GetUri()

	event = model.Event{
		Title:       title,
		Description: description,
		URL:         uri,
	}

	// Handle start time
	if proto.GetStartTime() != nil {
		event.StartTime = proto.GetStartTime().AsTime()
	}

	// Handle end time
	if proto.GetEndTime() != nil {
		endTime := proto.GetEndTime().AsTime()
		event.EndTime = &endTime
	}

	// Handle geo (convert LatLng to string for now)
	if proto.GetGeo() != nil {
		event.Geo = model.LatLng{
			Latitude:  proto.GetGeo().GetLatitude(),
			Longitude: proto.GetGeo().GetLongitude(),
		}
	}

	// Handle location
	if proto.GetLocation() != "" {
		event.Location = proto.GetLocation()
	}

	// Handle recurrence rule
	if proto.GetRecurrenceRule() != "" {
		event.RecurrenceRule = &proto.RecurrenceRule
	}

	// Handle recurrence time
	if proto.GetOverridenStartTime() != nil {
		recurrenceTime := proto.GetOverridenStartTime().AsTime()
		event.OverridenStartTime = &recurrenceTime
	}

	// Handle excluded times
	if len(proto.GetExcludedTimes()) > 0 {
		excludedDates := make([]time.Time, len(proto.GetExcludedTimes()))
		for i, t := range proto.GetExcludedTimes() {
			excludedDates[i] = t.AsTime()
		}
		event.ExcludedDates = excludedDates
	}

	// Handle additional times
	if len(proto.GetAdditionalTimes()) > 0 {
		additionalDates := make([]time.Time, len(proto.GetAdditionalTimes()))
		for i, t := range proto.GetAdditionalTimes() {
			additionalDates[i] = t.AsTime()
		}
		event.AdditionalDates = additionalDates
	}

	// Handle alarms (basic conversion for now)
	if len(proto.GetAlarms()) > 0 {
		alarms := make([]*model.Alarm, len(proto.GetAlarms()))
		for i, alarmProto := range proto.GetAlarms() {
			alarm := &model.Alarm{
				AlarmId: alarmProto.GetAlarmId(),
			}
			// Handle trigger conversion
			if trigger := alarmProto.GetTrigger(); trigger != nil {
				modelTrigger := &model.Trigger{}
				if duration := trigger.GetDuration(); duration != nil {
					durationVal := duration.AsDuration()
					modelTrigger.Duration = &durationVal
				} else if dateTime := trigger.GetDateTime(); dateTime != nil {
					dateTimeVal := dateTime.AsTime()
					modelTrigger.DateTime = &dateTimeVal
				}
				alarm.Trigger = modelTrigger
			}
			alarms[i] = alarm
		}
		event.Alarms = alarms
	}

	// Parse parent from name if provided
	if proto.GetName() != "" {
		nameIndex, err = s.eventNamer.Parse(proto.GetName(), &event)
		if err != nil {
			return 0, model.Event{}, err
		}
	}

	// Handle recurring event id
	if proto.GetParentEvent() != "" {
		e := model.Event{}
		_, err = s.eventNamer.Parse(proto.GetParentEvent(), &e)
		if err != nil {
			return 0, model.Event{}, err
		}
		event.ParentEventId = &e.Id.EventId
	}

	return nameIndex, event, nil
}

// EventToProto converts a model Event to a proto Event
func (s *CalendarService) EventToProto(event model.Event, options ...namer.FormatReflectNamerOption) (*pb.Event, error) {
	proto := &pb.Event{}

	// Handle title
	if event.Title != "" {
		proto.Title = event.Title
	}

	// Handle description
	if event.Description != "" {
		proto.Description = event.Description
	}

	// Handle URL
	if event.URL != "" {
		proto.Uri = event.URL
	}

	// Handle start time
	if !event.StartTime.IsZero() {
		proto.StartTime = timestamppb.New(event.StartTime)
	}

	// Handle end time
	if event.EndTime != nil {
		proto.EndTime = timestamppb.New(*event.EndTime)
	}

	// Handle geo
	if event.Geo.Latitude != 0 && event.Geo.Longitude != 0 {
		proto.Geo = &latlng.LatLng{
			Latitude:  event.Geo.Latitude,
			Longitude: event.Geo.Longitude,
		}
	}

	// Handle location
	if event.Location != "" {
		proto.Location = event.Location
	}

	if event.RecurrenceEndTime != nil {
		proto.RecurrenceEndTime = timestamppb.New(*event.RecurrenceEndTime)
	}

	// Handle recurrence rule
	if event.RecurrenceRule != nil {
		proto.RecurrenceRule = *event.RecurrenceRule
	}

	// Handle recurrence time
	if event.OverridenStartTime != nil {
		proto.OverridenStartTime = timestamppb.New(*event.OverridenStartTime)
	}

	// Handle excluded times
	if len(event.ExcludedDates) > 0 {
		excludedTimes := make([]*timestamppb.Timestamp, len(event.ExcludedDates))
		for i, t := range event.ExcludedDates {
			excludedTimes[i] = timestamppb.New(t)
		}
		proto.ExcludedTimes = excludedTimes
	}

	// Handle additional times
	if len(event.AdditionalDates) > 0 {
		additionalTimes := make([]*timestamppb.Timestamp, len(event.AdditionalDates))
		for i, t := range event.AdditionalDates {
			additionalTimes[i] = timestamppb.New(t)
		}
		proto.AdditionalTimes = additionalTimes
	}

	// Handle alarms
	if len(event.Alarms) > 0 {
		alarms := make([]*pb.Event_Alarm, len(event.Alarms))
		for i, alarm := range event.Alarms {
			alarmProto := &pb.Event_Alarm{
				AlarmId: alarm.AlarmId,
			}
			// Handle trigger conversion
			if alarm.Trigger != nil {
				alarmProto.Trigger = &pb.Event_Alarm_Trigger{}
				if alarm.Trigger.Duration != nil {
					alarmProto.Trigger.Trigger = &pb.Event_Alarm_Trigger_Duration{
						Duration: durationpb.New(*alarm.Trigger.Duration),
					}
				} else if alarm.Trigger.DateTime != nil {
					alarmProto.Trigger.Trigger = &pb.Event_Alarm_Trigger_DateTime{
						DateTime: timestamppb.New(*alarm.Trigger.DateTime),
					}
				}
			}
			alarms[i] = alarmProto
		}
		proto.Alarms = alarms
	}

	// Generate name
	if event.Id.EventId != 0 {
		name, err := s.eventNamer.Format(event, options...)
		if err != nil {
			return nil, err
		}
		proto.Name = name
	}

	// Handle parent event
	if event.ParentEventId != nil {
		e := model.Event{
			Parent: model.EventParent{
				CalendarId: event.Parent.CalendarId,
			},
			Id: model.EventId{
				EventId: *event.ParentEventId,
			},
		}
		name, err := s.eventNamer.Format(e, options...)
		if err != nil {
			return nil, err
		}
		proto.ParentEvent = name
	}

	return proto, nil
}
