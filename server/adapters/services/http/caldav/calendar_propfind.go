package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// Calendar specific property structures
type CalendarProp struct {
	ResourceType                  *ResourceType                  `xml:"D:resourcetype,omitempty"`
	DisplayName                   string                         `xml:"D:displayname,omitempty"`
	GetETag                       int64                          `xml:"D:getetag,omitempty"`
	GetCTag                       int64                          `xml:"CS:getctag,omitempty"`
	GetLastModified               string                         `xml:"D:getlastmodified,omitempty"`
	CalendarDescription           string                         `xml:"C:calendar-description,omitempty"`
	SupportedCalendarComponentSet *SupportedCalendarComponentSet `xml:"C:supported-calendar-component-set,omitempty"`
	SupportedCalendarData         *SupportedCalendarData         `xml:"C:supported-calendar-data,omitempty"`
	SupportedReportSet            *SupportedReportSet            `xml:"D:supported-report-set,omitempty"`
	CurrentUserPrivilegeSet       *PrivilegeSet                  `xml:"D:current-user-privilege-set,omitempty"`
	SyncToken                     string                         `xml:"D:sync-token,omitempty"`
	GetContentType                string                         `xml:"D:getcontenttype,omitempty"`
	// Timezone?
	Raw []RawXMLValue `xml:",any"`
}

type CalendarPropNames struct {
	ResourceType                  *struct{} `xml:"D:resourcetype,omitempty"`
	DisplayName                   *struct{} `xml:"D:displayname,omitempty"`
	GetETag                       *struct{} `xml:"D:getetag,omitempty"`
	GetCTag                       *struct{} `xml:"CS:getctag,omitempty"`
	GetLastModified               *struct{} `xml:"D:getlastmodified,omitempty"`
	CalendarDescription           *struct{} `xml:"C:calendar-description,omitempty"`
	SupportedCalendarComponentSet *struct{} `xml:"C:supported-calendar-component-set,omitempty"`
	SupportedCalendarData         *struct{} `xml:"C:supported-calendar-data,omitempty"`
	SupportedReportSet            *struct{} `xml:"D:supported-report-set,omitempty"`
	CurrentUserPrivilegeSet       *struct{} `xml:"D:current-user-privilege-set,omitempty"`
	SyncToken                     *struct{} `xml:"D:sync-token,omitempty"`
	GetContentType                *struct{} `xml:"D:getcontenttype,omitempty"`
}

func (s *Service) CalendarPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarPropFind called")

	depthStr := r.Header.Get("Depth")
	if depthStr == "infinity" {
		depthStr = "1"
	}

	depth, err := strconv.Atoi(depthStr)
	if err != nil {
		s.log.Error().Err(err).Str("depth", depthStr).Msg("Invalid Depth header in Calendars")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if depth > 1 {
		depth = 1
	}

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

	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var responses []Response

	// Build the response based on what was requested
	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeProp:
		responses, err = s.buildCalendarPropResponse(r.Context(), authAccount, calendarID, propFindRequest.Prop, depth)
	case PropFindRequestTypeAllProp:
		responses, err = s.buildCalendarAllPropResponse(r.Context(), authAccount, calendarID, depth)
	case PropFindRequestTypePropName:
		responses, err = s.buildCalendarPropNameResponse(r.Context(), authAccount, calendarID, depth)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
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

	s.log.Info().Msgf("calendar propfind response: %s", string(responseBytes))

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(responseBytes)
}

func (s *Service) buildCalendarPropResponse(ctx context.Context, authAccount model.AuthAccount, calendarID int64, prop *Prop, depth int) ([]Response, error) {
	// Get calendar from domain
	calendar, err := s.domain.GetCalendar(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain")
		return nil, err
	}

	return s._buildCalendarPropResponse(ctx, authAccount, calendar, prop, depth)
}

