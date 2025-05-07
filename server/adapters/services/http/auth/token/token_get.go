package token

import (
	"fmt"
	"net/http"
	"strings"
)

// GetToken - Get a token for a user based on the token key in the url
func (s *Service) GetToken(w http.ResponseWriter, r *http.Request) {
	// Retrieve the token key from the url
	tokenKey := strings.TrimSuffix(r.URL.Path[len("/auth/token/"):], "/")

	// Call the domain to retrieve the token
	token, err := s.domain.RetrieveToken(r.Context(), tokenKey)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}
