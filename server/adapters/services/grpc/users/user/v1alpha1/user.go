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
)

var (
	userMaxPageSize     int32 = 1000
	userDefaultPageSize int32 = 100
)

// GetUser -
func (s *UserService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(request.GetName(), &mUser)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	fieldMask := s.userFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.userFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mUser, err = s.domain.GetUser(ctx, authAccount, mUser.Id, readMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user, err := convert.UserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return user, nil
}

// UpdateUser -
func (s *UserService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	userProto := request.GetUser()
	var mUser model.User
	_, err = s.userNamer.Parse(userProto.GetName(), &mUser)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", userProto.GetName())
	}

	fieldMask := request.GetUpdateMask()
	updateMask, err := s.userFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mUser, err = convert.ProtoToUser(s.userNamer, userProto)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mUser, err = s.domain.UpdateUser(ctx, authAccount, mUser, updateMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userProto, err = convert.UserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return userProto, nil
}

// ListUsers -
func (s *UserService) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	fieldMask := s.userFieldMasker.GetFieldMaskFromCtx(ctx)
	readMask, err := s.userFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: userDefaultPageSize,
		MaxPageSize:     userMaxPageSize,
	})
	if err != nil {
		return nil, err
	}
	request.PageSize = pageSize

	users, err := s.domain.ListUsers(ctx, authAccount, request.GetPageSize(), pageToken.Offset, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userProtos, err := convert.UserListToProto(s.userNamer, users)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	response := &pb.ListUsersResponse{
		Users: userProtos,
	}

	if len(userProtos) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	return response, nil
}
