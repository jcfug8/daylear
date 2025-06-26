package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser -
func (s *UserService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.userNamer)
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
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user, err := convert.UserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return user, nil
}

// UpdateUser -
func (s *UserService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.userNamer)
	if err != nil {
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userNamer.Parse(request.GetUser().GetName(), &mUser)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetUser().GetName())
	}

	fieldMask := s.userFieldMasker.GetFieldMaskFromCtx(ctx)

	updateMask, err := s.userFieldMasker.GetWriteMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mUser, err = convert.ProtoToUser(s.userNamer, request.GetUser())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	mUser, err = s.domain.UpdateUser(ctx, authAccount, mUser, updateMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user, err := convert.UserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return user, nil
}

// GetUser -
func (s *UserService) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
