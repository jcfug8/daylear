package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShareCircle -
func (s *CircleService) ShareCircle(ctx context.Context, request *pb.ShareCircleRequest) (*pb.ShareCircleResponse, error) {
	var mCircle model.Circle
	_, err := s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	tokenUser, ok := ctx.Value(headers.UserKey).(model.User)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	mCircle.Parent.UserId = tokenUser.Id.UserId

	// Recipients are not defined in proto, so just pass empty for now
	err = s.domain.ShareCircle(ctx, mCircle.Parent, nil, mCircle.Id, 0)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	mCircle, err = s.domain.GetCircle(ctx, mCircle.Parent, mCircle.Id, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err := convert.CircleToProto(s.circleNamer, s.publicCircleNamer, mCircle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return &pb.ShareCircleResponse{Circle: circleProto}, nil
}
