package caldav

import (
	"fmt"
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
	w.Header().Set("DAV", "1, 2, calendar-access, calendar-schedule")
	w.Header().Set("CalDAV", "calendar-access, calendar-schedule")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) WellKnownPropFind(w http.ResponseWriter, r *http.Request, authAccount model.AuthAccount) {
	w.Header().Set("Location", fmt.Sprintf("/caldav/principals/%d", authAccount.AuthUserId))
	w.WriteHeader(http.StatusPermanentRedirect)
}
