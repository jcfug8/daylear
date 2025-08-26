package caldav

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/icalendar"
	"github.com/jcfug8/daylear/server/core/model"
)

// GetCalendarObject returns a specific event
func (b *Service) GetCalendarObject(ctx context.Context, path string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	b.log.Info().Str("path", path).Msg("GetCalendarObject called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse auth data in GetCalendarObject")
		return nil, fmt.Errorf("error parsing auth data: %w", err)
	}

	b.log.Info().Str("path", path).Int64("authUserID", authAccount.AuthUserId).Msg("Parsing calendar object path")
	userID, calendarID, eventID, err := parseCalendarObjectPath(path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse calendar object path")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Int64("authUserID", authAccount.AuthUserId).Msg("Parsed calendar object path")

	if userID != authAccount.AuthUserId {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("User ID mismatch in GetCalendarObject")
		return nil, fmt.Errorf("user ID does not match")
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Getting calendar from domain")
	calendar, err := b.domain.GetCalendar(ctx, model.AuthAccount{AuthUserId: userID}, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Failed to get calendar from domain")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Getting event from domain")
	event, err := b.domain.GetEvent(ctx, model.AuthAccount{AuthUserId: userID}, model.EventParent{CalendarId: calendarID, UserId: userID}, model.EventId{EventId: eventID}, []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Failed to get event from domain")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Str("eventTitle", event.Title).Msg("GetCalendarObject returning")

	return &caldav.CalendarObject{
		Path:    path,
		ModTime: event.UpdateTime,
		// ContentLength: 0,
		ETag: strconv.FormatInt(event.UpdateTime.UnixNano(), 10),
		Data: icalendar.ToICalendar(calendar, []model.Event{event}),
	}, nil
}

// ListCalendarObjects returns all events in a calendar
func (b *Service) ListCalendarObjects(ctx context.Context, path string, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	b.log.Info().Str("path", path).Msg("ListCalendarObjects called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse auth data in ListCalendarObjects")
		return nil, fmt.Errorf("error parsing auth data: %w", err)
	}

	b.log.Info().Str("path", path).Int64("authUserID", authAccount.AuthUserId).Msg("Parsing calendar path for ListCalendarObjects")
	userID, calendarID, err := parseCalendarPath(path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse calendar path in ListCalendarObjects")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("authUserID", authAccount.AuthUserId).Msg("Parsed calendar path for ListCalendarObjects")

	if userID != authAccount.AuthUserId {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("User ID mismatch in ListCalendarObjects")
		return nil, fmt.Errorf("user ID does not match")
	}

	//TODO: handle AllProps and Props

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Getting calendar from domain for ListCalendarObjects")
	calendar, err := b.domain.GetCalendar(ctx, model.AuthAccount{AuthUserId: userID}, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain for ListCalendarObjects")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Listing events from domain for ListCalendarObjects")
	events, err := b.domain.ListEvents(ctx, model.AuthAccount{AuthUserId: userID}, model.EventParent{CalendarId: calendarID, UserId: userID}, 1000, 0, "", []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to list events from domain for ListCalendarObjects")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int("eventCount", len(events)).Msg("ListCalendarObjects returning")

	return []caldav.CalendarObject{
		{
			Path:    path,
			ModTime: calendar.UpdateTime,
			// ContentLength: 0,
			ETag: strconv.FormatInt(calendar.UpdateTime.UnixNano(), 10),
			Data: icalendar.ToICalendar(calendar, events),
		},
	}, nil
}

// QueryCalendarObjects handles calendar-query REPORT requests
func (b *Service) QueryCalendarObjects(ctx context.Context, path string, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	b.log.Info().Str("path", path).Msg("QueryCalendarObjects called")

	calendarObjects, err := b.ListCalendarObjects(ctx, path, &caldav.CalendarCompRequest{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to list calendar objects in QueryCalendarObjects")
		return nil, err
	}

	b.log.Info().Str("path", path).Int("objectCount", len(calendarObjects)).Msg("QueryCalendarObjects returning")
	return calendarObjects, nil
}

// PutCalendarObject creates or updates an event
func (b *Service) PutCalendarObject(ctx context.Context, path string, input *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (*caldav.CalendarObject, error) {
	b.log.Info().Str("path", path).Msg("PutCalendarObject called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse auth data in PutCalendarObject")
		return nil, fmt.Errorf("error parsing auth data: %w", err)
	}

	b.log.Info().Str("path", path).Int64("authUserID", authAccount.AuthUserId).Msg("Parsing calendar object path for PutCalendarObject")
	userID, calendarID, _, err := parseCalendarObjectPath(path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse calendar object path in PutCalendarObject")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("authUserID", authAccount.AuthUserId).Msg("Parsed calendar object path for PutCalendarObject")

	if userID != authAccount.AuthUserId {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("User ID mismatch in PutCalendarObject")
		return nil, fmt.Errorf("user ID does not match")
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Converting iCalendar to domain models")
	calendar, events, err := icalendar.FromICalendar(input)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to convert iCalendar in PutCalendarObject")
		return nil, err
	}

	if calendar.CalendarId.CalendarId != calendarID {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("inputCalendarID", calendar.CalendarId.CalendarId).Msg("Calendar ID mismatch in PutCalendarObject")
		return nil, fmt.Errorf("calendar ID does not match")
	}

	if len(events) == 0 {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("No events found in iCalendar in PutCalendarObject")
		return nil, fmt.Errorf("no events found")
	}

	if len(events) > 1 {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int("eventCount", len(events)).Msg("Multiple events found in iCalendar in PutCalendarObject")
		return nil, fmt.Errorf("multiple events found")
	}

	event := events[0]
	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Str("eventTitle", event.Title).Msg("Processing single event from iCalendar")

	event.Parent = model.EventParent{CalendarId: calendarID, UserId: userID}

	// TODO: handle IfNoneMatch and IfMatch options
	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Str("eventTitle", event.Title).Msg("Updating event in domain")
	event, err = b.domain.UpdateEvent(ctx, model.AuthAccount{AuthUserId: userID}, event, []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Str("eventTitle", event.Title).Msg("Failed to update event in domain")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Str("eventTitle", event.Title).Msg("PutCalendarObject returning")
	return &caldav.CalendarObject{
		Path:    path,
		ModTime: event.UpdateTime,
		// ContentLength: 0,
		ETag: strconv.FormatInt(event.UpdateTime.UnixNano(), 10),
		Data: icalendar.ToICalendar(calendar, []model.Event{event}),
	}, nil
}

// DeleteCalendarObject deletes an event
func (b *Service) DeleteCalendarObject(ctx context.Context, path string) error {
	b.log.Info().Str("path", path).Msg("DeleteCalendarObject called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse auth data in DeleteCalendarObject")
		return fmt.Errorf("error parsing auth data: %w", err)
	}

	b.log.Info().Str("path", path).Int64("authUserID", authAccount.AuthUserId).Msg("Parsing calendar object path for DeleteCalendarObject")
	userID, calendarID, eventID, err := parseCalendarObjectPath(path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse calendar object path in DeleteCalendarObject")
		return err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Int64("authUserID", authAccount.AuthUserId).Msg("Parsed calendar object path for DeleteCalendarObject")

	if userID != authAccount.AuthUserId {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("User ID mismatch in DeleteCalendarObject")
		return fmt.Errorf("user ID does not match")
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Deleting event from domain")
	_, err = b.domain.DeleteEvent(ctx, model.AuthAccount{AuthUserId: userID}, model.EventParent{CalendarId: calendarID, UserId: userID}, model.EventId{EventId: eventID})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("Failed to delete event from domain")
		return err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("eventID", eventID).Msg("DeleteCalendarObject completed successfully")
	return nil
}

// TODO: figure out how other objects are handled
func parseCalendarObjectPath(path string) (userID, calendarID, objectID int64, err error) {
	// Note: This function doesn't have access to the logger, so we can't log here
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[0] != "calendars" {
		return 0, 0, 0, fmt.Errorf("invalid event path: %s", path)
	}

	userID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	calendarID, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	// Remove .ics extension
	eventIDStr := strings.TrimSuffix(parts[3], ".ics")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}

	return userID, calendarID, eventID, nil
}
