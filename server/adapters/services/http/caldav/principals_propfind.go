package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
)

// Principal specific property structures
type PrincipalProp struct {
	ResourceType            *ResourceType `xml:"D:resourcetype,omitempty"`
	DisplayName             string        `xml:"D:displayname,omitempty"`
	CurrentUserPrivilegeSet *PrivilegeSet `xml:"D:current-user-privilege-set,omitempty"`
	Raw                     []RawXMLValue `xml:",any"`
}

type PrincipalPropNames struct {
	ResourceType            *struct{} `xml:"D:resourcetype,omitempty"`
	DisplayName             *struct{} `xml:"D:displayname,omitempty"`
	CurrentUserPrivilegeSet *struct{} `xml:"D:current-user-privilege-set,omitempty"`
}

func (s *Service) PrincipalsPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("PrincipalsPropFind called")

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
		responses, err = s.buildPrincipalAllPropResponse(r.Context(), authAccount)
	case PropFindRequestTypePropName:
		responses, err = s.buildPrincipalPropNameResponse(r.Context(), authAccount)
	case PropFindRequestTypeProp:
		responses, err = s.buildPrincipalPropResponse(r.Context(), authAccount, propFindRequest.Prop)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build principal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	multistatus := ResponseBuilder{}.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in PrincipalsPropFind")
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

func (s *Service) buildPrincipalPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop) ([]Response, error) {
	var foundP PrincipalProp
	var notFoundP PrincipalProp

	// Get user from domain
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in buildPrincipalPropResponse")
		return nil, err
	}

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Principal: &Principal{},
			}
		case raw.XMLName.Local == "displayname":
			foundP.DisplayName = user.GetFullName()
		case raw.XMLName.Local == "current-user-privilege-set":
			privileges := []Privilege{
				{Name: "D:read"},
			}
			foundP.CurrentUserPrivilegeSet = &PrivilegeSet{Privileges: privileges}
		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	principalPath := s.formatPrincipalPath(authAccount.AuthUserId)
	response := Response{Href: principalPath}
	builder := ResponseBuilder{}

	// Add propstat for found properties
	if hasAnyPrincipalPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}
	// Add propstat for not found properties
	if hasAnyPrincipalPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildPrincipalAllPropResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	// Get user from domain
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in buildPrincipalAllPropResponse")
		return nil, err
	}

	privileges := []Privilege{
		{Name: "D:read"},
	}

	foundP := PrincipalProp{
		ResourceType: &ResourceType{
			Principal: &Principal{},
		},
		DisplayName:             user.GetFullName(),
		CurrentUserPrivilegeSet: &PrivilegeSet{Privileges: privileges},
	}

	principalPath := s.formatPrincipalPath(authAccount.AuthUserId)
	response := Response{Href: principalPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildPrincipalPropNameResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	foundP := PrincipalPropNames{
		ResourceType:            &struct{}{},
		DisplayName:             &struct{}{},
		CurrentUserPrivilegeSet: &struct{}{},
	}

	principalPath := s.formatPrincipalPath(authAccount.AuthUserId)
	response := Response{Href: principalPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func hasAnyPrincipalPropProperties(prop PrincipalProp) bool {
	return prop.ResourceType != nil ||
		prop.DisplayName != "" ||
		prop.CurrentUserPrivilegeSet != nil ||
		len(prop.Raw) > 0
}

func (s *Service) formatPrincipalPath(userID int64) string {
	return fmt.Sprintf("/caldav/principals/%d", userID)
}
