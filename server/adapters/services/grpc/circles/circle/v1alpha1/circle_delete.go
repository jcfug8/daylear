package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/metadata"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteCircle -
func (s *CircleService) DeleteCircle(ctx context.Context, request *pb.DeleteCircleRequest) (*pb.Circle, error) {
	// Extract the Authorization header from the gRPC context
	authToken, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing or invalid authorization token")
	}

	parent, id, err := s.circleNamer.Parse(request.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	if s.domain.AuthorizeCircleParent(ctx, authToken, parent) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "user not authorized")
	}

	mCircle, err := s.domain.DeleteCircle(ctx, parent, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err := convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return circleProto, nil
}
