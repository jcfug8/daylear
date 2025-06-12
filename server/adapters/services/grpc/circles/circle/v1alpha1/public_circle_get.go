package v1alpha1

import (
	"context"

	// TODO: implement metadata if needed for auth, as in user
	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPublicCircle -
func (s *CircleService) GetPublicCircle(ctx context.Context, request *pb.GetPublicCircleRequest) (*pb.PublicCircle, error) {
	var mCircle model.Circle
	_, err := s.publicCircleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	fieldMask := s.publicCircleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.publicCircleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = s.domain.GetCircle(ctx, model.CircleParent{}, mCircle.Id, readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circle, err := convert.PublicCircleToProto(s.publicCircleNamer, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return circle, nil
}
