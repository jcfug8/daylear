package token

import (
	"fmt"
	"net/http"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
)

// GetToken - Get a token for a user based on the token key in the url
func (s *Service) CheckToken(w http.ResponseWriter, r *http.Request) {
	tokenUser, ok := r.Context().Value(headers.UserKey).(model.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"userId": %d}`, tokenUser.Id.UserId)))
}
