package caldav

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/jcfug8/daylear/server/core/icalendar"
	"github.com/jcfug8/daylear/server/core/model"
)

type EventProp struct {
	GetETag         int64         `xml:"D:getetag,omitempty"`
	GetLastModified string        `xml:"D:getlastmodified,omitempty"`
	CalendarData    string        `xml:"C:calendar-data,omitempty"`
	GetContentType  string        `xml:"D:getcontenttype,omitempty"`
	Raw             []RawXMLValue `xml:",any"`
}

type EventPropNames struct {
	GetETag         *struct{} `xml:"D:getetag,omitempty"`
	GetLastModified *struct{} `xml:"D:getlastmodified,omitempty"`
	CalendarData    *struct{} `xml:"C:calendar-data,omitempty"`
	GetContentType  *struct{} `xml:"D:getcontenttype,omitempty"`
}

func (s *Service) _buildEventPropResponse(_ context.Context, authAccount model.AuthAccount, event model.Event, prop *Prop) ([]Response, error) {
	var foundP EventProp
	var notFoundP EventProp

	eventPath := formatEventPath(authAccount.AuthUserId, event.Parent.CalendarId, event.Id.EventId)

	if event.DeleteTime != nil {
		return []Response{
			{
				Href:   eventPath,
				Status: &Status{Status: "HTTP/1.1 404 Not Found"},
			},
		}, nil
	}

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "getetag":
			foundP.GetETag = event.UpdateTime.UnixNano()
		case raw.XMLName.Local == "getlastmodified":
			foundP.GetLastModified = event.UpdateTime.Format(time.RFC1123)
		case raw.XMLName.Local == "calendar-data":
			cal := icalendar.ToICalendar(model.Calendar{}, []model.Event{event})
			var buf bytes.Buffer
			if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
				return []Response{}, err
			}
			foundP.CalendarData = buf.String()
		case raw.XMLName.Local == "getcontenttype":
			foundP.GetContentType = "text/calendar; charset=utf-8"
		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	response := Response{Href: eventPath}
	builder := ResponseBuilder{}

	if hasAnyEventPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}

	if hasAnyEventPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	return []Response{response}, nil
}

func (s *Service) _buildEventAllPropResponse(ctx context.Context, authAccount model.AuthAccount, event model.Event) ([]Response, error) {
	cal := icalendar.ToICalendar(model.Calendar{}, []model.Event{event})
	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		return []Response{}, err
	}

	foundP := EventProp{
		GetETag:         event.UpdateTime.UnixNano(),
		GetLastModified: event.UpdateTime.Format(time.RFC1123),
		CalendarData:    buf.String(),
		GetContentType:  "text/calendar; charset=utf-8",
	}

	eventPath := formatEventPath(authAccount.AuthUserId, event.Parent.CalendarId, event.Id.EventId)

	response := Response{Href: eventPath}
	builder := ResponseBuilder{}

	response = builder.AddPropertyStatus(response, foundP, 200)

	return []Response{response}, nil
}

func (s *Service) buildEventPropNameResponse(ctx context.Context, authAccount model.AuthAccount, calendarID, eventID int64) ([]Response, error) {
	eventPath := formatEventPath(authAccount.AuthUserId, calendarID, eventID)

	response := Response{Href: eventPath}
	builder := ResponseBuilder{}

	response = builder.AddPropertyStatus(response, EventPropNames{
		GetETag:         &struct{}{},
		GetLastModified: &struct{}{},
		CalendarData:    &struct{}{},
		GetContentType:  &struct{}{},
	}, 200)

	return []Response{response}, nil
}

func hasAnyEventPropProperties(prop EventProp) bool {
	return prop.GetETag != 0 ||
		prop.CalendarData != "" ||
		prop.GetContentType != "" ||
		prop.GetLastModified != "" ||
		len(prop.Raw) > 0
}

func formatEventPath(userID, calendarID, eventID int64) string {
	return fmt.Sprintf("/caldav/principals/%d/calendars/%d/events/%d.ics", userID, calendarID, eventID)
}

func parseEventPath(path string) (int64, int64, int64, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 8 {
		return 0, 0, 0, fmt.Errorf("invalid event path")
	}
	userId, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	calendarId, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	eventId, err := strconv.ParseInt(strings.TrimSuffix(parts[7], ".ics"), 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return userId, calendarId, eventId, nil
}
