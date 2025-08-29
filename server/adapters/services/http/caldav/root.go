package caldav

import (
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
)

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
