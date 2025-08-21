package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jcfug8/daylear/server/core/logutil"
)

var (
	userMaxPageSize     int32 = 1000
	userDefaultPageSize int32 = 100
)

var userFieldMap = map[string][]string{
	"name":        {model.UserField_Parent, model.UserField_Id},
	"username":    {model.UserField_Username},
	"given_name":  {model.UserField_GivenName},
	"family_name": {model.UserField_FamilyName},
	"image_uri":   {model.UserField_ImageUri},
	"bio":         {model.UserField_Bio},

	"access": {model.UserField_AccessName, model.UserField_AccessPermissionLevel, model.UserField_AccessState},
}

// GetUser -
func (s *UserService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetUser called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(request.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid user name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mUser, err = s.domain.GetUser(ctx, authAccount, mUser.Parent, mUser.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, err := convert.UserToProto(s.userNamer, s.accessNamer, mUser)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC GetUser returning successfully")
	return user, nil
}

// UpdateUser -
func (s *UserService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateUser called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	userProto := request.GetUser()
	var mUser model.User
	_, err = s.userNamer.Parse(userProto.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", userProto.GetName()).Msg("invalid user name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", userProto.GetName())
	}

	fieldMask := request.GetUpdateMask()
	updateMask := s.userFieldMasker.Convert(fieldMask.GetPaths())

	mUser, err = convert.ProtoToUser(s.userNamer, s.accessNamer, userProto)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.Internal, err.Error())
	}

	mUser, err = s.domain.UpdateUser(ctx, authAccount, mUser, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	userProto, err = convert.UserToProto(s.userNamer, s.accessNamer, mUser)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC UpdateUser returning successfully")
	return userProto, nil
}

// ListUsers -
func (s *UserService) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListUsers called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUserParent := model.UserParent{}
	_, err = s.userNamer.ParseParent(request.GetParent(), &mUserParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: userDefaultPageSize,
		MaxPageSize:     userMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}
	request.PageSize = pageSize

	users, err := s.domain.ListUsers(ctx, authAccount, mUserParent, request.GetPageSize(), pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListUsers failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	userProtos := make([]*pb.User, len(users))
	for i, user := range users {
		userProto, err := convert.UserToProto(s.userNamer, s.accessNamer, user)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		userProtos[i] = userProto
	}

	response := &pb.ListUsersResponse{
		Users: userProtos,
	}

	if len(userProtos) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListUsers returning successfully")
	return response, nil
}

// FavoriteUser -
func (s *UserService) FavoriteUser(ctx context.Context, request *pb.FavoriteUserRequest) (*pb.FavoriteUserResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC FavoriteUser called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(request.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid user name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	err = s.domain.FavoriteUser(ctx, authAccount, mUser.Parent, mUser.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.FavoriteUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC FavoriteUser returning successfully")
	return &pb.FavoriteUserResponse{}, nil
}

// UnfavoriteUser -
func (s *UserService) UnfavoriteUser(ctx context.Context, request *pb.UnfavoriteUserRequest) (*pb.UnfavoriteUserResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UnfavoriteUser called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(request.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid user name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	err = s.domain.UnfavoriteUser(ctx, authAccount, mUser.Parent, mUser.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.UnfavoriteUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC UnfavoriteUser returning successfully")
	return &pb.UnfavoriteUserResponse{}, nil
}
