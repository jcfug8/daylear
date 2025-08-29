package caldav

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
)

func (s *Service) Principals(w http.ResponseWriter, r *http.Request) {
	s.log.Info().Msg("Principals called")

	userID, err := strconv.ParseInt(mux.Vars(r)["userID"], 10, 64)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse userID in Principals")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authAccount, err := headers.ParseAuthData(r.Context())
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse auth data in Principals")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userID != authAccount.AuthUserId {
		s.log.Error().Int64("userID", userID).Int64("authUserID", authAccount.AuthUserId).Msg("UserID does not match authUserID in Principals")
		w.WriteHeader(http.StatusForbidden)
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
	w.Header().Set("Allow", "PROPFIND,OPTIONS")
	w.WriteHeader(http.StatusNoContent)
}
