package icalendar

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ical "github.com/arran4/golang-ical"
	"github.com/jcfug8/daylear/server/core/model"
)

// ToICalendar converts a Calendar model to iCalendar format
func ToICalendar(cal *model.Calendar, events []model.Event) string {
	icalCal := ical.NewCalendar()
	icalCal.SetMethod(ical.MethodPublish)
	icalCal.SetProductId("-//Daylear//Calendar//EN")
	icalCal.SetCalscale("GREGORIAN")

	// Add calendar properties
	if cal.Title != "" {
		icalCal.SetXWRCalName(cal.Title)
	}
	if cal.Description != "" {
		icalCal.SetXWRCalDesc(cal.Description)
	}

	// Add events
	for _, event := range events {
		icalEvent := toICalEvent(event)

		icalCal.AddVEvent(icalEvent)
	}

	return icalCal.Serialize()
}

func toICalEvent(event model.Event) *ical.VEvent {
	icalEvent := ical.NewEvent(fmt.Sprintf("%d", event.Id.EventId))

	// Set recurrence rule
	if event.RecurrenceRule != nil {
		icalEvent.AddRrule(*event.RecurrenceRule)
	}

	// Set excluded dates
	for _, date := range event.ExcludedDates {
		icalEvent.AddExdate(date.Format("20060102T150405Z"))
	}

	// Set additional dates
	for _, date := range event.AdditionalDates {
		icalEvent.AddRdate(date.Format("20060102T150405Z"))
	}

	if event.OverridenStartTime != nil {
		icalEvent.SetProperty(ical.ComponentProperty(ical.PropertyRecurrenceId), event.OverridenStartTime.Format("20060102T150405Z"))
	}

	// Set creation time
	if event.CreateTime.IsZero() {
		icalEvent.SetCreatedTime(event.CreateTime)
	} else {
		icalEvent.SetCreatedTime(time.Now())
	}

	// Set start time
	if event.IsAllDay {
		icalEvent.SetAllDayStartAt(event.StartTime)
	} else {
		icalEvent.SetStartAt(event.StartTime)
	}

	// Set end time
	if event.EndTime != nil {
		if event.IsAllDay {
			icalEvent.SetAllDayEndAt(*event.EndTime)
		} else {
			icalEvent.SetEndAt(*event.EndTime)
		}
	}

	// Set summary
	if event.Title != "" {
		icalEvent.SetSummary(event.Title)
	}

	// Set description
	if event.Description != "" {
		icalEvent.SetDescription(event.Description)
	}

	// Set location
	if event.Location.Latitude != 0 && event.Location.Longitude != 0 {
		icalEvent.SetLocation(fmt.Sprintf("%f,%f", event.Location.Latitude, event.Location.Longitude))
	}

	// Set URL
	if event.URL != "" {
		icalEvent.SetURL(event.URL)
	}
	// Add alarms
	for _, alarm := range event.Alarms {
		icalAlarm := icalEvent.AddAlarm()
		icalAlarm.SetProperty(ical.ComponentProperty(ical.PropertyAction), "DISPLAY")

		// Set alarm UID
		icalAlarm.SetProperty(ical.ComponentProperty(ical.PropertyUid), alarm.AlarmId)

		// Set creation time
		if alarm.CreateTime != nil {
			icalAlarm.SetCreatedTime(*alarm.CreateTime)
		}

		// Set description
		if alarm.Description != nil {
			icalAlarm.SetDescription(*alarm.Description)
		}

		// Set summary
		if alarm.Summary != nil {
			icalAlarm.SetSummary(*alarm.Summary)
		}

		// Set trigger
		if alarm.Trigger != nil {
			if alarm.Trigger.Duration != nil {
				icalAlarm.SetProperty(ical.ComponentProperty(ical.PropertyTrigger), alarm.Trigger.Duration.String())
			} else if alarm.Trigger.DateTime != nil {
				icalAlarm.SetProperty(ical.ComponentProperty(ical.PropertyTrigger), "VALUE=DATE-TIME:"+alarm.Trigger.DateTime.Format("20060102T150405Z"))
			}
		}
	}

	return icalEvent
}

