package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
)

// Calendar home set specific property structures
type CalendarHomeSetProp struct {
	ResourceType    *ResourceType `xml:"D:resourcetype,omitempty"`
	CalendarHomeSet *ResponseHref `xml:"C:calendar-home-set,omitempty"`
	Raw             []RawXMLValue `xml:",any"`
}

type CalendarHomeSetPropNames struct {
	ResourceType    *struct{} `xml:"D:resourcetype,omitempty"`
	CalendarHomeSet *struct{} `xml:"C:calendar-home-set,omitempty"`
}

func (s *Service) CalendarHomeSetPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("CalendarHomeSetPropFind called")

	// Parse the PROPFIND request
	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var responses []Response

	// Build response based on prop type
	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeAllProp:
		responses, err = s.buildCalendarHomeSetAllPropResponse(r.Context(), authAccount)
	case PropFindRequestTypePropName:
		responses, err = s.buildCalendarHomeSetPropNameResponse(r.Context(), authAccount)
	case PropFindRequestTypeProp:
		responses, err = s.buildCalendarHomeSetPropResponse(r.Context(), authAccount, propFindRequest.Prop)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build calendar home set response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	multistatus := ResponseBuilder{}.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in CalendarHomeSetPropFind")
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

func (s *Service) buildCalendarHomeSetPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop) ([]Response, error) {
	var foundP CalendarHomeSetProp
	var notFoundP CalendarHomeSetProp

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Collection: &Collection{},
			}
		case raw.XMLName.Local == "calendar-home-set":
			foundP.CalendarHomeSet = s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId))
		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	builder := ResponseBuilder{}

	// Add propstat for found properties
	if hasAnyCalendarHomeSetPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}
	// Add propstat for not found properties
	if hasAnyCalendarHomeSetPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildCalendarHomeSetAllPropResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	foundP := CalendarHomeSetProp{
		ResourceType: &ResourceType{
			Collection: &Collection{},
		},
		CalendarHomeSet: s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId)),
	}

	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildCalendarHomeSetPropNameResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	foundP := CalendarHomeSetPropNames{
		ResourceType:    &struct{}{},
		CalendarHomeSet: &struct{}{},
	}

	calendarHomeSetPath := s.formatCalendarHomeSetPath(authAccount.AuthUserId)
	response := Response{Href: calendarHomeSetPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func hasAnyCalendarHomeSetPropProperties(prop CalendarHomeSetProp) bool {
	return prop.ResourceType != nil ||
		prop.CalendarHomeSet != nil ||
		len(prop.Raw) > 0
}

func (s *Service) formatCalendarHomeSetPath(userID int64) string {
	return fmt.Sprintf("/caldav/principals/%d/calendar-home-set", userID)
}
