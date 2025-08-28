package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) CalendarReport(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarReport called")

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

	reportRequest, err := NewReportRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse REPORT request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var responses []Response

	switch reportRequest.GetRequestType() {
	case ReportRequestTypeCalendarQuery:
		responses, err = s.buildCalendarQueryResponse()
	case ReportRequestTypeCalendarMultiget:
		responses, err = s.buildCalendarMultigetResponse(r.Context(), authAccount, calendarID, reportRequest.CalendarMultiget)
	case ReportRequestTypeSyncCollection:
		responses, err = s.buildSyncCollectionResponse(r.Context(), authAccount, calendarID, reportRequest.SyncCollection)
	default:
		s.log.Error().Msg("Invalid REPORT request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build calendar response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	multistatus := ResponseBuilder{}.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in CalendarPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes = addXMLDeclaration(responseBytes)

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(responseBytes)
}

func (s *Service) buildCalendarQueryResponse() ([]Response, error) {
	return []Response{}, nil
}

func (s *Service) buildCalendarMultigetResponse(ctx context.Context, authAccount model.AuthAccount, calendarID int64, calendarMultiget *CalendarMultigetReport) ([]Response, error) {
	var filter string
	var hrefCalendarId int64
	eventIds := []string{}
	if len(calendarMultiget.Hrefs) == 0 {
		return []Response{}, fmt.Errorf("no hrefs provided")
	}

	for _, href := range calendarMultiget.Hrefs {
		userId, calId, eventId, err := s.parseCalendarMultigetHref(href)
		if err != nil {
			return []Response{}, err
		}
		if userId != authAccount.AuthUserId || calId != calendarID {
			return []Response{}, fmt.Errorf("invalid user id or calendar id in event path")
		}
		if eventId == 0 && calId != 0 {
			if calId != calendarID {
				return []Response{}, fmt.Errorf("invalid calendar id in event path")
			}
			hrefCalendarId = calendarID
		} else {
			eventIds = append(eventIds, strconv.FormatInt(eventId, 10))
		}
	}
	responses := []Response{}

	if hrefCalendarId != 0 {
		calendarResponses, err := s.buildCalendarPropResponse(ctx, authAccount, hrefCalendarId, calendarMultiget.Prop, 0)
		if err != nil {
			return []Response{}, err
		}
		responses = append(responses, calendarResponses...)
	}

	filter = fmt.Sprintf("any(event_id,%s)", strings.Join(eventIds, ","))

	events, err := s.domain.ListEvents(ctx, authAccount, model.EventParent{UserId: authAccount.AuthUserId, CalendarId: calendarID}, 0, 0, filter, []string{})
	if err != nil {
		return []Response{}, err
	}

	for _, event := range events {
		eventResponses, err := s._buildEventPropResponse(ctx, authAccount, event, calendarMultiget.Prop)
		if err != nil {
			return []Response{}, err
		}
		responses = append(responses, eventResponses...)
	}

	return responses, nil
}

func (s *Service) buildSyncCollectionResponse(ctx context.Context, authAccount model.AuthAccount, calendarID int64, syncCollection *SyncCollectionReport) ([]Response, error) {
	filter := ""
	if syncCollection.SyncToken != nil && *syncCollection.SyncToken != "" {
		// turn nano into time.Time
		syncTime, err := time.Parse(time.RFC3339Nano, *syncCollection.SyncToken)
		if err != nil {
			return []Response{}, err
		}
		filter = fmt.Sprintf("update_time >= '%s' OR delete_time >= '%s'", syncTime.UTC().Format(time.RFC3339), syncTime.UTC().Format(time.RFC3339))
	}

	events, err := s.domain.ListEvents(ctx, authAccount, model.EventParent{UserId: authAccount.AuthUserId, CalendarId: calendarID}, 0, 0, filter, []string{})
	if err != nil {
		return []Response{}, err
	}

	responses := []Response{}

	for _, event := range events {
		eventResponses, err := s._buildEventPropResponse(ctx, authAccount, event, syncCollection.Prop)
		if err != nil {
			return []Response{}, err
		}
		responses = append(responses, eventResponses...)
	}

	calendar, err := s.domain.GetCalendar(ctx, authAccount, model.CalendarParent{UserId: authAccount.AuthUserId}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		return []Response{}, err
	}

	calendarProp := CalendarProp{
		SyncToken: calendar.EventUpdateTime.UTC().Format(time.RFC3339Nano),
	}

	calendarPath := s.formatCalendarPath(authAccount.AuthUserId, calendarID)
	response := Response{Href: calendarPath}
	builder := ResponseBuilder{}

	response = builder.AddPropertyStatus(response, calendarProp, 200)
	responses = append(responses, response)

	return responses, nil
}

func (s *Service) parseCalendarMultigetHref(path string) (int64, int64, int64, error) {
	s.log.Info().Msgf("Parsing event path: %s", path)
	path = strings.TrimPrefix(path, s.apiPath)
	s.log.Info().Msgf("Trimmed event path by removing api path %s: %s", s.apiPath, path)
	parts := strings.Split(path, "/")
	if len(parts) == 8 && parts[2] == "principals" && parts[4] == "calendars" && parts[6] == "events" {
		return s.parseEventPath(path)
	} else if len(parts) == 7 && parts[2] == "principals" && parts[4] == "calendars" {
		userId, calendarId, err := s.parseCalendarPath(path)
		return userId, calendarId, 0, err
	}
	return 0, 0, 0, fmt.Errorf("invalid calendar multiget href")

}
