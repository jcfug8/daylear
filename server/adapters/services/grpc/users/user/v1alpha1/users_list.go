package v1alpha1

import (
	"context"

	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser -
func (s *UserService) ListUsers(ctx context.Context, request *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")	
}
