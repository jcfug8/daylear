package caldav

import (
	"io"
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

func (s *Service) WellKnown(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("WellKnown called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in UserPrincipal")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.WellKnownOptions(w, r)
		return
	case "PROPFIND":
		s.WellKnownPropFind(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) WellKnownOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET,OPTIONS")
	w.Header().Set("DAV", "1, 2, calendar-access")
	w.Header().Set("CalDAV", "calendar-access")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) WellKnownPropFind(w http.ResponseWriter, r *http.Request, _ model.AuthAccount) {
	body := r.Body
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to read body in WellKnownPropFind")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	s.log.Info().Msg("Body: " + string(bodyBytes))
	// w.Header().Set("Location", s.NewResponseHref(fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId)).Href)
	w.Header().Set("Location", s.NewResponseHref("/caldav/").Href)
	w.WriteHeader(http.StatusPermanentRedirect)
}
