package headers

import (
	"context"
	"net/http"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	namer "github.com/jcfug8/daylear/server/core/namer"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type keyType string

const UserKey keyType = "auth-token-user"
const CircleNameKey keyType = "auth-circle-name"

const authorizationHeaderKey = "Authorization"
const actingAsCircleHeaderKey = "X-Daylear-Circle"

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

			// Get circle name if acting as a circle
			circleName := GetCircleName(r)

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)
			if circleName != "" {
				ctx = context.WithValue(ctx, CircleNameKey, circleName)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthToken(r *http.Request) (string, error) {
	// Retrieve the auth-token cookie
	headers := r.Header[authorizationHeaderKey]
	if len(headers) != 1 {
		// For any other error, return a bad request status
		return "", fileretriever.ErrInvalidArgument{Msg: "missing or invalid authorization token"}
	}

	authToken := strings.TrimPrefix(headers[0], "Bearer ")

	return authToken, nil
}

func GetCircleName(r *http.Request) string {
	headers := r.Header[actingAsCircleHeaderKey]
	if len(headers) != 1 {
		return ""
	}
	return headers[0]
}

func ParseAuthData(ctx context.Context, circleNamer namer.ReflectNamer) (model.AuthAccount, error) {
	user, ok := ctx.Value(UserKey).(model.User)
	if !ok {
		return model.AuthAccount{}, status.Error(codes.Unauthenticated, "user not found")
	}

	var circleID model.CircleId
	circleName := ctx.Value(CircleNameKey)
	if circleName != nil {
		circleNameStr, ok := circleName.(string)
		if !ok {
			return model.AuthAccount{}, status.Error(codes.InvalidArgument, "invalid circle name")
		}
		_, err := circleNamer.Parse(circleNameStr, &circleID)
		if err != nil {
			return model.AuthAccount{}, status.Errorf(codes.InvalidArgument, "invalid circle name: %v", err)
		}
	}

	return model.AuthAccount{UserId: user.Id.UserId, CircleId: circleID.CircleId}, nil
}