// FromICalendar converts iCalendar format to a Calendar model
func FromICalendar(icalString string) (model.Calendar, []model.Event, error) {
	cal, err := ical.ParseCalendar(strings.NewReader(icalString))
	if err != nil {
		return model.Calendar{}, nil, err
	}

	result := model.Calendar{}
	events := []model.Event{}

	// Parse calendar properties
	for _, prop := range cal.CalendarProperties {
		switch prop.IANAToken {
		case "X-WR-CALNAME":
			result.Title = prop.Value
		case "X-WR-CALDESC":
			result.Description = prop.Value
		}
	}

	// Parse events
	for _, icalEvent := range cal.Events() {
		event := model.Event{}

		// Parse UID
		uidId := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyUid))
		if uidId == nil {
			continue
		}
		event.Id.EventId, err = strconv.ParseInt(uidId.Value, 10, 64)
		if err != nil {
			continue
		}

		// Parse creation time
		if created := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyCreated)); created != nil {
			if t, err := time.Parse("20060102T150405Z", created.Value); err == nil {
				event.CreateTime = t
			}
		}

		// Parse start time
		if start := icalEvent.GetProperty(ical.ComponentProperty("DTSTART")); start != nil {
			if t, err := time.Parse("20060102T150405Z", start.Value); err == nil {
				event.StartTime = t
				// Check if it's all-day based on format
				if start.GetValueType() == ical.ValueDataTypeDate {
					event.IsAllDay = true
				}
			}
		}

		// Parse end time
		if end := icalEvent.GetProperty(ical.ComponentProperty("DTEND")); end != nil {
			if t, err := time.Parse("20060102T150405Z", end.Value); err == nil {
				event.EndTime = &t
			}
		}

		// Parse summary
		if summary := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertySummary)); summary != nil {
			event.Title = summary.Value
		}

		// Parse description
		if description := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyDescription)); description != nil {
			event.Description = description.Value
		}

		// Parse location
		if location := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyLocation)); location != nil {
			splitLatLng := strings.Split(location.Value, ",")
			if len(splitLatLng) == 2 {
				latitude, err := strconv.ParseFloat(splitLatLng[0], 64)
				if err == nil {
					longitude, err := strconv.ParseFloat(splitLatLng[1], 64)
					if err == nil {
						event.Location = model.LatLng{
							Latitude:  latitude,
							Longitude: longitude,
						}
					}
				}
			}
		}

		// Parse URL
		if url := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyUrl)); url != nil {
			event.URL = url.Value
		}

		// Parse recurrence ID
		if rid := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyRecurrenceId)); rid != nil {
			if t, err := time.Parse("20060102T150405Z", rid.Value); err == nil {
				event.OverridenStartTime = &t
			}
			events = append(events, event)
		} else {
			// Parse recurrence rule
			if rrule := icalEvent.GetProperty(ical.ComponentProperty(ical.PropertyRrule)); rrule != nil {
				event.RecurrenceRule = &rrule.Value
			}

			// Parse excluded dates
			for _, exdate := range icalEvent.GetProperties(ical.ComponentProperty(ical.PropertyExdate)) {
				if t, err := time.Parse("20060102T150405Z", exdate.Value); err == nil {
					event.ExcludedDates = append(event.ExcludedDates, t)
				}
			}

			// Parse additional dates
			for _, rdate := range icalEvent.GetProperties(ical.ComponentProperty(ical.PropertyRdate)) {
				if t, err := time.Parse("20060102T150405Z", rdate.Value); err == nil {
					event.AdditionalDates = append(event.AdditionalDates, t)
				}
			}
			events = append(events, event)
		}

		// Parse alarms
		for _, icalAlarm := range icalEvent.Alarms() {
			alarm := &model.Alarm{}

			// Parse alarm UID
			if uid := icalAlarm.GetProperty(ical.ComponentProperty(ical.PropertyUid)); uid != nil {
				alarm.AlarmId = uid.Value
			}

			// Parse creation time
			if created := icalAlarm.GetProperty(ical.ComponentProperty(ical.PropertyCreated)); created != nil {
				if t, err := time.Parse("20060102T150405Z", created.Value); err == nil {
					alarm.CreateTime = &t
				}
			}

			// Parse description
			if description := icalAlarm.GetProperty(ical.ComponentProperty(ical.PropertyDescription)); description != nil {
				descStr := description.Value
				alarm.Description = &descStr
			}

			// Parse summary
			if summary := icalAlarm.GetProperty(ical.ComponentProperty(ical.PropertySummary)); summary != nil {
				summaryStr := summary.Value
				alarm.Summary = &summaryStr
			}

			// Parse trigger
			if trigger := icalAlarm.GetProperty(ical.ComponentProperty(ical.PropertyTrigger)); trigger != nil {
				if trigger.GetValueType() == ical.ValueDataTypeDuration {
					// Relative trigger
					if duration, err := parseDuration(trigger.Value); err == nil {
						alarm.Trigger = &model.Trigger{
							Duration: &duration,
						}
					}
				} else {
					// Absolute trigger
					if t, err := time.Parse("20060102T150405Z", trigger.Value); err == nil {
						alarm.Trigger = &model.Trigger{
							DateTime: &t,
						}
					}
				}
			}

			event.Alarms = append(event.Alarms, alarm)
		}
	}

	return result, events, nil
}

