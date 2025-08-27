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
type CalendarHomeSetMultiStatus struct {
	XMLName  xml.Name                `xml:"D:multistatus"`
	XMLNSD   string                  `xml:"xmlns:D,attr"`
	XMLNSC   string                  `xml:"xmlns:C,attr"`
	Response CalendarHomeSetResponse `xml:"D:response"`
}

type CalendarHomeSetResponse struct {
	Href     string                  `xml:"D:href"`
	Propstat CalendarHomeSetPropstat `xml:"D:propstat"`
}

type CalendarHomeSetPropstat struct {
	Prop   CalendarHomeSetProp `xml:"D:prop"`
	Status Status              `xml:"D:status"`
}

type CalendarHomeSetProp struct {
	ResourceType    ResourceType `xml:"D:resourcetype"`
	CalendarHomeSet ResponseHref `xml:"C:calendar-home-set"`
}

func (s *Service) CalendarHomeSet(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("CalendarHomeSet called")

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
		s.CalendarHomeSetOptions(w, r)
		return
	case "PROPFIND":
		s.CalendarHomeSetPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) CalendarHomeSetOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) CalendarHomeSetPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	response := CalendarHomeSetMultiStatus{
		XMLNSD: "DAV:",
		XMLNSC: "urn:ietf:params:xml:ns:caldav",
		Response: CalendarHomeSetResponse{
			Href: r.URL.Path,
			Propstat: CalendarHomeSetPropstat{
				Prop: CalendarHomeSetProp{
					ResourceType: ResourceType{
						Collection: &Collection{},
					},
					CalendarHomeSet: ResponseHref{
						Href: fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId),
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

func (s *Service) CalendarHomeSetPath(ctx context.Context) (string, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/caldav/principals/%d/calendars", authAccount.AuthUserId), nil
}
