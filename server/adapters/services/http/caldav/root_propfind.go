package caldav

import (
	"context"
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
)

// Root specific property structures
type RootProp struct {
	ResourceType            *ResourceType `xml:"D:resourcetype,omitempty"`
	PrincipalCollectionSet  *ResponseHref `xml:"C:principal-collection-set,omitempty"`
	CurrentUserPrivilegeSet *PrivilegeSet `xml:"D:current-user-privilege-set,omitempty"`
	DisplayName             string        `xml:"D:displayname,omitempty"`
	Raw                     []RawXMLValue `xml:",any"`
}

type RootPropNames struct {
	ResourceType            *struct{} `xml:"D:resourcetype,omitempty"`
	PrincipalCollectionSet  *struct{} `xml:"C:principal-collection-set,omitempty"`
	CurrentUserPrivilegeSet *struct{} `xml:"D:current-user-privilege-set,omitempty"`
	DisplayName             *struct{} `xml:"D:displayname,omitempty"`
}

func (s *Service) RootPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	s.log.Info().Msg("RootPropFind called")

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
		responses, err = s.buildRootAllPropResponse(r.Context(), authAccount)
	case PropFindRequestTypePropName:
		responses, err = s.buildRootPropNameResponse(r.Context(), authAccount)
	case PropFindRequestTypeProp:
		responses, err = s.buildRootPropResponse(r.Context(), authAccount, propFindRequest.Prop)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build root response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	multistatus := ResponseBuilder{}.BuildMultiStatusResponse(responses)

	// Marshal and send response
	responseBytes, err := xml.MarshalIndent(multistatus, "", "  ")
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal response in RootPropFind")
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

func (s *Service) buildRootPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop) ([]Response, error) {
	var foundP RootProp
	var notFoundP RootProp

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Collection: &Collection{},
			}
		case raw.XMLName.Local == "principal-collection-set":
			foundP.PrincipalCollectionSet = s.NewResponseHrefPointer("/caldav/principals")
		case raw.XMLName.Local == "current-user-privilege-set":
			privileges := []Privilege{
				{Name: "D:read"},
			}
			foundP.CurrentUserPrivilegeSet = &PrivilegeSet{Privileges: privileges}
		case raw.XMLName.Local == "displayname":
			foundP.DisplayName = "Daylear Calendar Server"
		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	rootPath := "/"
	response := Response{Href: rootPath}
	builder := ResponseBuilder{}

	// Add propstat for found properties
	if hasAnyRootPropProperties(foundP) {
		response = builder.AddPropertyStatus(response, foundP, 200)
	}
	// Add propstat for not found properties
	if hasAnyRootPropProperties(notFoundP) {
		response = builder.AddPropertyStatus(response, notFoundP, 404)
	}

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildRootAllPropResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	privileges := []Privilege{
		{Name: "D:read"},
	}

	foundP := RootProp{
		ResourceType: &ResourceType{
			Collection: &Collection{},
		},
		PrincipalCollectionSet:  s.NewResponseHrefPointer("/caldav/principals"),
		CurrentUserPrivilegeSet: &PrivilegeSet{Privileges: privileges},
		DisplayName:             "Daylear Calendar Server",
	}

	rootPath := "/"
	response := Response{Href: rootPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func (s *Service) buildRootPropNameResponse(ctx context.Context, authAccount model.AuthAccount) ([]Response, error) {
	foundP := RootPropNames{
		ResourceType:            &struct{}{},
		PrincipalCollectionSet:  &struct{}{},
		CurrentUserPrivilegeSet: &struct{}{},
		DisplayName:             &struct{}{},
	}

	rootPath := "/"
	response := Response{Href: rootPath}
	response = ResponseBuilder{}.AddPropertyStatus(response, foundP, 200)

	responses := []Response{response}

	return responses, nil
}

func hasAnyRootPropProperties(prop RootProp) bool {
	return prop.ResourceType != nil ||
		prop.PrincipalCollectionSet != nil ||
		prop.CurrentUserPrivilegeSet != nil ||
		prop.DisplayName != "" ||
		len(prop.Raw) > 0
}
