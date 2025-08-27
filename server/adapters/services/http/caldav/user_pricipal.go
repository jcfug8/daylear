package caldav

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

// User Principal specific property structures
type UserPrincipalProp struct {
	ResourceType            *ResourceType `xml:"D:resourcetype,omitempty"`
	DisplayName             string        `xml:"D:displayname,omitempty"`
	CurrentUserPrincipal    *Href         `xml:"D:current-user-principal,omitempty"`
	CalendarHomeSet         *Href         `xml:"C:calendar-home-set,omitempty"`
	PrincipalURL            *Href         `xml:"C:principal-URL,omitempty"`
	Owner                   string        `xml:"D:owner,omitempty"`
	CurrentUserPrivilegeSet *PrivilegeSet `xml:"D:current-user-privilege-set,omitempty"`
	Raw                     []RawXMLValue `xml:",any"`
}

type UserPrincipalPropNames struct {
	ResourceType            struct{} `xml:"D:resourcetype"`
	DisplayName             struct{} `xml:"D:displayname"`
	CurrentUserPrincipal    struct{} `xml:"D:current-user-principal"`
	CalendarHomeSet         struct{} `xml:"C:calendar-home-set"`
	PrincipalURL            struct{} `xml:"C:principal-URL"`
	Owner                   struct{} `xml:"D:owner"`
	CurrentUserPrivilegeSet struct{} `xml:"D:current-user-privilege-set"`
}

func (s *Service) UserPrincipal(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("UserPrincipal called")

	userID, err := strconv.ParseInt(mux.Vars(r)["userID"], 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in UserPrincipal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in UserPrincipal")
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
		s.UserPrincipalOptions(w, r)
		return
	case "PROPFIND":
		s.UserPrincipalPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) UserPrincipalOptions(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) UserPrincipalPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	propFindRequest, err := NewPropFindRequestFromReader(r.Body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse PROPFIND request")
		http.Error(w, "Invalid XML", http.StatusBadRequest)
		return
	}

	var foundProps UserPrincipalProp
	var notFoundProps UserPrincipalProp
	var propNames *UserPrincipalPropNames

	// Build the response based on what was requested
	switch propFindRequest.GetRequestType() {
	case PropFindRequestTypeProp:
		foundProps, notFoundProps, err = s.buildUserPrincipalPropResponse(r.Context(), r.URL.Path, authAccount, propFindRequest.Prop)
	case PropFindRequestTypeAllProp:
		foundProps, err = s.buildUserPrincipalAllPropResponse(r.Context(), r.URL.Path, authAccount)
	case PropFindRequestTypePropName:
		propNames = s.buildUserPrincipalPropNameResponse()
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

	response := Response{Href: r.URL.Path}
	builder := ResponseBuilder{}

	// Add propstat for found properties
	if hasAnyUserPrincipalProperties(foundProps) {
		response = builder.AddPropertyStatus(response, foundProps, 200)
	}

	// Add propstat for not found properties
	if hasAnyUserPrincipalProperties(notFoundProps) {
		response = builder.AddPropertyStatus(response, notFoundProps, 404)
	}

	if propNames != nil {
		response = builder.AddPropertyStatus(response, propNames, 200)
	}

	multistatus := builder.BuildMultiStatusResponse([]Response{response})

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

func (s *Service) buildUserPrincipalPropResponse(ctx context.Context, href string, authAccount model.AuthAccount, prop *Prop) (foundProps UserPrincipalProp, notFoundProps UserPrincipalProp, err error) {
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		return foundProps, notFoundProps, err
	}

	// Check each requested property
	for _, raw := range prop.Raw {
		switch {
		case raw.XMLName.Local == "resourcetype":
			foundProps.ResourceType = &ResourceType{
				Collection: &Collection{},
				Principal:  &Principal{},
			}

		case raw.XMLName.Local == "displayname":
			foundProps.DisplayName = user.GetFullName()

		case raw.XMLName.Local == "current-user-principal":
			foundProps.CurrentUserPrincipal = &Href{
				Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
			}

		case raw.XMLName.Local == "principal-URL":
			foundProps.PrincipalURL = &Href{
				Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
			}

		case raw.XMLName.Local == "calendar-home-set":
			foundProps.CalendarHomeSet = &Href{
				Href: fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId),
			}

		case raw.XMLName.Local == "owner":
			foundProps.Owner = fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)

		case raw.XMLName.Local == "current-user-privilege-set":
			foundProps.CurrentUserPrivilegeSet = &PrivilegeSet{
				Privileges: []Privilege{
					{Name: "D:read"},
					{Name: "D:write"},
					{Name: "D:write-acl"},
				},
			}

		default:
			notFoundProps.Raw = append(notFoundProps.Raw, raw)
		}
	}

	return foundProps, notFoundProps, nil
}

func (s *Service) buildUserPrincipalAllPropResponse(ctx context.Context, href string, authAccount model.AuthAccount) (UserPrincipalProp, error) {
	user, err := s.domain.GetOwnUser(ctx, authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		return UserPrincipalProp{}, err
	}

	return UserPrincipalProp{
		ResourceType: &ResourceType{
			Collection: &Collection{},
			Principal:  &Principal{},
		},
		DisplayName: user.GetFullName(),
		CurrentUserPrincipal: &Href{
			Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
		},
		CalendarHomeSet: &Href{
			Href: fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId),
		},
		PrincipalURL: &Href{
			Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
		},
		Owner: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
		CurrentUserPrivilegeSet: &PrivilegeSet{
			Privileges: []Privilege{
				{Name: "D:read"},
				{Name: "D:write"},
				{Name: "D:write-acl"},
			},
		},
	}, nil
}

func (s *Service) buildUserPrincipalPropNameResponse() *UserPrincipalPropNames {
	return &UserPrincipalPropNames{
		ResourceType:            struct{}{},
		DisplayName:             struct{}{},
		CurrentUserPrincipal:    struct{}{},
		CalendarHomeSet:         struct{}{},
		PrincipalURL:            struct{}{},
		Owner:                   struct{}{},
		CurrentUserPrivilegeSet: struct{}{},
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
