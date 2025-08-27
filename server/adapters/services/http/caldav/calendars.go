package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// Calendar specific property structures
type CalendarProperties struct {
	// used for both calendar and calendar collection
	ResourceType *ResourceType `xml:"D:resourcetype,omitempty"`
	// used for both calendar and calendar collection
	DisplayName                   string                         `xml:"D:displayname,omitempty"`
	GetETag                       int64                          `xml:"D:getetag,omitempty"`
	GetCTag                       int64                          `xml:"CS:getctag,omitempty"`
	GetLastModified               string                         `xml:"D:getlastmodified,omitempty"`
	CalendarDescription           string                         `xml:"C:calendar-description,omitempty"`
	SupportedCalendarComponentSet *SupportedCalendarComponentSet `xml:"C:supported-calendar-component-set,omitempty"`
	SupportedCalendarData         *SupportedCalendarData         `xml:"C:supported-calendar-data,omitempty"`
	SupportedReportSet            *SupportedReportSet            `xml:"D:supported-report-set,omitempty"`
	CurrentUserPrivilegeSet       *PrivilegeSet                  `xml:"D:current-user-privilege-set,omitempty"`
	Raw                           []RawXMLValue                  `xml:",any"`
}

func (c *CalendarProperties) ContainsRaw(raw RawXMLValue) bool {
	for _, r := range c.Raw {
		if r.XMLName.Local == raw.XMLName.Local {
			return true
		}
	}
	return false
}

type CalendarPropertyNames struct {
	// used for both calendar and calendar collection
	ResourceType *struct{} `xml:"D:resourcetype,omitempty"`
	// used for both calendar and calendar collection
	DisplayName                   *struct{} `xml:"D:displayname,omitempty"`
	GetETag                       *struct{} `xml:"D:getetag,omitempty"`
	GetCTag                       *struct{} `xml:"CS:getctag,omitempty"`
	GetLastModified               *struct{} `xml:"D:getlastmodified,omitempty"`
	CalendarDescription           *struct{} `xml:"C:calendar-description,omitempty"`
	SupportedCalendarComponentSet *struct{} `xml:"C:supported-calendar-component-set,omitempty"`
	SupportedCalendarData         *struct{} `xml:"C:supported-calendar-data,omitempty"`
	SupportedReportSet            *struct{} `xml:"D:supported-report-set,omitempty"`
	CurrentUserPrivilegeSet       *struct{} `xml:"D:current-user-privilege-set,omitempty"`
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
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in Calendars")
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
	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var foundProps map[string]CalendarProperties
	var notFoundProps map[string]CalendarProperties
	var propNames map[string]CalendarPropertyNames

	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeProp:
		foundProps, notFoundProps, err = s.buildCalendarsPropResponse(r.Context(), authAccount, propFindRequest.Prop)
	case PropFindRequestTypeAllProp:
		foundProps, err = s.buildCalendarsAllPropResponse(r.Context(), authAccount)
	case PropFindRequestTypePropName:
		propNames, err = s.buildCalendarsPropNameResponse(r.Context(), authAccount)
	}

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build calendars response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responses := []Response{}
	builder := ResponseBuilder{}

	for href, foundP := range foundProps {
		response := Response{Href: href}
		// Add propstat for found properties
		if hasAnyCalendarProperties(foundP) {
			response = builder.AddPropertyStatus(response, foundP, 200)
		}
		// Add propstat for not found properties
		if notFoundP, ok := notFoundProps[href]; ok && hasAnyCalendarProperties(notFoundP) {
			response = builder.AddPropertyStatus(response, notFoundP, 404)
		}
		responses = append(responses, response)
	}

	for href, pNames := range propNames {
		response := Response{Href: href}
		response = builder.AddPropertyStatus(response, pNames, 200)
		responses = append(responses, response)
	}

	multistatus := builder.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in UserPrincipalPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(addXMLDeclaration(responseBytes))
}

