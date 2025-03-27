package metadata

import (
	"context"
	"strings"

	"github.com/jcfug8/daylear/server/core/errz"
	"google.golang.org/grpc/metadata"
)

// GetAuthToken extracts the Authorization token from the gRPC context.
func GetAuthToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errz.NewInternal("missing metadata")
	}

	authHeaders := md.Get("grpcgateway-authorization")
	if len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer "), nil
		}
		return "", errz.NewInternal("invalid authorization header format")
	}

	return "", errz.NewInternal("missing authorization token")
}
