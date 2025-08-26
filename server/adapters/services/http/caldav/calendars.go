package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/emersion/go-webdav/caldav"
	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

// CalDAV XML response structures
type CalendarsMultiStatus struct {
	XMLName  xml.Name            `xml:"D:multistatus"`
	XMLNSD   string              `xml:"xmlns:D,attr"`
	XMLNSC   string              `xml:"xmlns:C,attr"`
	XMLNSCS  string              `xml:"xmlns:CS,attr"`
	Response []CalendarsResponse `xml:"D:response"`
}

type CalendarsResponse struct {
	Href     string            `xml:"D:href"`
	Propstat CalendarsPropstat `xml:"D:propstat"`
}

type CalendarsPropstat struct {
	Prop   CalendarsProp `xml:"D:prop"`
	Status Status        `xml:"D:status"`
}

type CalendarsProp struct {
	ResourceType                  ResourceType                   `xml:"D:resourcetype"`
	DisplayName                   string                         `xml:"D:displayname"`
	GetETag                       int64                          `xml:"D:getetag,omitempty"`
	GetCTag                       int64                          `xml:"CS:getctag,omitempty"`
	GetLastModified               string                         `xml:"D:getlastmodified,omitempty"`
	CalendarDescription           string                         `xml:"C:calendar-description,omitempty"`
	SupportedCalendarComponentSet *SupportedCalendarComponentSet `xml:"C:supported-calendar-component-set,omitempty"`
	SupportedCalendarData         *SupportedCalendarData         `xml:"C:supported-calendar-data,omitempty"`
	SupportedReportSet            *SupportedReportSet            `xml:"D:supported-report-set,omitempty"`
}

type SupportedCalendarComponentSet struct {
	CalendarComponents []CalendarComponent `xml:"C:comp,omitempty"`
}

type CalendarComponent struct {
	Name string `xml:"name,attr,omitempty"`
}

type SupportedCalendarData struct {
	CalendarData []CalendarData `xml:"C:calendar-data,omitempty"`
}

type CalendarData struct {
	ContentType string `xml:"content-type,attr,omitempty"`
	Version     string `xml:"version,attr,omitempty"`
}

type SupportedReportSet struct {
	SupportedReports []SupportedReport `xml:"D:supported-report"`
}

type SupportedReport struct {
	Report CalendarReportType `xml:"D:report"`
}

type CalendarReportType struct {
	CalendarQuery    *CalendarQuery    `xml:"C:calendar-query,omitempty"`
	CalendarMultiget *CalendarMultiget `xml:"C:calendar-multiget,omitempty"`
	SyncCollection   *SyncCollection   `xml:"D:sync-collection,omitempty"`
}

type CalendarQuery struct{}
type CalendarMultiget struct{}
type SyncCollection struct{}

func (s *Service) Calendars(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Calendars called")

	userID, err := strconv.ParseInt(mux.Vars(r)["userID"], 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in Calendars")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Calendars")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userID != authAccount.AuthUserId {
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in UserPrincipal")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.CalendarsOptions(w, r)
		return
	case "PROPFIND":
		s.CalendarsPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) CalendarsOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CalendarsPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	calendars, err := s.domain.ListCalendars(r.Context(), model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := CalendarsMultiStatus{
		XMLNSD:  "DAV:",
		XMLNSC:  "urn:ietf:params:xml:ns:caldav",
		XMLNSCS: "http://calendarserver.org/ns/",
		Response: []CalendarsResponse{
			{
				Href: r.URL.Path,
				Propstat: CalendarsPropstat{
					Prop: CalendarsProp{
						ResourceType: ResourceType{
							Collection: &Collection{},
						},
						DisplayName: "Calendars",
					},
					Status: Status{
						Status: "HTTP/1.1 200 OK",
					},
				},
			},
		},
	}

	for _, calendar := range calendars {
		response.Response = append(response.Response, CalendarsResponse{
			Href: fmt.Sprintf("/caldav/principals/%d/calendars/%d", authAccount.AuthUserId, calendar.CalendarId.CalendarId),
			Propstat: CalendarsPropstat{
				Prop: CalendarsProp{
					ResourceType: ResourceType{
						Calendar:   &Calendar{},
						Collection: &Collection{},
					},
					DisplayName:         calendar.Title,
					GetETag:             calendar.UpdateTime.UnixNano(),
					GetCTag:             calendar.UpdateTime.UnixNano(),
					GetLastModified:     calendar.UpdateTime.Format(time.RFC1123),
					CalendarDescription: calendar.Description,
					SupportedCalendarComponentSet: &SupportedCalendarComponentSet{
						CalendarComponents: []CalendarComponent{
							{Name: "VEVENT"},
						},
					},
					SupportedCalendarData: &SupportedCalendarData{
						CalendarData: []CalendarData{
							{ContentType: "text/calendar", Version: "2.0"},
						},
					},
					SupportedReportSet: &SupportedReportSet{
						SupportedReports: []SupportedReport{
							{Report: CalendarReportType{CalendarQuery: &CalendarQuery{}}},
							{Report: CalendarReportType{CalendarMultiget: &CalendarMultiget{}}},
							{Report: CalendarReportType{SyncCollection: &SyncCollection{}}},
						},
					},
				},
				Status: Status{
					Status: "HTTP/1.1 200 OK",
				},
			},
		})
	}

	responseBytes, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in CalendarsPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(addXMLDeclaration(responseBytes))
}

func (s *Service) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	return nil, nil
}
