package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/emersion/go-webdav/caldav"
	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/icalendar"
	"github.com/jcfug8/daylear/server/core/model"
)

type CalendarReportRequest struct {
	SyncCollection *SyncCollectionRequest `xml:"sync-collection,omitempty"`
}

type SyncCollectionRequest struct{}

func (b *Service) Calendar(w http.ResponseWriter, r *http.Request) {
	b.log.Info().Msg("Calendar called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to parse auth data in Calendar")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		setCalDAVHeaders(w)
		w.Header().Set("Allow", "PROPFIND,OPTIONS,REPORT")
		w.WriteHeader(http.StatusNoContent)
		return
	case "PROPFIND":
		b.CalendarPropFind(w, r, authAccount)
		return
	case "REPORT":
		b.CalendarReport(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) CalendarPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarPropFind called")

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

	_, err = strconv.ParseInt(calendarIDStr, 10, 64)
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

	// // Get calendar from domain
	// calendar, err := s.domain.GetCalendar(r.Context(), model.AuthAccount{AuthUserId: userID}, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	// if err != nil {
	// 	s.log.Error().Err(err).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// // Create response for calendar properties
	// response := CalendarsMultiStatus{
	// 	// XMLNSD:  "DAV:",
	// 	// XMLNSC:  "urn:ietf:params:xml:ns:caldav",
	// 	// XMLNSCS: "http://calendarserver.org/ns/",
	// 	// Response: []CalendarsResponse{
	// 	// 	{
	// 	// 		Href: r.URL.Path,
	// 	// 		Propstat: CalendarsPropstat{
	// 	// 			Prop: CalendarsProp{
	// 	// 				ResourceType: ResourceType{
	// 	// 					Calendar:   &Calendar{},
	// 	// 					Collection: &Collection{},
	// 	// 				},
	// 	// 				GetETag:             calendar.UpdateTime.UnixNano(),
	// 	// 				GetCTag:             calendar.UpdateTime.UnixNano(),
	// 	// 				GetLastModified:     calendar.UpdateTime.Format(time.RFC1123),
	// 	// 				DisplayName:         calendar.Title,
	// 	// 				CalendarDescription: calendar.Description,
	// 	// 				SupportedCalendarComponentSet: &SupportedCalendarComponentSet{
	// 	// 					CalendarComponents: []CalendarComponent{
	// 	// 						{Name: "VEVENT"},
	// 	// 					},
	// 	// 				},
	// 	// 				SupportedCalendarData: &SupportedCalendarData{
	// 	// 					CalendarData: []CalendarData{
	// 	// 						{ContentType: "text/calendar", Version: "2.0"},
	// 	// 					},
	// 	// 				},
	// 	// 				SupportedReportSet: &SupportedReportSet{
	// 	// 					SupportedReports: []SupportedReport{
	// 	// 						{Report: CalendarReportType{CalendarQuery: &CalendarQuery{}}},
	// 	// 						{Report: CalendarReportType{CalendarMultiget: &CalendarMultiget{}}},
	// 	// 						{Report: CalendarReportType{SyncCollection: &SyncCollection{}}},
	// 	// 					},
	// 	// 				},
	// 	// 			},
	// 	// 			Status: Status{
	// 	// 				Status: "HTTP/1.1 200 OK",
	// 	// 			},
	// 	// 		},
	// 	// 	},
	// 	// },
	// }

	// responseBytes, err := xml.MarshalIndent(response, "", "  ")
	// if err != nil {
	// 	s.log.Error().Err(err).Msg("Failed to marshal response in CalendarPropFind")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// setCalDAVHeaders(w)
	// w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	// w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	// w.WriteHeader(http.StatusMultiStatus)
	// w.Write(addXMLDeclaration(responseBytes))
}

