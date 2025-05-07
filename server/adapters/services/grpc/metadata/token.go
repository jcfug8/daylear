package metadata

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetAuthToken extracts the Authorization token from the gRPC context.
func GetAuthToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Internal, "missing metadata")
	}

	authHeaders := md.Get("grpcgateway-authorization")
	if len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer "), nil
		}
		return "", status.Errorf(codes.Internal, "invalid authorization header format")
	}

	return "", status.Errorf(codes.Internal, "missing authorization token")
}