func (s *Service) _buildCalendarPropResponse(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, prop *Prop, depth int) ([]Response, error) {
	var foundP CalendarProp
	var notFoundP CalendarProp

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Calendar:   &Calendar{},
				Collection: &Collection{},
			}

		case raw.XMLName.Local == "displayname":
			foundP.DisplayName = calendar.Title

		case raw.XMLName.Local == "getetag":
			foundP.GetETag = calendar.UpdateTime.UTC().UnixNano()

		case raw.XMLName.Local == "getctag":
			foundP.GetCTag = calendar.EventUpdateTime.UTC().UnixNano()

		case raw.XMLName.Local == "getlastmodified":
			foundP.GetLastModified = calendar.UpdateTime.UTC().Format(time.RFC1123)

		case raw.XMLName.Local == "calendar-description":
			foundP.CalendarDescription = calendar.Description

		case raw.XMLName.Local == "supported-calendar-component-set":
			foundP.SupportedCalendarComponentSet = &SupportedCalendarComponentSet{
				CalendarComponents: []CalendarComponent{
					{Name: "VEVENT"},
				},
			}

		case raw.XMLName.Local == "supported-calendar-data":
			foundP.SupportedCalendarData = &SupportedCalendarData{
				CalendarData: []CalendarData{
					{ContentType: "text/calendar", Version: "2.0"},
				},
			}

		case raw.XMLName.Local == "supported-report-set":
			foundP.SupportedReportSet = &SupportedReportSet{
				SupportedReports: []SupportedReport{
					{Report: CalendarReportType{CalendarQuery: &CalendarQuery{}}},
					{Report: CalendarReportType{CalendarMultiget: &CalendarMultiget{}}},
					{Report: CalendarReportType{SyncCollection: &SyncCollection{}}},
				},
			}

		case raw.XMLName.Local == "current-user-privilege-set":
			privileges := []Privilege{
				{Name: "D:read"},
			}
			if calendar.CalendarAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE {
				privileges = append(privileges, Privilege{Name: "D:write"})
				privileges = append(privileges, Privilege{Name: "D:write-acl"})
			}
			foundP.CurrentUserPrivilegeSet = &PrivilegeSet{Privileges: privileges}

		case raw.XMLName.Local == "sync-token":
			foundP.SyncToken = calendar.EventUpdateTime.UTC().Format(time.RFC3339Nano)

		case raw.XMLName.Local == "getcontenttype":
			foundP.GetContentType = "text/calendar; charset=utf-8"

		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	calendarPath := s.formatCalendarPath(authAccount.AuthUserId, calendar.CalendarId.CalendarId)

	response := Response{Href: calendarPath}
	builder := ResponseBuilder{}

	// Add propstat for found properties
	if hasAnyCalendarPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}
	// Add propstat for not found properties
	if hasAnyCalendarPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	events, err := s.domain.ListEvents(ctx, authAccount, model.EventParent{UserId: authAccount.AuthUserId, CalendarId: calendar.CalendarId.CalendarId}, 0, 0, "delete_time = null", []string{})
	if err != nil {
		return []Response{}, err
	}

	eventResponses, err := s._buildEventPropResponse(ctx, authAccount, events, prop)
	if err != nil {
		return []Response{}, err
	}
	responses = append(responses, eventResponses...)

	return responses, nil
}

func (s *Service) buildCalendarAllPropResponse(ctx context.Context, authAccount model.AuthAccount, calendarID int64, depth int) ([]Response, error) {
	// Get calendar from domain
	calendar, err := s.domain.GetCalendar(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, model.CalendarId{CalendarId: calendarID}, []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendarID).Msg("Failed to get calendar from domain")
		return nil, err
	}

	return s._buildCalendarAllPropResponse(ctx, authAccount, calendar, depth)
}

