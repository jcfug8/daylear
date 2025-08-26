package headers

import (
	"context"
	"net/http"
	"strconv"

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

			userId, err := strconv.ParseInt(username, 10, 64)
			if err != nil {
				w.Header().Set("WWW-Authenticate", `Basic realm="Access Key"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := domain.AuthenticateByAccessKey(r.Context(), userId, password)
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
