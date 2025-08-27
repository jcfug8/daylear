package caldav

import (
	"bytes"
	"context"
	"fmt"
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