func (s *Service) _buildCalendarAllPropResponse(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, depth int) ([]Response, error) {
	privileges := []Privilege{
		{Name: "D:read"},
	}
	if calendar.CalendarAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		privileges = append(privileges, Privilege{Name: "D:write"})
		privileges = append(privileges, Privilege{Name: "D:write-acl"})
	}

	foundP := CalendarProp{
		ResourceType: &ResourceType{
			Calendar:   &Calendar{},
			Collection: &Collection{},
		},
		DisplayName:     calendar.Title,
		GetETag:         calendar.UpdateTime.UTC().UnixNano(),
		GetCTag:         calendar.EventUpdateTime.UTC().UnixNano(),
		GetLastModified: calendar.UpdateTime.UTC().Format(time.RFC1123),
		// CalendarDescription: calendar.Description, Should not be returned by an allProp request via RFC4791
		// SupportedCalendarComponentSet: &SupportedCalendarComponentSet{ Should not be returned by an allProp request via RFC4791
		// 	CalendarComponents: []CalendarComponent{
		// 		{Name: "VEVENT"},
		// 	},
		// },
		// SupportedCalendarData: &SupportedCalendarData{ Should not be returned by an allProp request via RFC4791
		// 	CalendarData: []CalendarData{
		// 		{ContentType: "text/calendar", Version: "2.0"},
		// 	},
		// },
		SupportedReportSet: &SupportedReportSet{
			SupportedReports: []SupportedReport{
				{Report: CalendarReportType{CalendarQuery: &CalendarQuery{}}},
				{Report: CalendarReportType{CalendarMultiget: &CalendarMultiget{}}},
				{Report: CalendarReportType{SyncCollection: &SyncCollection{}}},
			},
		},
		CurrentUserPrivilegeSet: &PrivilegeSet{Privileges: privileges},
		SyncToken:               calendar.EventUpdateTime.UTC().Format(time.RFC3339Nano),
		GetContentType:          "text/calendar; charset=utf-8",
	}

	calendarPath := s.formatCalendarPath(authAccount.AuthUserId, calendar.CalendarId.CalendarId)

	response := Response{Href: calendarPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	events, err := s.domain.ListEvents(ctx, authAccount, model.EventParent{UserId: authAccount.AuthUserId, CalendarId: calendar.CalendarId.CalendarId}, 0, 0, "delete_time = null", []string{})
	if err != nil {
		return []Response{}, err
	}

	for _, event := range events {
		eventResponses, err := s._buildEventAllPropResponse(ctx, authAccount, event)
		if err != nil {
			return []Response{}, err
		}
		responses = append(responses, eventResponses...)
	}

	return responses, nil
}

func (s *Service) buildCalendarPropNameResponse(ctx context.Context, authAccount model.AuthAccount, calendarID int64, depth int) ([]Response, error) {
	calendarPath := s.formatCalendarPath(authAccount.AuthUserId, calendarID)

	response := Response{Href: calendarPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, CalendarPropNames{
		ResourceType:                  &struct{}{},
		DisplayName:                   &struct{}{},
		GetETag:                       &struct{}{},
		GetCTag:                       &struct{}{},
		GetLastModified:               &struct{}{},
		CalendarDescription:           &struct{}{},
		SupportedCalendarComponentSet: &struct{}{},
		SupportedCalendarData:         &struct{}{},
		SupportedReportSet:            &struct{}{},
		CurrentUserPrivilegeSet:       &struct{}{},
	}, 200)

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	events, err := s.domain.ListEvents(ctx, authAccount, model.EventParent{UserId: authAccount.AuthUserId, CalendarId: calendarID}, 0, 0, "delete_time = null", []string{})
	if err != nil {
		return []Response{}, err
	}

	for _, event := range events {
		eventResponses, err := s.buildEventPropNameResponse(ctx, authAccount, calendarID, event.Id.EventId)
		if err != nil {
			return []Response{}, err
		}
		responses = append(responses, eventResponses...)
	}

	return responses, nil
}

func hasAnyCalendarPropProperties(prop CalendarProp) bool {
	return prop.ResourceType != nil ||
		prop.DisplayName != "" ||
		prop.GetETag != 0 ||
		prop.GetCTag != 0 ||
		prop.GetLastModified != "" ||
		prop.CalendarDescription != "" ||
		(prop.SupportedCalendarComponentSet != nil && len(prop.SupportedCalendarComponentSet.CalendarComponents) > 0) ||
		(prop.SupportedCalendarData != nil && len(prop.SupportedCalendarData.CalendarData) > 0) ||
		(prop.SupportedReportSet != nil && len(prop.SupportedReportSet.SupportedReports) > 0) ||
		(prop.CurrentUserPrivilegeSet != nil && len(prop.CurrentUserPrivilegeSet.Privileges) > 0) ||
		len(prop.Raw) > 0
}

func (s *Service) formatCalendarPath(userID, calendarID int64) string {
	return path.Join(s.apiPath, fmt.Sprintf("/caldav/principals/%d/calendars/%d/", userID, calendarID))
}

func (s *Service) parseCalendarPath(path string) (int64, int64, error) {
	s.log.Info().Msgf("Parsing calendar path: %s", path)
	path = strings.TrimPrefix(path, s.apiPath)
	path = strings.TrimSuffix(path, "/")
	s.log.Info().Msgf("Trimmed calendar path by removing api path %s: %s", s.apiPath, path)
	parts := strings.Split(path, "/")
	if len(parts) != 6 || parts[2] != "principals" || parts[4] != "calendars" {
		return 0, 0, fmt.Errorf("invalid calendar path")
	}
	userId, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	calendarId, err := strconv.ParseInt(parts[5], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return userId, calendarId, nil
}
