package model

import (
	"fmt"
	"time"

	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"github.com/teambition/rrule-go"
)

// Event represents a VEVENT component in iCalendar.
// This can be either a single event or a recurring event with exceptions.
type EventSet struct {
	Parent EventSetParent
	Id     int64 `aip_pattern:"event_set"`

	// RecurrenceRule defines how this event repeats (if recurring)
	RecurrenceRule  *string
	ExcludedDates   []time.Time
	AdditionalDates []time.Time

	EventData
}

type EventSetParent struct {
	UserId   int64 `aip_pattern:"user"`
	CircleId int64 `aip_pattern:"circle"`
}

type EventInstance struct {
	Parent EventInstanceParent
	Id     int64 `aip_pattern:"event_instance"`

	// RecurrenceId indicates this is an exception instance of a recurring event
	// If set, this event overrides the regular instance at this date/time
	RecurrenceId *time.Time

	EventData
}

type EventInstanceParent struct {
	EventSetParent
	EventSetId int64 `aip_pattern:"event_set"`
}

type EventData struct {
	// DTStamp is the creation timestamp of the event
	CreateTime *time.Time
	// UpdateTime is the last update timestamp of the event
	UpdateTime *time.Time
	// DTStart is the start date/time of the event
	// Required for all events
	StartTime time.Time
	// EndTime is the end date/time of the event
	// Optional - if not specified, event is assumed to be all-day
	EndTime *time.Time
	// IsAllDay indicates if this is an all-day event
	IsAllDay bool
	// Title is the title/summary of the event
	Title *string
	// Description is the detailed description of the event
	Description *string
	// Location is the location of the event
	Location *string
	// Status indicates the status of the event
	Status pb.Event_State
	// Class indicates the classification of the event
	Class pb.Event_Class
	// URL is a URL associated with the event
	URL *string

	// Alarms are the VALARM components associated with this event
	Alarms []*Alarm
}

// GenerateInstances generates a list of event instances based on the event's RecurrenceRule ExcludedDates AdditionalDates.
// It should generate instance within now and now + duration.
func (r EventSet) GenerateInstances(startingTime, endingTime time.Time) ([]EventInstance, error) {
	var eventDuration time.Duration
	if r.EndTime != nil {
		eventDuration = r.EndTime.Sub(r.StartTime)
	}
	var instances []EventInstance

	// If no recurrence rule, this is a single event
	if r.RecurrenceRule == nil || *r.RecurrenceRule == "" {
		// Check if the event falls within the specified duration
		if r.StartTime.After(startingTime) && r.StartTime.Before(endingTime) {
			instance := EventInstance{
				Parent: EventInstanceParent{
					EventSetParent: r.Parent,
					EventSetId:     r.Id,
				},
				EventData: r.EventData,
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
		instance := EventInstance{
			Parent: EventInstanceParent{
				EventSetParent: r.Parent,
				EventSetId:     r.Id,
			},
			EventData: r.EventData,
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
