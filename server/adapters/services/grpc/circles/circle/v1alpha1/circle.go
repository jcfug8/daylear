package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/pagination"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	circleMaxPageSize     int32 = 1000
	circleDefaultPageSize int32 = 100
)

// CreateCircle -
func (s *CircleService) CreateCircle(ctx context.Context, request *pb.CreateCircleRequest) (response *pb.Circle, err error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// convert proto to model
	circleProto := request.GetCircle()
	circleProto.Name = ""
	mCircle, err := convert.ProtoToCircle(nil, circleProto)
	if err != nil {
		s.log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// create circle
	mCircle, err = s.domain.CreateCircle(ctx, authAccount, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	circleProto, err = convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	fieldbehavior.ClearFields(circleProto, annotations.FieldBehavior_INPUT_ONLY)

	return circleProto, nil
}

// DeleteCircle -
func (s *CircleService) DeleteCircle(ctx context.Context, request *pb.DeleteCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	var mCircle model.Circle
	_, err = s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCircle, err = s.domain.DeleteCircle(ctx, authAccount, mCircle.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err := convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return circleProto, nil
}

// GetCircle -
func (s *CircleService) GetCircle(ctx context.Context, request *pb.GetCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	var mCircle model.Circle
	_, err = s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.circleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = s.domain.GetCircle(ctx, authAccount, mCircle.Id, readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err := convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return circleProto, nil
}

// UpdateCircle -
func (s *CircleService) UpdateCircle(ctx context.Context, request *pb.UpdateCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	circleProto := request.GetCircle()
	var mCircle model.Circle
	_, err = s.circleNamer.Parse(circleProto.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", circleProto.GetName())
	}

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	updateMask, err := s.circleFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = convert.ProtoToCircle(s.circleNamer, circleProto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	mCircle, err = s.domain.UpdateCircle(ctx, authAccount, mCircle, updateMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	circleProto, err = convert.CircleToProto(s.circleNamer, mCircle)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return circleProto, nil
}

// ListCircles -
func (s *CircleService) ListCircles(ctx context.Context, request *pb.ListCirclesRequest) (*pb.ListCirclesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.circleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if request.GetPageSize() == 0 {
		request.PageSize = circleDefaultPageSize
	}
	request.PageSize = min(request.PageSize, circleMaxPageSize)

	circles, err := s.domain.ListCircles(ctx, authAccount, request.GetPageSize(), pageToken.Offset, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// convert models to protos
	circleProtos, err := convert.CircleListToProto(s.circleNamer, circles)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// prepare response
	return &pb.ListCirclesResponse{
		NextPageToken: pageToken.Next(request).String(),
		Circles:       circleProtos,
	}, nil
}
