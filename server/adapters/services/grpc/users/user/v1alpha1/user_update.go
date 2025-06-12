package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser -
func (s *UserService) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.User, error) {
	tokenUser, ok := ctx.Value(headers.UserKey).(cmodel.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "user not found")
	}

	mUser := model.User{}
	_, err := s.userNamer.Parse(request.GetUser().GetName(), &mUser)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetUser().GetName())
	}

	if tokenUser.Id != mUser.Id {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
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

	mUser, err = s.domain.UpdateUser(ctx, mUser, updateMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user, err := convert.UserToProto(s.userNamer, s.publicUserNamer, mUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return user, nil
}
