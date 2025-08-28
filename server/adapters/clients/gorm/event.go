package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateEvent creates a new event in the database
func (c *Client) CreateEvent(ctx context.Context, event model.Event, fields []string) (model.Event, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Strs("fields", fields).
		Logger()

	mEvent, mEventData, err := convert.EventFromCoreModel(event)
	if err != nil {
		log.Error().Err(err).Msg("invalid event when creating event row")
		return model.Event{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid event: %v", err)}
	}

	now := time.Now().UTC()

	mEventData.CreateTime = &now
	mEventData.UpdateTime = &now

	res := c.db.WithContext(ctx).
		Select(gmodel.EventDataFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&mEventData)
	if res.Error != nil {
		log.Error().Err(err).Msg("unable to create event data row")
		return model.Event{}, ConvertGormError(err)
	}

	mEvent.EventDataId = mEventData.EventDataId

	eventFields := gmodel.EventFieldMasker.Convert(fields)

	res = c.db.WithContext(ctx).
		Select(eventFields).
		Clauses(clause.Returning{}).
		Create(&mEvent)
	if res.Error != nil {
		log.Error().Err(err).Msg("unable to create event row")
		return model.Event{}, ConvertGormError(err)
	}

	res = c.db.WithContext(ctx).
		Select(gmodel.CalendarColumn_EventUpdateTime).
		Where("calendar_id = ?", mEventData.CalendarId).
		Updates(&gmodel.Calendar{
			EventUpdateTime: now,
		})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update calendar row")
		return model.Event{}, ConvertGormError(res.Error)
	}

	event, err = convert.EventToCoreModel(mEvent, mEventData)
	if err != nil {
		log.Error().Err(err).Msg("invalid event row when creating event")
		return model.Event{}, fmt.Errorf("unable to read event: %v", err)
	}

	return event, nil
}

// DeleteEvent deletes an event from the database
func (c *Client) DeleteEvent(ctx context.Context, id model.EventId) (model.Event, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Int64("event_id", id.EventId).
		Logger()

	var event gmodel.Event
	deleteTime := time.Now().UTC()

	mEventData := gmodel.EventData{
		DeleteTime: &deleteTime,
	}

	eventRes := c.db.WithContext(ctx).
		Where("event_id = ?", id.EventId).
		Clauses(clause.Returning{}).
		Find(&event)

	if eventRes.Error != nil {
		log.Error().Err(eventRes.Error).Msg("unable to delete event")
		return model.Event{}, ConvertGormError(eventRes.Error)
	}

	eventDataRes := c.db.WithContext(ctx).
		Where("event_data_id = ?", event.EventDataId).
		Clauses(clause.Returning{}).
		Updates(&mEventData)

	if eventDataRes.Error != nil {
		log.Error().Err(eventDataRes.Error).Msg("unable to delete event data")
		return model.Event{}, ConvertGormError(eventDataRes.Error)
	}

	res := c.db.WithContext(ctx).
		Select(gmodel.CalendarColumn_EventUpdateTime).
		Where("calendar_id = ?", mEventData.CalendarId).
		Updates(&gmodel.Calendar{
			EventUpdateTime: time.Now().UTC(),
		})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update calendar row")
		return model.Event{}, ConvertGormError(res.Error)
	}

	m, err := convert.EventToCoreModel(event, mEventData)
	if err != nil {
		log.Error().Err(err).Msg("invalid event row when deleting event")
		return model.Event{}, fmt.Errorf("unable to read event: %v", err)
	}

	return m, nil
}

// DeleteChildEvents deletes all events that have the given parent event id from the database
func (c *Client) DeleteChildEvents(ctx context.Context, id model.EventId) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Int64("event_id", id.EventId).
		Logger()

	var events []gmodel.Event
	var eventDataIds []int64

	eventRes := c.db.WithContext(ctx).
		Where("parent_event_id = ?", id.EventId).
		Clauses(clause.Returning{}).
		Find(&events)

	if eventRes.Error != nil {
		log.Error().Err(eventRes.Error).Msg("unable to delete child events")
		return ConvertGormError(eventRes.Error)
	}

	for _, event := range events {
		eventDataIds = append(eventDataIds, event.EventDataId)
	}

	deleteTime := time.Now()

	eventDataRes := c.db.WithContext(ctx).
		Where("event_data_id IN (?)", eventDataIds).
		Updates(&gmodel.EventData{
			DeleteTime: &deleteTime,
		})
	if eventDataRes.Error != nil {
		log.Error().Err(eventDataRes.Error).Msg("unable to delete event data")
		return ConvertGormError(eventDataRes.Error)
	}

	mEventData := gmodel.EventData{}
	res := c.db.WithContext(ctx).
		Joins("JOIN event ON event.event_data_id = event_data.event_data_id").
		Where("event.event_id = ?", id.EventId).
		First(&mEventData)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to get event data")
		return ConvertGormError(res.Error)
	}

	res = c.db.WithContext(ctx).
		Select(gmodel.CalendarColumn_EventUpdateTime).
		Where("calendar_id = ?", mEventData.CalendarId).
		Updates(&gmodel.Calendar{
			EventUpdateTime: time.Now().UTC(),
		})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update calendar row")
		return ConvertGormError(res.Error)
	}

	return nil
}

