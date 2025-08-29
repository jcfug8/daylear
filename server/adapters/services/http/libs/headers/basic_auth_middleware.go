package headers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/ports/domain"
)

func NewBasicAuthMiddleware(domain domain.Domain) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="Access Key"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// if the username has an @, use the part before the @
			username = strings.Split(username, "@")[0]

			user, err := domain.AuthenticateByAccessKey(r.Context(), username, password)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Access Key"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
