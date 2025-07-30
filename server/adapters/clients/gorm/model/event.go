package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/masks"
)

// EventMap maps the Event fields to their corresponding fields in the core model.
var EventMap = masks.NewFieldMap().
	MapFieldToFields("event_id",
		EventFields.EventId).
	MapFieldToFields("parent_event_id",
		EventFields.ParentEventId).
	MapFieldToFields("recurrence_time",
		EventFields.RecurrenceTime).
	MapFieldToFields("event_data_id",
		EventFields.EventDataId)

// EventFields defines the event fields in the GORM model.
var EventFields = eventFields{
	EventId:        "event.event_id",
	ParentEventId:  "event.parent_event_id",
	RecurrenceTime: "event.recurrence_time",
	EventDataId:    "event.event_data_id",
}

type eventFields struct {
	EventId        string
	ParentEventId  string
	RecurrenceTime string
	EventDataId    string
}

// Map maps the event fields to their corresponding model values.
func (fields eventFields) Map(m Event) map[string]any {
	return map[string]any{
		fields.EventId:        m.EventId,
		fields.ParentEventId:  m.ParentEventId,
		fields.RecurrenceTime: m.RecurrenceTime,
		fields.EventDataId:    m.EventDataId,
	}
}

// Mask returns a FieldMask for the event fields.
func (fields eventFields) Mask() []string {
	return []string{
		fields.EventId,
		fields.ParentEventId,
		fields.RecurrenceTime,
		fields.EventDataId,
	}
}

// Event is the GORM model for an event.
// This table contains unique fields that distinguish between parent events and instances.
type Event struct {
	EventId        int64      `gorm:"primaryKey;autoIncrement;<-:false"`
	ParentEventId  *int64     `gorm:"index"`          // NULL for parent events, points to parent for instances
	RecurrenceTime *time.Time `gorm:"index"`          // NULL for parent events, set for instances/overrides
	EventDataId    int64      `gorm:"not null;index"` // Foreign key to event_data table

	RecurrenceRule  *string     `gorm:"column:recurrence_rule"`                  // only set for parent events
	ExcludedDates   []time.Time `gorm:"column:excluded_dates;serializer:json"`   // only set for parent events
	AdditionalDates []time.Time `gorm:"column:additional_dates;serializer:json"` // only set for parent events
}

// TableName sets the table name for the Event model.
func (Event) TableName() string {
	return "event"
}

// IsRecurringSet returns true if this event represents a recurring event set
func (e Event) IsRecurringSet() bool {
	return e.ParentEventId == nil && e.RecurrenceTime == nil
}

// IsInstance returns true if this event represents a specific instance
func (e Event) IsInstance() bool {
	return e.ParentEventId != nil || e.RecurrenceTime != nil
}

// IsOverride returns true if this event is an override of a recurring event
func (e Event) IsOverride() bool {
	return e.ParentEventId != nil && e.RecurrenceTime != nil
}
