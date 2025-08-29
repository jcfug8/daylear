package caldav

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/core/icalendar"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) EventGet(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("EventGet called")

	// Parse path parameters
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	calendarIDStr := vars["calendarID"]
	eventIDStr := strings.TrimSuffix(vars["eventID"], ".ics")

	// Parse user ID
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in EventGet")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse calendar ID
	calendarID, err := strconv.ParseInt(calendarIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse calendarID in EventGet")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse event ID
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse eventID in EventGet")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Verify user owns this calendar
	if userID != authAccount.AuthUserId {
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in EventGet")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get the event from the domain
	event, err := s.domain.GetEvent(r.Context(), authAccount, model.EventParent{UserId: userID, CalendarId: calendarID}, model.EventId{EventId: eventID}, nil)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get event in EventGet")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Check if event is deleted
	if event.DeleteTime != nil {
		s.log.Error().Msg("Event is deleted in EventGet")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Convert to iCalendar format
	cal := icalendar.ToICalendar(model.Calendar{}, []model.Event{event})
	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		s.log.Error().Err(err).Msg("Failed to encode event to iCalendar in EventGet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	w.Header().Set("ETag", fmt.Sprintf("\"%d\"", event.UpdateTime.UTC().UnixNano()))
	w.Header().Set("Last-Modified", event.UpdateTime.UTC().Format(time.RFC1123))

	// Return the event
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
