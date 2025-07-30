package model

import (
	"fmt"
	"time"

	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"github.com/teambition/rrule-go"
)

// Event represents a VEVENT component in iCalendar.
// This can be either a single event or a recurring event with exceptions.
type Event struct {
	Parent EventParent
	Id     EventId

	// these are set if the event is part of a recurring event set
	// if they are set, and the RecurringEventId is not set, then
	// you treat this as the parent event of a recurring event set
	// and you should get the instances or generate the instances
	// from the recurrence rule
	RecurrenceRule  *string
	ExcludedDates   []time.Time
	AdditionalDates []time.Time
	// this is set if the event is an instance of a recurring event set
	// it points to the parent event id of the recurring event set
	RecurringEventId int64

	// this is set if the event is an override of an event in
	// the recurring event set
	RecurrenceTime *time.Time

	CreateTime  *time.Time
	UpdateTime  *time.Time
	StartTime   time.Time
	EndTime     *time.Time
	IsAllDay    bool
	Title       *string
	Description *string
	Location    *string
	Status      pb.Event_State
	Class       pb.Event_Class
	URL         *string

	Alarms []*Alarm
}

type EventParent struct {
	UserId     int64 `aip_pattern:"user"`
	CircleId   int64 `aip_pattern:"circle"`
	CalendarId int64 `aip_pattern:"calendar"`
}

type EventId struct {
	// the event id
	EventId int64 `aip_pattern:"event"`
}

// GenerateInstances generates a list of event instances based on the event's RecurrenceRule ExcludedDates AdditionalDates.
// It should generate instance within now and now + duration.
func (r Event) GenerateInstances(startingTime, endingTime time.Time) ([]Event, error) {
	var eventDuration time.Duration
	if r.EndTime != nil {
		eventDuration = r.EndTime.Sub(r.StartTime)
	}
	var instances []Event

	// If no recurrence rule, this is a single event
	if r.RecurrenceRule == nil || *r.RecurrenceRule == "" {
		// Check if the event falls within the specified duration
		if r.StartTime.After(startingTime) && r.StartTime.Before(endingTime) {
			instance := Event{
				Parent: r.Parent,
				Id:     r.Id,
			}
			instances = append(instances, instance)
		}
		return instances, nil
	}

	// Parse the recurrence rule
	rule, err := rrule.StrToRRuleSet(*r.RecurrenceRule)
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
		instance := Event{
			Parent:           r.Parent,
			RecurringEventId: r.RecurringEventId,
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
