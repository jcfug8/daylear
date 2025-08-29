package caldav

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emersion/go-ical"
	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/core/icalendar"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) CalendarGet(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarGet called")

	// Parse path parameters
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	calendarIDStr := vars["calendarID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in CalendarPropFind")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	calendarID, err := strconv.ParseInt(calendarIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse calendarID in CalendarPropFind")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID != authAccount.AuthUserId {
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in CalendarPropFind")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	calendar, err := s.domain.GetCalendar(r.Context(), authAccount, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get calendar in CalendarGet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	events, err := s.domain.ListEvents(r.Context(), authAccount, model.EventParent{CalendarId: calendarID}, 0, 0, "delete_time = null", []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to list events in CalendarGet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c := icalendar.ToICalendar(calendar, events)
	var buf bytes.Buffer
	if err := ical.NewEncoder(&buf).Encode(c); err != nil {
		s.log.Error().Err(err).Msg("Failed to encode calendar in CalendarGet")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/calendar")
	w.Header().Set("ETag", fmt.Sprintf("%d", calendar.UpdateTime.UTC().UnixNano()))
	w.Header().Set("Last-Modified", calendar.UpdateTime.UTC().Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
