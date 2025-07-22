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

	mUserAccess, err := s.ProtoToUserAccess(pbAccess)
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
	pbAccess, err = s.UserAccessToProto(mUserAccess)
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

	pbAccess, err := s.UserAccessToProto(mUserAccess)
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
		accessProto, err := s.UserAccessToProto(access)
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

	pbAccess, err := s.UserAccessToProto(mUserAccess)
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
func (s *UserService) ProtoToUserAccess(pbAccess *pb.Access) (model.UserAccess, error) {
	mUserAccess := model.UserAccess{
		Level: pbAccess.GetLevel(),
		State: pbAccess.GetState(),
	}

	if pbAccess.GetName() != "" {
		_, err := s.accessNamer.Parse(pbAccess.GetName(), &mUserAccess)
		if err != nil {
			return model.UserAccess{}, err
		}
	}

	if pbAccess.GetRequester().GetName() != "" {
		_, err := s.userNamer.Parse(pbAccess.GetRequester().GetName(), &mUserAccess.Requester)
		if err != nil {
			return model.UserAccess{}, err
		}
	}

	if pbAccess.GetRecipient().GetName() != "" {
		_, err := s.userNamer.Parse(pbAccess.GetRecipient().GetName(), &mUserAccess.Recipient)
		if err != nil {
			return model.UserAccess{}, err
		}
	}

	return mUserAccess, nil
}

// UserAccessToProto converts a model UserAccess to a proto Access
func (s *UserService) UserAccessToProto(mUserAccess model.UserAccess) (*pb.Access, error) {
	pbAccess := &pb.Access{
		Level: mUserAccess.Level,
		State: mUserAccess.State,
	}

	if mUserAccess.UserId.UserId != 0 && mUserAccess.UserAccessId.UserAccessId != 0 {
		name, err := s.accessNamer.Format(mUserAccess)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format access: %v", err)
		}
		pbAccess.Name = name
	}

	if mUserAccess.Requester != 0 {
		userName, err := s.userNamer.Format(mUserAccess.Requester)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format requester: %v", err)
		}
		pbAccess.Requester = &pb.Access_User{
			Name: userName,
		}
	}

	if mUserAccess.Recipient != 0 {
		userName, err := s.userNamer.Format(mUserAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_User{
			Name:       userName,
			Username:   mUserAccess.RecipientUsername,
			GivenName:  mUserAccess.RecipientGivenName,
			FamilyName: mUserAccess.RecipientFamilyName,
		}
	}

	return pbAccess, nil
}
