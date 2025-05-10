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

// ListCircles -
func (s *CircleService) ListCircles(ctx context.Context, request *pb.ListCirclesRequest) (*pb.ListCirclesResponse, error) {
	parent, err := s.circleNamer.ParseParent(request.GetParent())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	tokenUser, ok := ctx.Value(headers.UserKey).(model.User)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	if s.domain.AuthorizeCircleParent(ctx, tokenUser, parent) != nil {
		return nil, status.Error(codes.PermissionDenied, "user not authorized")
	}

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.circleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	// TODO: handle pagination and page_token
	circles, err := s.domain.ListCircles(ctx, nil, parent, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProtos, err := convert.CircleListToProto(s.circleNamer, circles)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return &pb.ListCirclesResponse{
		Circles: circleProtos,
		// NextPageToken: ...
	}, nil
}
