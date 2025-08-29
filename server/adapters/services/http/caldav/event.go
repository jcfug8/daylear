package caldav

import (
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
)

func (s *Service) Event(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Event called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Event")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.EventOptions(w, r)
		return
	case "GET":
		s.EventGet(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) EventOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "OPTIONS,GET")
	w.WriteHeader(http.StatusNoContent)
}
