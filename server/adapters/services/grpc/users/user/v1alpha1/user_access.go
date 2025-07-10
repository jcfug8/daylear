package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/jcfug8/daylear/server/core/logutil"
)

var (
	userAccessMaxPageSize     int32 = 1000
	userAccessDefaultPageSize int32 = 100
)

// CreateAccess -
func (s *UserService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
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

	// parse parent user name
	var mUserParent model.User
	_, err = s.userNamer.Parse(request.GetParent(), &mUserParent)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.GetParent()).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	pbAccess.Name = ""

	mUserAccess, err := ProtoToUserAccess(s.userNamer, s.accessNamer, pbAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	mUserAccess.UserId = mUserParent.Id

	// create access
	mUserAccess, err = s.domain.CreateUserAccess(ctx, authAccount, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateUserAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = UserAccessToProto(s.userNamer, s.accessNamer, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC CreateAccess returning successfully")
	return pbAccess, nil
}

// DeleteAccess -
func (s *UserService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUserAccess := model.UserAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &mUserAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	err = s.domain.DeleteUserAccess(ctx, authAccount, mUserAccess.UserAccessParent, mUserAccess.UserAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteUserAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC DeleteAccess returning successfully")
	return &emptypb.Empty{}, nil
}

// GetAccess -
func (s *UserService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUserAccess := model.UserAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &mUserAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mUserAccess, err = s.domain.GetUserAccess(ctx, authAccount, mUserAccess.UserAccessParent, mUserAccess.UserAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetUserAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbAccess, err := UserAccessToProto(s.userNamer, s.accessNamer, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC GetAccess returning successfully")
	return pbAccess, nil
}

// ListAccesses -
func (s *UserService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListAccesses called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent user name
	var mUserParent model.UserAccessParent
	_, err = s.userNamer.Parse(request.GetParent(), &mUserParent)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.GetParent()).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: userAccessDefaultPageSize,
		MaxPageSize:     userAccessMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("failed to setup pagination")
		return nil, err
	}
	request.PageSize = pageSize

	accesses, err := s.domain.ListUserAccesses(ctx, authAccount, mUserParent, request.GetPageSize(), pageToken.Offset, request.GetFilter())
	if err != nil {
		log.Error().Err(err).Msg("domain.ListUserAccesses failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessProtos := make([]*pb.Access, 0, len(accesses))
	for _, access := range accesses {
		accessProto, err := UserAccessToProto(s.userNamer, s.accessNamer, access)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		// check field behavior
		grpc.ProcessResponseFieldBehavior(accessProto)
		accessProtos = append(accessProtos, accessProto)
	}

	response := &pb.ListAccessesResponse{
		Accesses: accessProtos,
	}

	if len(accessProtos) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListAccesses returning successfully")
	return response, nil
}

// UpdateAccess -
func (s *UserService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	accessProto := request.GetAccess()
	var mUserAccess model.UserAccess
	_, err = s.accessNamer.Parse(accessProto.GetName(), &mUserAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", accessProto.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", accessProto.GetName())
	}

	// TODO: update mask

	mUserAccess, err = ProtoToUserAccess(s.userNamer, s.accessNamer, accessProto)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, status.Error(codes.Internal, err.Error())
	}

	mUserAccess, err = s.domain.UpdateUserAccess(ctx, authAccount, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateUserAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessProto, err = UserAccessToProto(s.userNamer, s.accessNamer, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(accessProto)
	log.Info().Msg("gRPC UpdateAccess returning successfully")
	return accessProto, nil
}

// AcceptAccess -
func (s *UserService) AcceptAccess(ctx context.Context, request *pb.AcceptAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC AcceptAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUserAccess := model.UserAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &mUserAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mUserAccess, err = s.domain.AcceptUserAccess(ctx, authAccount, mUserAccess.UserAccessParent, mUserAccess.UserAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.AcceptUserAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbAccess, err := UserAccessToProto(s.userNamer, s.accessNamer, mUserAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC AcceptAccess returning successfully")
	return pbAccess, nil
}

// ProtoToUserAccess converts a proto Access to a model UserAccess
func ProtoToUserAccess(userNamer, accessNamer interface{}, pbAccess *pb.Access) (model.UserAccess, error) {
	// TODO: Implement actual conversion logic
	// This should parse the name, extract requester and recipient info
	// and convert the permission level and state

	var mUserAccess model.UserAccess

	// Parse recipient from the proto (assuming it's a user name)
	if pbAccess.GetRecipient().GetName() != "" {
		var recipientUser model.User
		// Note: This needs proper implementation based on how recipient is structured
		// For now, assuming it's a user name that can be parsed
		mUserAccess.Recipient = recipientUser.Id.UserId
	}

	// Convert permission level and state
	mUserAccess.Level = pbAccess.GetLevel()
	mUserAccess.State = pbAccess.GetState()

	return mUserAccess, nil
}

// UserAccessToProto converts a model UserAccess to a proto Access
func UserAccessToProto(userNamer, accessNamer interface{}, mUserAccess model.UserAccess) (*pb.Access, error) {
	// TODO: Implement actual conversion logic
	// This should generate the proper name, requester, and recipient fields
	// and convert the permission level and state

	pbAccess := &pb.Access{
		// Name should be generated using the accessNamer
		// Name: accessNamer.Name(&mUserAccess),
		Level: mUserAccess.Level,
		State: mUserAccess.State,
		// Requester and Recipient should be converted from user IDs to names
		// Requester: userNamer.Name(&model.User{Id: mUserAccess.requester}),
		// Recipient: userNamer.Name(&model.User{Id: mUserAccess.Recipient}),
	}

	return pbAccess, nil
}
