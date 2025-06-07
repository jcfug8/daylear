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

// UpdateCircle -
func (s *CircleService) UpdateCircle(ctx context.Context, request *pb.UpdateCircleRequest) (*pb.Circle, error) {
	tokenUser, ok := ctx.Value(headers.UserKey).(model.User)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	circleProto := request.GetCircle()
	var mCircle model.Circle
	_, err := s.circleNamer.Parse(circleProto.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", circleProto.GetName())
	}

	mCircle.Parent.UserId = tokenUser.Id.UserId

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	updateMask, err := s.circleFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = convert.ProtoToCircle(s.circleNamer, circleProto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	mCircle, err = s.domain.UpdateCircle(ctx, mCircle, updateMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err = convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return circleProto, nil
}
