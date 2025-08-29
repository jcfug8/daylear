package caldav

import (
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
)

func (s *Service) Calendar(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Calendar called")

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Calendar")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "OPTIONS":
		s.CalendarOptions(w, r)
		return
	case "PROPFIND":
		s.CalendarPropFind(w, r, authAccount)
		return
	case "REPORT":
		s.CalendarReport(w, r, authAccount)
		return
	case "GET":
		s.CalendarGet(w, r, authAccount)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Service) CalendarOptions(w http.ResponseWriter, r *http.Request) {
	setCalDAVHeaders(w)
	w.Header().Set("Allow", "PROPFIND,OPTIONS,REPORT,GET")
	w.WriteHeader(http.StatusNoContent)
}
