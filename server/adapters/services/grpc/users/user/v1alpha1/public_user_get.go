package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPublicUser -
func (s *UserService) GetPublicUser(ctx context.Context, request *pb.GetPublicUserRequest) (*pb.PublicUser, error) {
	tokenUser, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "user not found")
	}

	id, err := s.userNamer.Parse(request.GetName())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	if tokenUser.Id != id {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}

	fieldMask := s.userFieldMasker.GetFieldMaskFromCtx(ctx)

	readMask, err := s.userFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid field mask")
	}

	mUser, err := s.domain.GetUser(ctx, id, readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user, err := convert.PublicUserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return user, nil
}
