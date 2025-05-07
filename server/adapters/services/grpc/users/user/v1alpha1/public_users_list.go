package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc/pagination"
	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IRIOMO:CUSTOM_CODE_SLOT_START recipeServiceListConstants

const (
	userMaxPageSize     int32 = 1000
	userDefaultPageSize int32 = 100
)

// IRIOMO:CUSTOM_CODE_SLOT_END

// Listrecipes -
func (s *UserService) ListPublicUsers(ctx context.Context, request *pb.ListPublicUsersRequest) (*pb.ListPublicUsersResponse, error) {
	fieldMask := s.userFieldMasker.GetFieldMaskFromCtx(ctx)

	readMask, err := s.userFieldMasker.GetReadMask(fieldMask)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid field mask")
	}

	pageToken, err := pagination.ParsePageToken[cmodel.User](request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if pageToken.PageSize == 0 {
		pageToken.PageSize = userDefaultPageSize
	}
	pageToken.PageSize = min(pageToken.PageSize, userMaxPageSize)

	res, err := s.domain.ListUsers(ctx, pageToken, request.GetFilter(), readMask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	users, err := convert.PublicUserListToProto(s.userNamer, res)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to prepare response")
	}

	return &pb.ListPublicUsersResponse{
		NextPageToken: pagination.EncodePageToken(pageToken.Next(res)),
		PublicUsers:   users,
	}, nil
}
