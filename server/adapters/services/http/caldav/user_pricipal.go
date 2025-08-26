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

// CalDAV XML response structures
type UserPrincipalMultiStatus struct {
	XMLName  xml.Name              `xml:"D:multistatus"`
	XMLNSD   string                `xml:"xmlns:D,attr"`
	XMLNSC   string                `xml:"xmlns:C,attr"`
	Response UserPrincipalResponse `xml:"D:response"`
}

type UserPrincipalResponse struct {
	Href     string                `xml:"D:href"`
	Propstat UserPrincipalPropstat `xml:"D:propstat"`
}

type UserPrincipalPropstat struct {
	Prop   UserPrincipalProp `xml:"D:prop"`
	Status Status            `xml:"D:status"`
}

type UserPrincipalProp struct {
	ResourceType         ResourceType `xml:"D:resourcetype"`
	CurrentUserPrincipal Href         `xml:"D:current-user-principal"`
	DisplayName          string       `xml:"D:displayname"`
	CalendarHomeSet      Href         `xml:"C:calendar-home-set"`
	PrincipalURL         Href         `xml:"C:principal-URL"`
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
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) UserPrincipalPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	user, err := s.domain.GetOwnUser(r.Context(), authAccount, model.UserId{UserId: authAccount.AuthUserId}, []string{})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get user in UserPrincipalPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the response structure
	response := UserPrincipalMultiStatus{
		XMLNSD: "DAV:",
		XMLNSC: "urn:ietf:params:xml:ns:caldav",
		Response: UserPrincipalResponse{
			Href: r.URL.Path,
			Propstat: UserPrincipalPropstat{
				Prop: UserPrincipalProp{
					ResourceType: ResourceType{
						Collection: &Collection{},
						Principal:  &Principal{},
					},
					DisplayName: user.GetFullName(),
					CurrentUserPrincipal: Href{
						Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
					},
					CalendarHomeSet: Href{
						Href: fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId),
					},
					PrincipalURL: Href{
						Href: fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId),
					},
				},
				Status: Status{
					Status: "HTTP/1.1 200 OK",
				},
			},
		},
	}

	responseBytes, err := xml.MarshalIndent(response, "", "  ")
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

func (s *Service) CurrentUserPrincipal(ctx context.Context) (string, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId), nil
}
