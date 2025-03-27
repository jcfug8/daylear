package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/grpc/metadata"
	convert "github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/core/errz"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser -
func (s *UserService) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.User, error) {
	token, err := metadata.GetAuthToken(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "unable to get token")
	}

	tokenUser, err := s.domain.ParseToken(ctx, token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "unable to parse token")
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
		return nil, errz.Sanitize(err)
	}

	user, err := convert.UserToProto(s.userNamer, mUser)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	return user, nil
}