// BulkDeleteEvents deletes a list of events from the database
func (c *Client) BulkDeleteEvents(ctx context.Context, ids []model.EventId) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Int("num_events", len(ids)).
		Logger()

	var events []gmodel.Event
	var eventIds []int64
	var eventDataIds []int64

	for _, id := range ids {
		eventIds = append(eventIds, id.EventId)
	}

	eventRes := c.db.WithContext(ctx).
		Where("event_id IN (?)", eventIds).
		Clauses(clause.Returning{}).
		Find(&events)

	if eventRes.Error != nil {
		log.Error().Err(eventRes.Error).Msg("unable to delete child events")
		return ConvertGormError(eventRes.Error)
	}

	for _, event := range events {
		eventDataIds = append(eventDataIds, event.EventDataId)
	}

	deleteTime := time.Now()

	eventDataRes := c.db.WithContext(ctx).
		Where("event_data_id IN (?)", eventDataIds).
		Updates(&gmodel.EventData{
			DeleteTime: &deleteTime,
		})

	if eventDataRes.Error != nil {
		log.Error().Err(eventDataRes.Error).Msg("unable to delete event data")
		return ConvertGormError(eventDataRes.Error)
	}

	return nil
}

// GetEvent retrieves an event from the database
func (c *Client) GetEvent(ctx context.Context, authAccount model.AuthAccount, id model.EventId, fields []string) (model.Event, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Strs("fields", fields).
		Logger()

	type Result struct {
		gmodel.Event
		gmodel.EventData
	}

	var result Result

	fields = append(gmodel.EventFieldMasker.Convert(fields), gmodel.EventDataFieldMasker.Convert(fields)...)

	res := c.db.WithContext(ctx).
		Table(gmodel.EventTable).
		Select(fields).
		Joins("JOIN event_data ON event.event_data_id = event_data.event_data_id").
		Where("event.event_id = ?", id.EventId).
		First(&result)

	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to get event")
		return model.Event{}, ConvertGormError(res.Error)
	}

	event, err := convert.EventToCoreModel(result.Event, result.EventData)
	if err != nil {
		log.Error().Err(err).Msg("invalid event row when getting event")
		return model.Event{}, fmt.Errorf("unable to read event: %v", err)
	}
	return event, nil
}

// ListEvents lists events from the database with pagination and filtering
func (c *Client) ListEvents(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, pageSize int32, offset int64, filter string, fields []string) ([]model.Event, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Strs("fields", fields).
		Int64("calendar_id", parent.CalendarId).
		Int32("page_size", pageSize).
		Int64("offset", offset).
		Str("filter", filter).
		Logger()

	type Result struct {
		gmodel.Event
		gmodel.EventData
	}

	var results []Result

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "event.start_time"},
		Desc:   false,
	}}

	fields = append(gmodel.EventFieldMasker.Convert(fields), gmodel.EventDataFieldMasker.Convert(fields)...)

	tx := c.db.WithContext(ctx).
		Table(gmodel.EventTable).
		Select(fields).
		Joins("JOIN event_data ON event.event_data_id = event_data.event_data_id").
		Where("event_data.calendar_id = ?", parent.CalendarId).
		Order(clause.OrderBy{Columns: orders})

	if pageSize > 0 {
		tx = tx.Limit(int(pageSize))
	}

	if offset > 0 {
		tx = tx.Offset(int(offset))
	}

	conversion, err := gmodel.EventSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing event rows")
		return []model.Event{}, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&results).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list events")
		return []model.Event{}, ConvertGormError(err)
	}

	events := make([]model.Event, len(results))
	for i, result := range results {
		event, err := convert.EventToCoreModel(result.Event, result.EventData)
		if err != nil {
			log.Error().Err(err).Msg("invalid event row when listing events")
			return []model.Event{}, fmt.Errorf("unable to read event: %v", err)
		}
		events[i] = event
	}
	return events, nil
}

// UpdateEvent updates an existing event in the database
func (c *Client) UpdateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event, fields []string) (model.Event, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx).With().
		Strs("fields", fields).
		Logger()

	mEvent, mEventData, err := convert.EventFromCoreModel(event)
	if err != nil {
		log.Error().Err(err).Msg("invalid event when updating event row")
		return model.Event{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid event: %v", err)}
	}

	eventFields := gmodel.EventFieldMasker.Convert(fields, fieldmask.OnlyUpdatable())

	if len(eventFields) > 0 {
		res := c.db.WithContext(ctx).
			Select(eventFields).
			Where("event_id = ?", mEvent.EventId).
			Clauses(clause.Returning{}).
			Updates(&mEvent)
		if res.Error != nil {
			log.Error().Err(res.Error).Msg("unable to update event row")
			return model.Event{}, ConvertGormError(res.Error)
		}
	}

	if mEvent.EventDataId == 0 {
		res := c.db.WithContext(ctx).
			Select(gmodel.EventField_EventDataId).
			Where("event_id = ?", mEvent.EventId).
			First(&mEvent)
		if res.Error != nil {
			log.Error().Err(res.Error).Msg("unable to get event data row")
		}
	}

	dataFields := append(gmodel.EventDataFieldMasker.Convert(fields, fieldmask.OnlyUpdatable()), gmodel.EventDataField_UpdateTime)
	now := time.Now().UTC()
	mEventData.UpdateTime = &now

	res := c.db.WithContext(ctx).
		Select(dataFields).
		Where("event_data_id = ?", mEvent.EventDataId).
		Clauses(clause.Returning{}).
		Updates(&mEventData)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update event data row")
		return model.Event{}, ConvertGormError(res.Error)
	}

	res = c.db.WithContext(ctx).
		Select(gmodel.CalendarColumn_EventUpdateTime).
		Where("calendar_id = ?", mEventData.CalendarId).
		Updates(&gmodel.Calendar{
			EventUpdateTime: time.Now().UTC(),
		})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update calendar row")
		return model.Event{}, ConvertGormError(res.Error)
	}

	event, err = convert.EventToCoreModel(mEvent, mEventData)
	if err != nil {
		log.Error().Err(err).Msg("invalid event row when updating event")
		return model.Event{}, fmt.Errorf("unable to read event: %v", err)
	}

	return event, nil
}
