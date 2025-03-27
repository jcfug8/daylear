package headers

import (
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/core/errz"
)

func GetAuthToken(r *http.Request) (string, error) {
	// Retrieve the auth-token cookie
	headers := r.Header["Authorization"]
	if len(headers) != 1 {
		// For any other error, return a bad request status
		return "", errz.NewInvalidArgument("missing or invalid authorization token")
	}

	authToken := strings.TrimPrefix(headers[0], "Bearer ")

	return authToken, nil
}
