package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	circleMaxPageSize     int32 = 1000
	circleDefaultPageSize int32 = 100
)

// CreateCircle -
func (s *CircleService) CreateCircle(ctx context.Context, request *pb.CreateCircleRequest) (response *pb.Circle, err error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		return nil, err
	}

	// convert proto to model
	circleProto := request.GetCircle()
	circleProto.Name = ""
	mCircle, err := s.ProtoToCircle(circleProto)
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
	circleProto, err = s.CircleToProto(mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)

	return circleProto, nil
}

// DeleteCircle -
func (s *CircleService) DeleteCircle(ctx context.Context, request *pb.DeleteCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx)
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
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err := s.CircleToProto(mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)

	return circleProto, nil
}

// GetCircle -
func (s *CircleService) GetCircle(ctx context.Context, request *pb.GetCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	var mCircle model.Circle
	_, err = s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCircle, err = s.domain.GetCircle(ctx, authAccount, mCircle.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err := s.CircleToProto(mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)

	return circleProto, nil
}

// UpdateCircle -
func (s *CircleService) UpdateCircle(ctx context.Context, request *pb.UpdateCircleRequest) (*pb.Circle, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	circleProto := request.GetCircle()
	var mCircle model.Circle
	_, err = s.circleNamer.Parse(circleProto.GetName(), &mCircle)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", circleProto.GetName())
	}

	fieldMask := request.GetUpdateMask()
	updateMask, err := s.circleFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = s.ProtoToCircle(circleProto)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mCircle, err = s.domain.UpdateCircle(ctx, authAccount, mCircle, updateMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err = s.CircleToProto(mCircle)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)

	return circleProto, nil
}

// ListCircles -
func (s *CircleService) ListCircles(ctx context.Context, request *pb.ListCirclesRequest) (*pb.ListCirclesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	fieldMask := s.circleFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.circleFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: circleDefaultPageSize,
		MaxPageSize:     circleMaxPageSize,
	})
	if err != nil {
		return nil, err
	}
	request.PageSize = pageSize

	circles, err := s.domain.ListCircles(ctx, authAccount, request.GetPageSize(), pageToken.Offset, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	circleProtos, err := s.CircleListToProto(circles)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	for _, circleProto := range circleProtos {
		grpc.ProcessResponseFieldBehavior(circleProto)
	}

	response := &pb.ListCirclesResponse{
		Circles: circleProtos,
	}

	if len(circleProtos) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	// prepare response
	return response, nil
}

// ProtoToCircle converts a protobuf Circle to a model Circle
func (s *CircleService) ProtoToCircle(proto *pb.Circle) (model.Circle, error) {
	circle := model.Circle{}
	if proto.Name != "" {
		_, err := s.circleNamer.Parse(proto.Name, &circle)
		if err != nil {
			return circle, err
		}
	}
	circle.Title = proto.Title
	circle.ImageURI = proto.ImageUri
	circle.VisibilityLevel = proto.Visibility
	circle.PermissionLevel = proto.Permission
	circle.AccessState = proto.State
	return circle, nil
}

// CircleToProto converts a model Circle to a protobuf Circle
func (s *CircleService) CircleToProto(circle model.Circle) (*pb.Circle, error) {
	proto := &pb.Circle{}
	name, err := s.circleNamer.Format(circle)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Title = circle.Title
	proto.ImageUri = circle.ImageURI
	proto.Visibility = circle.VisibilityLevel
	proto.Permission = circle.PermissionLevel
	proto.State = circle.AccessState
	return proto, nil
}

// CircleListToProto converts a slice of model Circles to a slice of protobuf Circles
func (s *CircleService) CircleListToProto(circles []model.Circle) ([]*pb.Circle, error) {
	protos := make([]*pb.Circle, len(circles))
	for i, circle := range circles {
		proto, err := s.CircleToProto(circle)
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToCircle converts a slice of protobuf Circles to a slice of model Circles
func (s *CircleService) ProtosToCircle(protos []*pb.Circle) ([]model.Circle, error) {
	res := make([]model.Circle, len(protos))
	for i, proto := range protos {
		circle, err := s.ProtoToCircle(proto)
		if err != nil {
			return nil, err
		}
		res[i] = circle
	}
	return res, nil
}
