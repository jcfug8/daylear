package icalendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/teambition/rrule-go"
)

// ToICalendar converts a calendar and its events to iCalendar format
func ToICalendar(cal model.Calendar, events []model.Event) *ical.Calendar {
	// Create new calendar
	calendar := ical.NewCalendar()

	// Set calendar properties using the correct API
	calendar.Props.SetText(ical.PropVersion, "2.0")
	calendar.Props.SetText(ical.PropProductID, "-//Daylear//Calendar//EN")
	// calendar.Props.SetText(ical.PropCalendarScale, "GREGORIAN")
	// calendar.Props.SetText(ical.PropMethod, "PUBLISH")

	// Add events
	for _, event := range events {
		eventComponent := eventToComponent(event)
		calendar.Children = append(calendar.Children, eventComponent)
	}

	return calendar
}

// FromICalendar converts iCalendar content to a calendar and events
func FromICalendar(calendar *ical.Calendar) (model.Calendar, []model.Event, error) {

	// Extract calendar properties
	cal := model.Calendar{
		Title:       getCalendarProperty(calendar, ical.PropName),
		Description: getCalendarProperty(calendar, ical.PropDescription),
		// Add other calendar properties as needed
	}

	// Extract events
	var events []model.Event
	for _, component := range calendar.Children {
		if component.Name == ical.CompEvent {
			event, err := componentToEvent(component)
			if err != nil {
				// Log error but continue processing other events
				continue
			}
			events = append(events, event)
		}
	}

	return cal, events, nil
}

// eventToComponent converts a model.Event to an ical.Component
func eventToComponent(event model.Event) *ical.Component {
	component := ical.NewComponent(ical.CompEvent)

	// Set event properties using the correct API
	component.Props.SetText(ical.PropUID, fmt.Sprintf("%d", event.Id.EventId))
	component.Props.SetDateTime(ical.PropDateTimeStamp, event.CreateTime)
	component.Props.SetDateTime(ical.PropDateTimeStart, event.StartTime)
	component.Props.SetDateTime(ical.PropDateTimeEnd, *event.EndTime)
	component.Props.SetText(ical.PropSummary, event.Title)

	if event.Description != "" {
		component.Props.SetText(ical.PropDescription, event.Description)
	}

	if event.Location != "" {
		component.Props.SetText(ical.PropLocation, event.Location)
	}

	// Set status using the correct constant
	component.Props.SetText(ical.PropStatus, string(ical.EventConfirmed))

	// TODO: handle sequence
	// Set sequence for conflict detection
	component.Props.Set(&ical.Prop{
		Name:  ical.PropSequence,
		Value: "1",
	})

	// Handle recurrence rules if you have them
	if event.RecurrenceRule != nil {
		if !strings.HasPrefix(*event.RecurrenceRule, "RRULE:") {
			*event.RecurrenceRule = "RRULE:" + *event.RecurrenceRule
		}
		rule, err := rrule.StrToRRuleSet(*event.RecurrenceRule)
		if err != nil {
			return nil
		}
		component.Props.SetRecurrenceRule(&rule.GetRRule().Options)
	}

	// TODO: handle alarms
	// Handle alarms if you have them
	// for _, alarm := range event.Alarms {
	// 	alarmComponent := ical.NewComponent(ical.CompAlarm)
	// 	alarmComponent.Props.SetText(ical.PropAction, "DISPLAY")
	// 	alarmComponent.Props.SetText(ical.PropDescription, *alarm.Description)
	// 	alarmComponent.Props.SetText(ical.PropSummary, *alarm.Summary)
	// 	alarmComponent.Props.SetText(ical.PropTrigger, *alarm.Trigger.Duration)

	// 	component.Children = append(component.Children, alarmComponent)
	// }

	return component
}

// componentToEvent converts an ical.Component to a model.Event
func componentToEvent(component *ical.Component) (model.Event, error) {
	event := model.Event{}

	// Extract UID
	if uid := component.Props.Get(ical.PropUID); uid != nil {
		uidText, err := uid.Text()
		if err == nil {
			// Parse UID to extract event ID
			if strings.Contains(uidText, "@daylear.com") {
				uidParts := strings.Split(uidText, "@")
				if len(uidParts) > 0 {
					// You might want to parse this to get the actual event ID
					// For now, we'll use a placeholder
				}
			}
		}
	}

	// Extract summary
	if summary := component.Props.Get(ical.PropSummary); summary != nil {
		if summaryText, err := summary.Text(); err == nil {
			event.Title = summaryText
		}
	}

	// Extract description
	if description := component.Props.Get(ical.PropDescription); description != nil {
		if descText, err := description.Text(); err == nil {
			event.Description = descText
		}
	}

	// Extract location
	if location := component.Props.Get(ical.PropLocation); location != nil {
		if locText, err := location.Text(); err == nil {
			event.Location = locText
		}
	}

	// Extract start time
	if startTime := component.Props.Get(ical.PropDateTimeStart); startTime != nil {
		if start, err := startTime.DateTime(time.UTC); err == nil {
			event.StartTime = start
		}
	}

	// Extract end time
	if endTime := component.Props.Get(ical.PropDateTimeEnd); endTime != nil {
		if end, err := endTime.DateTime(time.UTC); err == nil {
			event.EndTime = &end
		}
	}

	// Extract recurrence rule
	if rrule := component.Props.Get(ical.PropRecurrenceRule); rrule != nil {
		if rruleText, err := rrule.Text(); err == nil {
			event.RecurrenceRule = &rruleText
		}
	}

	// Set creation time to now if not specified
	if event.CreateTime.IsZero() {
		event.CreateTime = time.Now()
	}

	// Set update time to now
	event.UpdateTime = time.Now()

	return event, nil
}

// getCalendarProperty safely extracts a property value from a calendar
func getCalendarProperty(calendar *ical.Calendar, propName string) string {
	if prop := calendar.Props.Get(propName); prop != nil {
		if text, err := prop.Text(); err == nil {
			return text
		}
	}
	return ""
}
