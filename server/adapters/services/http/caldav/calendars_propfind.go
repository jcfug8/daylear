package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
)

// Calendar specific property structures
type CalendarCollectionProp struct {
	// used for both calendar and calendar collection
	ResourceType *ResourceType `xml:"D:resourcetype,omitempty"`
	// used for both calendar and calendar collection
	DisplayName string        `xml:"D:displayname,omitempty"`
	Raw         []RawXMLValue `xml:",any"`
}

type CalendarCollectionPropNames struct {
	// used for both calendar and calendar collection
	ResourceType *struct{} `xml:"D:resourcetype,omitempty"`
	// used for both calendar and calendar collection
	DisplayName *struct{} `xml:"D:displayname,omitempty"`
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

func (s *Service) CalendarsPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
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

	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var responses []Response

	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeProp:
		responses, err = s.buildCalendarsPropResponse(r.Context(), authAccount, propFindRequest.Prop, depth)
	case PropFindRequestTypeAllProp:
		responses, err = s.buildCalendarsAllPropResponse(r.Context(), authAccount, depth)
	case PropFindRequestTypePropName:
		responses, err = s.buildCalendarsPropNameResponse(r.Context(), authAccount, depth)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build calendars response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	multistatus := ResponseBuilder{}.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in UserPrincipalPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBytes = addXMLDeclaration(responseBytes)

	s.log.Info().Msgf("calendars propfind response: %s", string(responseBytes))

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(responseBytes)
}

func (s *Service) buildCalendarsPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop, depth int) ([]Response, error) {
	var foundP CalendarCollectionProp
	var notFoundP CalendarCollectionProp

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

	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	builder := ResponseBuilder{}

	if hasAnyCalendarCollectionPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}
	if hasAnyCalendarCollectionPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	depth--

	// Get calendars from domain
	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return nil, err
	}

	for _, calendar := range calendars {
		calendarResponses, err := s._buildCalendarPropResponse(ctx, authAccount, calendar, prop, depth)
		if err != nil {
			s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendar.CalendarId.CalendarId).Msg("Failed to build calendar response")
			return nil, err
		}
		responses = append(responses, calendarResponses...)
	}

	return responses, nil
}

func (s *Service) buildCalendarsAllPropResponse(ctx context.Context, authAccount model.AuthAccount, depth int) ([]Response, error) {
	foundP := CalendarCollectionProp{
		ResourceType: &ResourceType{
			Collection: &Collection{},
		},
		DisplayName: "Calendars",
	}

	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	depth--

	// Get calendars from domain
	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return nil, err
	}

	for _, calendar := range calendars {
		calendarResponses, err := s._buildCalendarAllPropResponse(ctx, authAccount, calendar, depth)
		if err != nil {
			s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendar.CalendarId.CalendarId).Msg("Failed to build calendar response")
			return nil, err
		}
		responses = append(responses, calendarResponses...)
	}

	return responses, nil
}

func (s *Service) buildCalendarsPropNameResponse(ctx context.Context, authAccount model.AuthAccount, depth int) ([]Response, error) {
	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, CalendarCollectionPropNames{
		ResourceType: &struct{}{},
		DisplayName:  &struct{}{},
	}, 200)

	responses := []Response{response}

	if depth == 0 {
		return responses, nil
	}

	depth--

	calendars, err := s.domain.ListCalendars(ctx, model.AuthAccount{AuthUserId: authAccount.AuthUserId}, model.CalendarParent{UserId: authAccount.AuthUserId}, 1000, 0, "", []string{model.CalendarField_CalendarId})
	if err != nil {
		s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Msg("Failed to list calendars from domain")
		return nil, err
	}

	for _, calendar := range calendars {
		calendarResponses, err := s.buildCalendarPropNameResponse(ctx, authAccount, calendar.CalendarId.CalendarId, depth)
		if err != nil {
			s.log.Error().Err(err).Int64("userID", authAccount.AuthUserId).Int64("calendarID", calendar.CalendarId.CalendarId).Msg("Failed to build calendar response")
			return nil, err
		}
		responses = append(responses, calendarResponses...)
	}

	return responses, nil
}

func hasAnyCalendarCollectionPropProperties(prop CalendarCollectionProp) bool {
	return prop.ResourceType != nil ||
		prop.DisplayName != "" ||
		len(prop.Raw) > 0
}

func (s *Service) formatCalendarHomeSetPath(userId int64) string {
	return path.Join(s.apiPath, fmt.Sprintf("/caldav/principals/%d/calendars", userId))
}
