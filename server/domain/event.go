package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateEvent creates a new event
func (d *Domain) CreateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event) (dbEvent model.Event, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if event.Parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	if event.StartTime.IsZero() {
		log.Error().Msg("start time is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "start time is required"}
	}

	if event.EndTime == nil || event.EndTime.IsZero() {
		log.Error().Msg("end time is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "end time is required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: event.Parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating event")
		return model.Event{}, err
	}

	if event.RecurrenceRule != nil && *event.RecurrenceRule != "" {
		event.RecurrenceEndTime = event.GetUntil()
	}

	dbEvent, err = d.repo.CreateEvent(ctx, event, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to create event"}
	}

	return dbEvent, nil
}

// DeleteEvent deletes an event
func (d *Domain) DeleteEvent(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, id model.EventId) (dbEvent model.Event, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// TODO: Implement event deletion logic
	// This is a stub implementation - actual business logic will be added later

	log.Info().Msg("DeleteEvent called - stub implementation")
	return model.Event{}, nil
}

// GetEvent retrieves an event
func (d *Domain) GetEvent(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, id model.EventId, fields []string) (dbEvent model.Event, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	if id.EventId == 0 {
		log.Error().Msg("event id is required when getting event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "event id is required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating event")
		return model.Event{}, err
	}

	dbEvent, err = d.repo.GetEvent(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to get event"}
	}

	return dbEvent, nil
}

// ListEvents lists events with pagination and filtering
func (d *Domain) ListEvents(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, pageSize int32, offset int64, filter string, fields []string) (dbEvents []model.Event, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when listing events")
		return []model.Event{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required when listing events")
		return []model.Event{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	dbCalendar, err := d.GetCalendar(ctx, authAccount, model.CalendarParent{}, model.CalendarId{CalendarId: parent.CalendarId}, []string{model.CalendarField_Visibility})
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar")
		return []model.Event{}, domain.ErrInternal{Msg: "unable to get calendar"}
	}

	_, err = d.determineCalendarAccess(
		ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
		withResourceVisibilityLevel(dbCalendar.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when listing events")
		return []model.Event{}, err
	}

	dbEvents, err = d.repo.ListEvents(ctx, authAccount, parent, pageSize, offset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list events")
		return []model.Event{}, domain.ErrInternal{Msg: "unable to list events"}
	}

	return dbEvents, nil
}

// UpdateEvent updates an existing event
func (d *Domain) UpdateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event, fields []string) (dbEvent model.Event, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if event.Parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	if event.StartTime.IsZero() {
		log.Error().Msg("start time is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "start time is required"}
	}

	if event.EndTime == nil || event.EndTime.IsZero() {
		log.Error().Msg("end time is required when creating event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "end time is required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: event.Parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating event")
		return model.Event{}, err
	}

	// remove duplicate excluded dates
	if event.ExcludedDates != nil {
		event.ExcludedDates = slices.Compact(event.ExcludedDates)
	}

	// remove duplicate additional dates
	if event.AdditionalDates != nil {
		event.AdditionalDates = slices.Compact(event.AdditionalDates)
	}

	//TODO: need to eventually add recurring logic

	if slices.Contains(fields, model.EventField_RecurrenceRule) {
		if event.RecurrenceRule != nil && *event.RecurrenceRule != "" {
			event.RecurrenceEndTime = event.GetUntil()
		}
	}

	dbEvent, err = d.repo.UpdateEvent(ctx, authAccount, event, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to update event"}
	}

	return dbEvent, nil
}
