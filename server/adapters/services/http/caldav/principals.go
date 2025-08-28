package caldav

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

// CalDAV XML response structures
type PrincipalsMultiStatus struct {
	XMLName  xml.Name             `xml:"D:multistatus"`
	XMLNSD   string               `xml:"xmlns:D,attr"`
	Response []PrincipalsResponse `xml:"D:response"`
}

type PrincipalsResponse struct {
	Href     string             `xml:"D:href"`
	Propstat PrincipalsPropstat `xml:"D:propstat"`
}

type PrincipalsPropstat struct {
	Prop   PrincipalsProp `xml:"D:prop"`
	Status Status         `xml:"D:status"`
}

type PrincipalsProp struct {
	ResourceType ResourceType `xml:"D:resourcetype"`
	DisplayName  string       `xml:"D:displayname"`
}

func (s *Service) Principals(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Principals called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Root")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.PrincipalsOptions(w, r)
		return
	case "PROPFIND":
		s.PrincipalsPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) PrincipalsOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS,HEAD")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) PrincipalsPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	user, err := s.domain.GetOwnUser(r.Context(), authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the response structure
	response := PrincipalsMultiStatus{
		XMLNSD: "DAV:",
		Response: []PrincipalsResponse{
			{
				Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
				Propstat: PrincipalsPropstat{
					Prop: PrincipalsProp{
						ResourceType: ResourceType{
							Principal: &Principal{},
						},
						DisplayName: user.GetFullName(),
					},
					Status: Status{
						Status: "HTTP/1.1 200 OK",
					},
				},
			},
		},
	}

	responseBytes, err := xml.MarshalIndent(response, "", "  ")
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
