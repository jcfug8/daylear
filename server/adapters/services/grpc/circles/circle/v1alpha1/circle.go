package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	circleMaxPageSize     int32 = 1000
	circleDefaultPageSize int32 = 100
)

var circleFieldMap = map[string][]string{
	"name":        {model.CircleField_Parent, model.CircleField_CircleId},
	"title":       {model.CircleField_Title},
	"description": {model.CircleField_Description},
	"handle":      {model.CircleField_Handle},
	"image_uri":   {model.CircleField_ImageURI},
	"visibility":  {model.CircleField_Visibility},

	"circle_access": {model.CircleField_CircleAccess},
}

// CreateCircle -
func (s *CircleService) CreateCircle(ctx context.Context, request *pb.CreateCircleRequest) (response *pb.Circle, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateCircle called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	circleProto := request.GetCircle()
	circleProto.Name = ""
	mCircle, err := s.ProtoToCircle(circleProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// create circle
	mCircle, err = s.domain.CreateCircle(ctx, authAccount, mCircle)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateCircle failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	circleProto, err = s.CircleToProto(mCircle)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)
	log.Info().Msg("gRPC CreateCircle returning successfully")
	return circleProto, nil
}

// DeleteCircle -
func (s *CircleService) DeleteCircle(ctx context.Context, request *pb.DeleteCircleRequest) (*pb.Circle, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteCircle called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mCircle model.Circle
	_, err = s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCircle, err = s.domain.DeleteCircle(ctx, authAccount, mCircle.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteCircle failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err := s.CircleToProto(mCircle)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)
	log.Info().Msg("gRPC DeleteCircle returning successfully")
	return circleProto, nil
}

// GetCircle -
func (s *CircleService) GetCircle(ctx context.Context, request *pb.GetCircleRequest) (*pb.Circle, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetCircle called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mCircle model.Circle
	nameIndex, err := s.circleNamer.Parse(request.GetName(), &mCircle)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mCircle, err = s.domain.GetCircle(ctx, authAccount, mCircle.Parent, mCircle.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetCircle failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err := s.CircleToProto(mCircle, namer.AsPatternIndex(nameIndex))
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)
	log.Info().Msg("gRPC GetCircle returning successfully")
	return circleProto, nil
}

// UpdateCircle -
func (s *CircleService) UpdateCircle(ctx context.Context, request *pb.UpdateCircleRequest) (*pb.Circle, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateCircle called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	circleProto := request.GetCircle()
	var mCircle model.Circle
	_, err = s.circleNamer.Parse(circleProto.GetName(), &mCircle)
	if err != nil {
		log.Warn().Err(err).Str("name", circleProto.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", circleProto.GetName())
	}

	fieldMask := request.GetUpdateMask()
	updateMask := s.circleFieldMasker.Convert(fieldMask.GetPaths())
	if err != nil {
		log.Warn().Err(err).Msg("invalid field mask")
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mCircle, err = s.ProtoToCircle(circleProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.Internal, err.Error())
	}

	mCircle, err = s.domain.UpdateCircle(ctx, authAccount, mCircle, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateCircle failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	circleProto, err = s.CircleToProto(mCircle)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(circleProto)
	log.Info().Msg("gRPC UpdateCircle returning successfully")
	return circleProto, nil
}

// ListCircles -
func (s *CircleService) ListCircles(ctx context.Context, request *pb.ListCirclesRequest) (*pb.ListCirclesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListCircles called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mCircleParent := model.CircleParent{}
	nameIndex, err := s.circleNamer.ParseParent(request.GetParent(), &mCircleParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: circleDefaultPageSize,
		MaxPageSize:     circleMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}
	request.PageSize = pageSize

	circles, err := s.domain.ListCircles(ctx, authAccount, mCircleParent, request.GetPageSize(), pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListCircles failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	circleProtos := make([]*pb.Circle, len(circles))
	for i, circle := range circles {
		circleProto, err := s.CircleToProto(circle, namer.AsPatternIndex(nameIndex))
		if err != nil {
			return nil, err
		}
		circleProtos[i] = circleProto
	}
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
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

	log.Info().Msg("gRPC ListCircles returning successfully")
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
	circle.Description = proto.Description
	circle.Handle = proto.Handle
	circle.ImageURI = proto.ImageUri
	circle.VisibilityLevel = proto.Visibility
	return circle, nil
}

// CircleToProto converts a model Circle to a protobuf Circle
func (s *CircleService) CircleToProto(circle model.Circle, nameIndex ...namer.FormatReflectNamerOption) (*pb.Circle, error) {
	proto := &pb.Circle{}
	name, err := s.circleNamer.Format(circle)
	if err != nil {
		return proto, err
	}
	proto.Name = name
	proto.Title = circle.Title
	proto.Description = circle.Description
	proto.Handle = circle.Handle
	proto.ImageUri = circle.ImageURI
	proto.Visibility = circle.VisibilityLevel

	if circle.CircleAccess.CircleAccessId.CircleAccessId != 0 {
		name, err := s.accessNamer.Format(circle.CircleAccess)
		if err != nil {
			return proto, err
		}
		proto.CircleAccess = &pb.Circle_CircleAccess{
			Name:            name,
			PermissionLevel: circle.CircleAccess.PermissionLevel,
			State:           circle.CircleAccess.State,
			AcceptTarget:    circle.CircleAccess.AcceptTarget,
		}
	}

	return proto, nil
}
