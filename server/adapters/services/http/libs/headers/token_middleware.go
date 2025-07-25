package headers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	circlePb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	var err error
	circleNamer, err = namer.NewReflectNamer[*circlePb.Circle]()
	if err != nil {
		panic(err)
	}
}

var circleNamer namer.ReflectNamer

type keyType string

const UserKey keyType = "auth-token-user"

const AuthorizationHeaderKey = "Authorization"

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

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthToken(r *http.Request) (string, error) {
	// Retrieve the auth-token cookie
	headers := r.Header[AuthorizationHeaderKey]
	if len(headers) != 1 {
		// For any other error, return a bad request status
		return "", fileretriever.ErrInvalidArgument{Msg: "missing or invalid authorization token"}
	}

	authToken := strings.TrimPrefix(headers[0], "Bearer ")

	return authToken, nil
}

func ParseAuthData(ctx context.Context) (model.AuthAccount, error) {
	user, ok := ctx.Value(UserKey).(model.User)
	if !ok {
		return model.AuthAccount{}, status.Error(codes.Unauthenticated, "user not found")
	}

	return model.AuthAccount{AuthUserId: user.Id.UserId}, nil
}
