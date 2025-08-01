package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

func (s *CircleService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateAccess called")
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
	pbAccess := request.GetAccess()
	pbAccess.Name = ""
	modelAccess, err := s.ProtoToCircleAccess(pbAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid access proto")
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &modelAccess)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.Parent).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// create access
	createdAccess, err := s.domain.CreateCircleAccess(ctx, authAccount, modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateCircleAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = s.CircleAccessToProto(createdAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC CreateAccess returning successfully")
	return pbAccess, nil
}

func (s *CircleService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.Name).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// delete access
	err = s.domain.DeleteCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteCircleAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC DeleteAccess returning successfully")
	return &emptypb.Empty{}, nil
}

func (s *CircleService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.Name).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// get access
	access, err := s.domain.GetCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetCircleAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.CircleAccessToProto(access)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC GetAccess returning successfully")
	return pbAccess, nil
}

func (s *CircleService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListAccesses called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent
	var circleAccessParent model.CircleAccessParent
	_, err = s.accessNamer.ParseParent(request.Parent, &circleAccessParent)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.Parent).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.Parent)
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: accessDefaultPageSize,
		MaxPageSize:     accessMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}
	request.PageSize = pageSize

	// list accesses
	accesses, err := s.domain.ListCircleAccesses(ctx, authAccount, circleAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter())
	if err != nil {
		log.Error().Err(err).Msg("domain.ListCircleAccesses failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	pbAccesses := make([]*pb.Access, len(accesses))
	for i, access := range accesses {
		pbAccess, err := s.CircleAccessToProto(access)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		pbAccesses[i] = pbAccess
	}

	// check field behavior
	for _, pbAccess := range pbAccesses {
		grpc.ProcessResponseFieldBehavior(pbAccess)
	}

	response := &pb.ListAccessesResponse{
		Accesses: pbAccesses,
	}

	if len(pbAccesses) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListAccesses returning successfully")
	return response, nil
}

func (s *CircleService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	modelAccess, err := s.ProtoToCircleAccess(pbAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// handle update mask
	updateMask := s.accessFieldMasker.Convert(request.GetUpdateMask().GetPaths())

	// update access
	updatedAccess, err := s.domain.UpdateCircleAccess(ctx, authAccount, modelAccess, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateCircleAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = s.CircleAccessToProto(updatedAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC UpdateAccess returning successfully")
	return pbAccess, nil
}

func (s *CircleService) AcceptAccess(ctx context.Context, request *pb.AcceptAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC AcceptAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.Name).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// accept access
	acceptedAccess, err := s.domain.AcceptCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.AcceptCircleAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.CircleAccessToProto(acceptedAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC AcceptAccess returning successfully")
	return pbAccess, nil
}

// Helper conversion functions (to be implemented)
func (s *CircleService) ProtoToCircleAccess(proto *pb.Access) (model.CircleAccess, error) {
	circleAccess := model.CircleAccess{
		PermissionLevel: proto.GetLevel(),
		State:           proto.GetState(),
	}
	if proto.GetName() != "" {
		_, err := s.accessNamer.Parse(proto.GetName(), &circleAccess)
		if err != nil {
			return circleAccess, err
		}
	}

	if proto.GetRecipient() != nil {
		_, err := s.userNamer.Parse(proto.GetRecipient().GetName(), &circleAccess.Recipient)
		if err != nil {
			return circleAccess, err
		}
	}

	switch pbrequester := proto.GetRequester().GetName().(type) {
	case *pb.Access_Requester_User:
		circleAccess.Requester = model.CircleRequester{}
		_, err := s.userNamer.Parse(pbrequester.User, &circleAccess.Requester)
		if err != nil {
			return circleAccess, err
		}
	case *pb.Access_Requester_Circle:
		circleAccess.Requester = model.CircleRequester{}
		_, err := s.circleNamer.Parse(proto.GetRequester().GetCircle(), &circleAccess.Requester)
		if err != nil {
			return circleAccess, err
		}
	}

	return circleAccess, nil
}

func (s *CircleService) CircleAccessToProto(circleAccess model.CircleAccess) (*pb.Access, error) {
	proto := &pb.Access{
		Level: circleAccess.PermissionLevel,
		State: circleAccess.State,
	}

	if circleAccess.CircleId.CircleId != 0 {
		name, err := s.accessNamer.Format(circleAccess)
		if err != nil {
			return nil, err
		}
		proto.Name = name
	}

	if circleAccess.Recipient.UserId != 0 {
		userName, err := s.userNamer.Format(circleAccess.Recipient)
		if err != nil {
			return nil, err
		}
		proto.Recipient = &pb.Access_User{
			Name:       userName,
			Username:   circleAccess.RecipientUsername,
			GivenName:  circleAccess.RecipientGivenName,
			FamilyName: circleAccess.RecipientFamilyName,
		}
	}

	if circleAccess.Requester.CircleId != 0 {
		name, err := s.circleNamer.Format(circleAccess.Requester)
		if err != nil {
			return nil, err
		}
		proto.Requester = &pb.Access_Requester{Name: &pb.Access_Requester_Circle{Circle: name}}
	} else if circleAccess.Requester.UserId != 0 {
		name, err := s.userNamer.Format(circleAccess.Requester)
		if err != nil {
			return nil, err
		}
		proto.Requester = &pb.Access_Requester{Name: &pb.Access_Requester_User{User: name}}
	}

	return proto, nil
}
