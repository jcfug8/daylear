package headers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
)

type userKeyType string
const UserKey userKeyType = "auth-token-user"

func NewAuthTokenMiddleware(domain domain.Domain) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetAuthToken(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := domain.ParseToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user)))
		})
	}
}

func GetAuthToken(r *http.Request) (string, error) {
	// Retrieve the auth-token cookie
	headers := r.Header["Authorization"]
	if len(headers) != 1 {
		// For any other error, return a bad request status
		return "", fileretriever.ErrInvalidArgument{Msg: "missing or invalid authorization token"}
	}

	authToken := strings.TrimPrefix(headers[0], "Bearer ")

	return authToken, nil
}
