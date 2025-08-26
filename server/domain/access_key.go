package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Domain) AuthenticateByAccessKey(ctx context.Context, userId int64, secretKey string) (model.User, error) {
	if userId == 0 || secretKey == "" {
		return model.User{}, status.Error(codes.InvalidArgument, "access key and secret key are required")
	}

	if userId != 1 || secretKey != "123456" {
		return model.User{}, status.Error(codes.InvalidArgument, "invalid access key or secret key")
	}

	return model.User{
		Id: model.UserId{UserId: userId},
	}, nil
}