func (s *Service) CalendarReport(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarReport called")

	// Parse path parameters
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	calendarIDStr := vars["calendarID"]

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in CalendarReport")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	calendarID, err := strconv.ParseInt(calendarIDStr, 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse calendarID in CalendarReport")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userID != authAccount.AuthUserId {
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in CalendarReport")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get calendar from domain
	calendar, err := s.domain.GetCalendar(r.Context(), model.AuthAccount{AuthUserId: userID}, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to read body in CalendarReport")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reportRequest := CalendarReportRequest{}
	err = xml.Unmarshal(body, &reportRequest)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to unmarshal body in CalendarReport")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if reportRequest.SyncCollection != nil {
		s.CalendarReportSyncCollection(w, r, authAccount, calendar)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (s *Service) CalendarReportSyncCollection(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount, calendar model.Calendar) {
	// List events from domain
	events, err := s.domain.ListEvents(r.Context(), authAccount, model.EventParent{CalendarId: calendar.CalendarId.CalendarId, UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendar.CalendarId.CalendarId).Msg("Failed to list events from domain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert to iCalendar format
	icalendar.ToICalendar(calendar, events)

	// Set CalDAV headers
	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/calendar; charset=utf-8")
	// w.Header().Set("Content-Length", strconv.Itoa(len(icalData.Bytes())))
	w.WriteHeader(http.StatusOK)
	// w.Write(icalData.Bytes())
}

// CreateCalendar creates a new calendar
func (b *Service) CreateCalendar(ctx context.Context, calendar *caldav.Calendar) error {
	b.log.Info().Str("calendarName", calendar.Name).Msg("CreateCalendar called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Msg("Failed to parse auth data in CreateCalendar")
		return err
	}

	b.log.Info().Int64("userID", authAccount.AuthUserId).Str("calendarName", calendar.Name).Msg("Creating calendar")
	_, err = b.domain.CreateCalendar(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.Calendar{Title: calendar.Name, Description: calendar.Description})
	if err != nil {
		b.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Str("calendarName", calendar.Name).Msg("Failed to create calendar")
	} else {
		b.log.Info().Int64("userID", authAccount.AuthUserId).Str("calendarName", calendar.Name).Msg("Calendar created successfully")
	}
	return err
}

// GetCalendar returns a specific calendar
func (b *Service) GetCalendar(ctx context.Context, path string) (*caldav.Calendar, error) {
	b.log.Info().Str("path", path).Msg("GetCalendar called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse auth data in GetCalendar")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("authUserID", authAccount.AuthUserId).Msg("Parsing calendar path")
	userID, calendarID, err := parseCalendarPath(path)
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Msg("Failed to parse calendar path")
		return nil, err
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Int64("authUserID", authAccount.AuthUserId).Msg("Parsed calendar path")

	if userID != authAccount.AuthUserId {
		b.log.Error().Str("path", path).Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("User ID mismatch in GetCalendar")
		return nil, fmt.Errorf("user ID does not match")
	}

	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Getting calendar from domain")
	calendar, err := b.domain.GetCalendar(ctx, model.AuthAccount{AuthUserId: userID}, model.CalendarParent{UserId: userID}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		b.log.Error().Err(err).Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain")
		return nil, err
	}

	caldavPath := formatCalendarPath(authAccount.AuthUserId, calendar.CalendarId.CalendarId)
	b.log.Info().Str("path", path).Int64("userID", userID).Int64("calendarID", calendarID).Str("caldavPath", caldavPath).Str("title", calendar.Title).Msg("GetCalendar returning")

	return &caldav.Calendar{
		Path:                  caldavPath,
		Name:                  calendar.Title,
		Description:           calendar.Description,
		MaxResourceSize:       1048576,
		SupportedComponentSet: []string{"VEVENT"},
	}, nil
}

func parseCalendarPath(path string) (userID, calendarID int64, err error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 3 || parts[0] != "calendars" {
		return 0, 0, fmt.Errorf("invalid calendar path: %s", path)
	}

	userID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	calendarID, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return userID, calendarID, nil
}

func formatCalendarPath(userID, calendarID int64) string {
	return fmt.Sprintf("/calendars/%d/%d/", userID, calendarID)
}
