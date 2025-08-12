package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
)

const (
	EventTable = "event"
)

const (
	EventField_EventId            = "event_id"
	EventField_ParentEventId      = "parent_event_id"
	EventField_OverridenStartTime = "overriden_start_time"
	EventField_EventDataId        = "event_data_id"
	EventField_RecurrenceRule     = "recurrence_rule"
	EventField_ExcludedDates      = "excluded_dates"
	EventField_AdditionalDates    = "additional_dates"
	EventField_StartTime          = "start_time"
	EventField_EndTime            = "end_time"
	EventField_IsAllDay           = "is_all_day"
)

var EventFieldMasker = fieldmask.NewSQLFieldMasker(Event{}, map[string][]fieldmask.Field{
	model.EventField_EventId:            {{Name: EventField_EventId, Table: EventTable}},
	model.EventField_Parent:             {{Name: EventField_ParentEventId, Table: EventTable}},
	model.EventField_OverridenStartTime: {{Name: EventField_OverridenStartTime, Table: EventTable}},
	fieldmask.AlwaysIncludeKey: {
		{Name: EventField_EventDataId, Table: EventTable},
		{Name: EventField_ParentEventId, Table: EventTable},
	},
	model.EventField_RecurrenceRule:  {{Name: EventField_RecurrenceRule, Table: EventTable, Updatable: true}},
	model.EventField_ExcludedDates:   {{Name: EventField_ExcludedDates, Table: EventTable, Updatable: true}},
	model.EventField_AdditionalDates: {{Name: EventField_AdditionalDates, Table: EventTable, Updatable: true}},
	model.EventField_StartTime:       {{Name: EventField_StartTime, Table: EventTable, Updatable: true}},
	model.EventField_EndTime:         {{Name: EventField_EndTime, Table: EventTable, Updatable: true}},
	model.EventField_IsAllDay:        {{Name: EventField_IsAllDay, Table: EventTable, Updatable: true}},
})

var EventSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"start_time": {Name: EventField_StartTime, Table: EventTable},
	"end_time":   {Name: EventField_EndTime, Table: EventTable},
}, true)

// Event is the GORM model for an event.
// This table contains unique fields that distinguish between parent events and instances.
type Event struct {
	EventId            int64      `gorm:"primaryKey;autoIncrement;<-:false"`
	ParentEventId      *int64     `gorm:"index"`          // NULL for parent events, points to parent for instances
	OverridenStartTime *time.Time `gorm:"index"`          // NULL for parent events, set for instances/overrides
	EventDataId        int64      `gorm:"not null;index"` // Foreign key to event_data table

	RecurrenceRule  *string     `gorm:"column:recurrence_rule"`                  // only set for parent events
	ExcludedDates   []time.Time `gorm:"column:excluded_dates;serializer:json"`   // only set for parent events
	AdditionalDates []time.Time `gorm:"column:additional_dates;serializer:json"` // only set for parent events

	StartTime time.Time  `gorm:"column:start_time;not null;index"`
	EndTime   *time.Time `gorm:"column:end_time;index"`
	IsAllDay  bool       `gorm:"column:is_all_day;not null;default:false"`
}

// TableName sets the table name for the Event model.
func (Event) TableName() string {
	return EventTable
}