func (s *Service) buildCalendarsPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop) (foundProps map[string]CalendarProperties, notFoundProps map[string]CalendarProperties, err error) {
	// Get calendars from domain
	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return foundProps, notFoundProps, err
	}

	var foundP CalendarProperties
	var notFoundP CalendarProperties

	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Collection: &Collection{},
			}
		case raw.XMLName.Local == "displayname":
			foundP.DisplayName = "Calendars"
		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	calendarHomeSetPath := fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId)
	foundProps = map[string]CalendarProperties{
		calendarHomeSetPath: foundP,
	}
	notFoundProps = map[string]CalendarProperties{
		calendarHomeSetPath: notFoundP,
	}

	for _, calendar := range calendars {
		calendarPath := fmt.Sprintf("%s/%d", calendarHomeSetPath, calendar.CalendarId)
		foundP = CalendarProperties{}
		notFoundP = CalendarProperties{}
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
				foundP.GetETag = calendar.UpdateTime.UnixNano()
			case raw.XMLName.Local == "getctag":
				foundP.GetCTag = calendar.UpdateTime.UnixNano()
			case raw.XMLName.Local == "getlastmodified":
				foundP.GetLastModified = calendar.UpdateTime.Format(time.RFC1123)
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
			default:
				notFoundP.Raw = append(notFoundP.Raw, raw)
			}
		}
		foundProps[calendarPath] = foundP
		notFoundProps[calendarPath] = notFoundP
	}

	return foundProps, notFoundProps, nil
}

func (s *Service) buildCalendarsAllPropResponse(ctx context.Context, authAccount model.AuthAccount) (foundProps map[string]CalendarProperties, err error) {
	// Get calendars from domain
	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return foundProps, err
	}

	calendarHomeSetPath := fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId)
	foundProps = map[string]CalendarProperties{
		calendarHomeSetPath: {
			ResourceType: &ResourceType{
				Collection: &Collection{},
			},
			DisplayName: "Calendars",
		},
	}

	for _, calendar := range calendars {
		calendarPath := fmt.Sprintf("%s/%d", calendarHomeSetPath, calendar.CalendarId)
		privileges := []Privilege{
			{Name: "D:read"},
		}
		if calendar.CalendarAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			privileges = append(privileges, Privilege{Name: "D:write"})
			privileges = append(privileges, Privilege{Name: "D:write-acl"})
		}
		foundProps[calendarPath] = CalendarProperties{
			ResourceType: &ResourceType{
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
			CurrentUserPrivilegeSet: &PrivilegeSet{Privileges: privileges},
		}
	}

	return foundProps, nil
}

func (s *Service) buildCalendarsPropNameResponse(ctx context.Context, authAccount model.AuthAccount) (map[string]CalendarPropertyNames, error) {
	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{model.CalendarField_CalendarId})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return nil, err
	}

	calendarHomeSetPath := fmt.Sprintf("/caldav/principals/%s/calendars", authAccount.AuthUserId)
	propNames := map[string]CalendarPropertyNames{
		calendarHomeSetPath: {
			ResourceType: &struct{}{},
			DisplayName:  &struct{}{},
		},
	}
	for _, calendar := range calendars {
		calendarPath := fmt.Sprintf("%s/%s", calendarHomeSetPath, calendar.CalendarId)
		propNames[calendarPath] = CalendarPropertyNames{
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
		}
	}

	return propNames, nil
}

func hasAnyCalendarProperties(prop CalendarProperties) bool {
	return prop.ResourceType != nil ||
		prop.DisplayName != "" ||
		prop.GetETag != 0 ||
		prop.GetCTag != 0 ||
		prop.GetLastModified != "" ||
		prop.CalendarDescription != "" ||
		(prop.SupportedCalendarComponentSet != nil && len(prop.SupportedCalendarComponentSet.CalendarComponents) > 0) ||
		(prop.SupportedCalendarData != nil && len(prop.SupportedCalendarData.CalendarData) > 0) ||
		(prop.SupportedReportSet != nil && len(prop.SupportedReportSet.SupportedReports) > 0) ||
		len(prop.Raw) > 0
}
