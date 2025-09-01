package domain

import (
	"context"
	"fmt"
	"slices"
	"time"

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

	// Validate that start time is before end time
	if event.StartTime.After(*event.EndTime) || event.StartTime.Equal(*event.EndTime) {
		log.Error().Msg("start time must be before end time")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "start time must be before end time"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: event.Parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating event")
		return model.Event{}, err
	}

	if event.RecurrenceRule != nil && *event.RecurrenceRule != "" {
		event.RecurrenceEndTime = event.GetLastOccurence(true)
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

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when deleting event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required when deleting event")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when deleting event")
		return model.Event{}, err
	}

	dbEvent, err = d.repo.DeleteEvent(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to delete event"}
	}

	err = d.repo.DeleteChildEvents(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete child events")
		return model.Event{}, domain.ErrInternal{Msg: "unable to delete child events"}
	}

	err = d.repo.BulkDeleteEventRecipes(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete event recipes")
		return model.Event{}, domain.ErrInternal{Msg: "unable to delete event recipes"}
	}

	return dbEvent, nil
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

	// Validate that start time is before end time
	if event.StartTime.After(*event.EndTime) || event.StartTime.Equal(*event.EndTime) {
		log.Error().Msg("start time must be before end time")
		return model.Event{}, domain.ErrInvalidArgument{Msg: "start time must be before end time"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: event.Parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating event")
		return model.Event{}, err
	}

	dbOldEvent, err := d.repo.GetEvent(ctx, authAccount, event.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get old event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to get old event"}
	}

	// remove duplicate excluded dates
	if event.ExcludedDates != nil {
		event.ExcludedDates = slices.Compact(event.ExcludedDates)
	}

	// remove duplicate additional dates
	if event.AdditionalDates != nil {
		event.AdditionalDates = slices.Compact(event.AdditionalDates)
	}

	if slices.Contains(fields, model.EventField_StartTime) {
		diff := event.StartTime.Sub(dbOldEvent.StartTime)
		for i := range event.ExcludedDates {
			event.ExcludedDates[i] = event.ExcludedDates[i].Add(diff)
		}
	}

	if slices.Contains(fields, model.EventField_RecurrenceRule) {
		event.RecurrenceEndTime = nil
		if event.RecurrenceRule != nil && *event.RecurrenceRule != "" {
			event.RecurrenceEndTime = event.GetLastOccurence(true)
			filter := fmt.Sprintf("parent_event_id = %d AND start_time > '%s'", event.Id.EventId, event.GetLastOccurence(false).UTC().Format(time.RFC3339))
			// list all child event ids before or equal to the last occurrence
			childEvents, err := d.repo.ListEvents(ctx, authAccount, event.Parent, 0, 0, filter, []string{model.EventField_EventId})
			if err != nil {
				log.Error().Err(err).Msg("unable to list child events")
				return model.Event{}, domain.ErrInternal{Msg: "unable to list child events"}
			}
			if len(childEvents) > 0 {
				childEventIds := make([]model.EventId, len(childEvents))
				for i, childEvent := range childEvents {
					childEventIds[i] = model.EventId{EventId: childEvent.Id.EventId}
				}
				err = d.repo.BulkDeleteEvents(ctx, childEventIds)
				if err != nil {
					log.Error().Err(err).Msg("unable to bulk delete child events")
					return model.Event{}, domain.ErrInternal{Msg: "unable to bulk delete child events"}
				}
			}
		}
	}

	dbEvent, err = d.repo.UpdateEvent(ctx, authAccount, event, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update event")
		return model.Event{}, domain.ErrInternal{Msg: "unable to update event"}
	}

	return dbEvent, nil
}
