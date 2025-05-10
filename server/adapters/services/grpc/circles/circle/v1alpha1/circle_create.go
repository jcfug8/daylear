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

// CreateCircle -
func (s *CircleService) CreateCircle(ctx context.Context, request *pb.CreateCircleRequest) (response *pb.Circle, err error) {
	tokenUser, ok := ctx.Value(headers.UserKey).(model.User)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	circleProto := request.GetCircle()

	err = s.fieldBehaviorValidator.Validate(circleProto)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	circleProto.Name = ""

	mCircle, err := convert.ProtoToCircle(s.circleNamer, circleProto)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	mCircle.Parent, err = s.circleNamer.ParseParent(request.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	if s.domain.AuthorizeCircleParent(ctx, tokenUser, mCircle.Parent) != nil {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	mCircle, err = s.domain.CreateCircle(ctx, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err = convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return circleProto, nil
}
