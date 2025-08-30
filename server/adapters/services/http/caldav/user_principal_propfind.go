package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
)

// User Principal specific property structures
type UserPrincipalProp struct {
	ResourceType            *ResourceType `xml:"D:resourcetype,omitempty"`
	DisplayName             string        `xml:"D:displayname,omitempty"`
	CurrentUserPrincipal    *ResponseHref `xml:"D:current-user-principal,omitempty"`
	CalendarHomeSet         *ResponseHref `xml:"C:calendar-home-set,omitempty"`
	PrincipalURL            *ResponseHref `xml:"D:principal-URL,omitempty"`
	Owner                   string        `xml:"D:owner,omitempty"`
	CurrentUserPrivilegeSet *PrivilegeSet `xml:"D:current-user-privilege-set,omitempty"`
	Raw                     []RawXMLValue `xml:",any"`
}

type UserPrincipalPropNames struct {
	ResourceType            struct{} `xml:"D:resourcetype"`
	DisplayName             struct{} `xml:"D:displayname"`
	CurrentUserPrincipal    struct{} `xml:"D:current-user-principal"`
	CalendarHomeSet         struct{} `xml:"C:calendar-home-set"`
	PrincipalURL            struct{} `xml:"D:principal-URL"`
	Owner                   struct{} `xml:"D:owner"`
	CurrentUserPrivilegeSet struct{} `xml:"D:current-user-privilege-set"`
}

func (s *Service) UserPrincipalPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var foundProps map[string]UserPrincipalProp
	var notFoundProps map[string]UserPrincipalProp
	var propNames map[string]UserPrincipalPropNames

	// Build the response based on what was requested
	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeProp:
		foundProps, notFoundProps, err = s.buildUserPrincipalPropResponse(r.Context(), authAccount, propFindRequest.Prop)
	case PropFindRequestTypeAllProp:
		foundProps, err = s.buildUserPrincipalAllPropResponse(r.Context(), authAccount)
	case PropFindRequestTypePropName:
		propNames = s.buildUserPrincipalPropNameResponse(authAccount)
	default:
		s.log.Error().Msg("Invalid PROPFIND request type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to build user principal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responses := []Response{}
	builder := ResponseBuilder{}

	for href, foundP := range foundProps {
		response := Response{Href: href}
		// Add propstat for found properties
		if hasAnyUserPrincipalProperties(foundP) {
			response = builder.AddPropertyStatus(response, foundP, 200)
		}
		// Add propstat for not found properties
		if notFoundP, ok := notFoundProps[href]; ok && hasAnyUserPrincipalProperties(notFoundP) {
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

	responseBytes = addXMLDeclaration(responseBytes)

	s.log.Info().Msg("UserPrincipalPropFind response: " + string(responseBytes))

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(responseBytes)
}

func (s *Service) buildUserPrincipalPropResponse(ctx context.Context, authAccount model.AuthAccount, prop *Prop) (foundProps map[string]UserPrincipalProp, notFoundProps map[string]UserPrincipalProp, err error) {
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		return foundProps, notFoundProps, err
	}

	var foundP UserPrincipalProp
	var notFoundP UserPrincipalProp

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundP.ResourceType = &ResourceType{
				Collection: &Collection{},
				Principal:  &Principal{},
			}

		case raw.XMLName.Local == "displayname":
			foundP.DisplayName = user.GetFullName()

		case raw.XMLName.Local == "current-user-principal":
			foundP.CurrentUserPrincipal = s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId))

		case raw.XMLName.Local == "principal-URL":
			foundP.PrincipalURL = s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId))

		case raw.XMLName.Local == "calendar-home-set":
			foundP.CalendarHomeSet = s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId))

		case raw.XMLName.Local == "owner":
			foundP.Owner = fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)

		case raw.XMLName.Local == "current-user-privilege-set":
			foundP.CurrentUserPrivilegeSet = &PrivilegeSet{
				Privileges: []Privilege{
					{Name: "D:read"},
					{Name: "D:write"},
					{Name: "D:write-acl"},
				},
			}

		default:
			notFoundP.Raw = append(notFoundP.Raw, raw)
		}
	}

	userPrincipalPath := fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)

	return map[string]UserPrincipalProp{
			userPrincipalPath: foundP,
		}, map[string]UserPrincipalProp{
			userPrincipalPath: notFoundP,
		}, nil
}

func (s *Service) buildUserPrincipalAllPropResponse(ctx context.Context, authAccount model.AuthAccount) (map[string]UserPrincipalProp, error) {
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		return nil, err
	}

	foundP := UserPrincipalProp{
		ResourceType: &ResourceType{
			Collection: &Collection{},
			Principal:  &Principal{},
		},
		DisplayName:          user.GetFullName(),
		CurrentUserPrincipal: s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)),
		CalendarHomeSet:      s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId)),
		PrincipalURL:         s.NewResponseHrefPointer(fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)),
		Owner:                fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
		CurrentUserPrivilegeSet: &PrivilegeSet{
			Privileges: []Privilege{
				{Name: "D:read"},
				{Name: "D:write"},
				{Name: "D:write-acl"},
			},
		},
	}

	userPrincipalPath := fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)

	return map[string]UserPrincipalProp{
		userPrincipalPath: foundP,
	}, nil
}

func (s *Service) buildUserPrincipalPropNameResponse(authAccount model.AuthAccount) map[string]UserPrincipalPropNames {
	userPrincipalPath := fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)

	return map[string]UserPrincipalPropNames{
		userPrincipalPath: {
			ResourceType:            struct{}{},
			DisplayName:             struct{}{},
			CurrentUserPrincipal:    struct{}{},
			CalendarHomeSet:         struct{}{},
			PrincipalURL:            struct{}{},
			Owner:                   struct{}{},
			CurrentUserPrivilegeSet: struct{}{},
		},
	}
}

func hasAnyUserPrincipalProperties(prop UserPrincipalProp) bool {
	return (prop.ResourceType != nil && prop.ResourceType.Collection != nil) ||
		(prop.ResourceType != nil && prop.ResourceType.Principal != nil) ||
		prop.DisplayName != "" ||
		(prop.CurrentUserPrincipal != nil && prop.CurrentUserPrincipal.Href != "") ||
		(prop.CalendarHomeSet != nil && prop.CalendarHomeSet.Href != "") ||
		(prop.PrincipalURL != nil && prop.PrincipalURL.Href != "") ||
		prop.Owner != "" ||
		(prop.CurrentUserPrivilegeSet != nil && len(prop.CurrentUserPrivilegeSet.Privileges) > 0) ||
		len(prop.Raw) > 0
}