// Helper function to parse duration string (e.g., "-PT15M")
func parseDuration(value string) (time.Duration, error) {
	// Remove leading minus sign if present
	isNegative := false
	if strings.HasPrefix(value, "-") {
		isNegative = true
		value = value[1:]
	}

	// Parse duration
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, err
	}

	if isNegative {
		duration = -duration
	}

	return duration, nil
}

// Helper functions for formatting
func formatIntSlice(ints []int) string {
	var parts []string
	for _, i := range ints {
		parts = append(parts, strconv.Itoa(i))
	}
	return strings.Join(parts, ",")
}

func formatWeekdaySlice(weekdays []time.Weekday) string {
	var parts []string
	for _, wd := range weekdays {
		parts = append(parts, formatWeekday(wd))
	}
	return strings.Join(parts, ",")
}

func formatWeekday(wd time.Weekday) string {
	switch wd {
	case time.Sunday:
		return "SU"
	case time.Monday:
		return "MO"
	case time.Tuesday:
		return "TU"
	case time.Wednesday:
		return "WE"
	case time.Thursday:
		return "TH"
	case time.Friday:
		return "FR"
	case time.Saturday:
		return "SA"
	default:
		return "MO"
	}
}

// Helper functions for parsing
func parseIntSlice(value string) []int {
	var result []int
	parts := strings.Split(value, ",")
	for _, part := range parts {
		if i, err := strconv.Atoi(part); err == nil {
			result = append(result, i)
		}
	}
	return result
}

func parseWeekdaySlice(value string) []time.Weekday {
	var result []time.Weekday
	parts := strings.Split(value, ",")
	for _, part := range parts {
		if wd := parseWeekday(part); wd != nil {
			result = append(result, *wd)
		}
	}
	return result
}

func parseWeekday(value string) *time.Weekday {
	switch strings.ToUpper(value) {
	case "SU":
		wd := time.Sunday
		return &wd
	case "MO":
		wd := time.Monday
		return &wd
	case "TU":
		wd := time.Tuesday
		return &wd
	case "WE":
		wd := time.Wednesday
		return &wd
	case "TH":
		wd := time.Thursday
		return &wd
	case "FR":
		wd := time.Friday
		return &wd
	case "SA":
		wd := time.Saturday
		return &wd
	default:
		return nil
	}
}
