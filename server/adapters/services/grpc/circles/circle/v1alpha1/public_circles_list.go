package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/grpc/pagination"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	circleMaxPageSize     int32 = 1000
	circleDefaultPageSize int32 = 100
)

// ListPublicCircles -
func (s *CircleService) ListPublicCircles(ctx context.Context, request *pb.ListPublicCirclesRequest) (*pb.ListPublicCirclesResponse, error) {
	fieldMask := s.publicCircleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.publicCircleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken[cmodel.Circle](request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if pageToken.PageSize == 0 {
		pageToken.PageSize = circleDefaultPageSize
	}
	pageToken.PageSize = min(pageToken.PageSize, circleMaxPageSize)

	res, err := s.domain.ListCircles(ctx, pageToken, cmodel.CircleParent{}, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circles := make([]*pb.PublicCircle, len(res))
	for i, c := range res {
		proto, err := convert.PublicCircleToProto(s.publicCircleNamer, c)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to prepare response")
		}
		circles[i] = proto
	}

	return &pb.ListPublicCirclesResponse{
		NextPageToken: pagination.EncodePageToken(pageToken.Next(res)),
		PublicCircles: circles,
	}, nil
}
