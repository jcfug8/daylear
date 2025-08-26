package caldav

import (
	"encoding/xml"
	"net/http"
	"strconv"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

// CalDAV XML response structures
type RootMultiStatus struct {
	XMLName  xml.Name     `xml:"D:multistatus"`
	XMLNS    string       `xml:"xmlns:D,attr"`
	Response RootResponse `xml:"D:response"`
}

type RootResponse struct {
	Href     string       `xml:"D:href"`
	Propstat RootPropstat `xml:"D:propstat"`
}

type RootPropstat struct {
	Prop   RootProp `xml:"D:prop"`
	Status Status   `xml:"D:status"`
}

type RootProp struct {
	PrincipalCollectionSet Href `xml:"C:principal-collection-set"`
}

func (s *Service) Root(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Root called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Root")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.RootOptions(w, r)
		return
	case "PROPFIND":
		s.RootPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) RootOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) RootPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	// Create the response structure
	response := RootMultiStatus{
		XMLNS: "DAV:",
		Response: RootResponse{
			Href: r.URL.Path,
			Propstat: RootPropstat{
				Prop: RootProp{
					PrincipalCollectionSet: Href{
						Href: "/caldav/principals",
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
		s.log.Error().Err(err).Msg("Failed to marshal response in RootPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setCalDAVHeaders(w)
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseBytes)))
	w.WriteHeader(http.StatusMultiStatus)
	w.Write(addXMLDeclaration(responseBytes))
}
