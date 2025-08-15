package model

import (
	"fmt"
	"time"

	"github.com/teambition/rrule-go"
)

const (
	EventField_Parent             = "parent"
	EventField_EventId            = "event_id"
	EventField_RecurrenceRule     = "recurrence_rule"
	EventField_ExcludedDates      = "excluded_dates"
	EventField_AdditionalDates    = "additional_dates"
	EventField_ParentEventId      = "parent_event_id"
	EventField_OverridenStartTime = "overriden_start_time"
	EventField_CreateTime         = "create_time"
	EventField_UpdateTime         = "update_time"
	EventField_StartTime          = "start_time"
	EventField_EndTime            = "end_time"
	EventField_IsAllDay           = "is_all_day"
	EventField_Title              = "title"
	EventField_Description        = "description"
	EventField_Location           = "location"
	EventField_Geo                = "geo"
	EventField_URL                = "url"
	EventField_Alarms             = "alarms"
	EventField_RecurrenceEndTime  = "recurrence_end_time"
)

// Event represents a VEVENT component in iCalendar.
// This can be either a single event or a recurring event with exceptions.
type Event struct {
	Parent EventParent
	Id     EventId

	// these are set if the event is part of a recurring event set
	// if they are set, and the ParentEventId is not set, then
	// you treat this as the parent event of a recurring event set
	// and you should get the instances or generate the instances
	// from the recurrence rule
	RecurrenceRule  *string
	ExcludedDates   []time.Time
	AdditionalDates []time.Time
	// this is set if the event is an instance of a recurring event set
	// it points to the parent event id of the recurring event set
	ParentEventId *int64

	// this is set if the event is an override of an event in
	// the recurring event set
	OverridenStartTime *time.Time

	CreateTime        time.Time
	UpdateTime        time.Time
	StartTime         time.Time
	EndTime           *time.Time
	IsAllDay          bool
	Title             string
	Description       string
	Location          string
	Geo               LatLng
	URL               string
	RecurrenceEndTime *time.Time

	Alarms []*Alarm
}

type LatLng struct {
	Latitude  float64
	Longitude float64
}

type EventParent struct {
	UserId     int64 `aip_pattern:"key=user"`
	CircleId   int64 `aip_pattern:"key=circle"`
	CalendarId int64 `aip_pattern:"key=calendar"`
}

type EventId struct {
	// the event id
	EventId int64 `aip_pattern:"key=event"`
}

// GenerateClones generates a list of event instances based on the event's RecurrenceRule ExcludedDates AdditionalDates.
// It should generate instance within now and now + duration.
func (r Event) GenerateClones(startingTime, endingTime time.Time) ([]Event, error) {
	var eventDuration time.Duration
	if r.EndTime != nil {
		eventDuration = r.EndTime.Sub(r.StartTime)
	}
	var instances []Event

	// If no recurrence rule, this is a single event
	if r.RecurrenceRule == nil || *r.RecurrenceRule == "" {
		return instances, nil
	}

	recurrenceRule := "RRULE:" + *r.RecurrenceRule

	// Parse the recurrence rule
	rule, err := rrule.StrToRRuleSet(recurrenceRule)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recurrence rule: %w", err)
	}

	// Set the start time for the rule
	rule.DTStart(r.StartTime)
	rule.SetExDates(r.ExcludedDates)
	rule.SetRDates(r.AdditionalDates)

	// Generate occurrences within the specified duration
	occurrences := rule.Between(startingTime, endingTime, true)

	// Process each occurrence
	for _, occurrence := range occurrences {
		if occurrence.Unix() == r.StartTime.Unix() {
			continue
		}

		instance := Event{
			Parent:        r.Parent,
			ParentEventId: &r.Id.EventId,
		}
		instance.StartTime = occurrence
		if r.EndTime != nil {
			endTime := occurrence.Add(eventDuration)
			instance.EndTime = &endTime
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

func (r Event) GetUntil() *time.Time {
	if r.RecurrenceRule == nil || *r.RecurrenceRule == "" {
		return r.EndTime
	}

	rr, err := rrule.StrToRRuleSet("RRULE:" + *r.RecurrenceRule)
	if err != nil {
		return r.EndTime
	}

	until := rr.GetRRule().GetUntil()
	if until.IsZero() {
		return nil
	}

	return &until
}
