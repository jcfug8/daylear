package domain

import (
	"context"
	"io"

	model "github.com/jcfug8/daylear/server/core/model"
)

type userDomain interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	IdentifyUser(ctx context.Context, user model.User) (model.User, error)
	AuthenticateByAccessKey(ctx context.Context, userId int64, secretKey string) (model.User, error)

	DeleteUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.User, error)
	GetUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId, fields []string) (model.User, error)
	GetOwnUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fields []string) (model.User, error)
	ListUsers(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, pageSize int32, offset int64, filter string, fields []string) ([]model.User, error)
	UpdateUser(ctx context.Context, authAccount model.AuthAccount, user model.User, fields []string) (model.User, error)

	FavoriteUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId) error
	UnfavoriteUser(ctx context.Context, authAccount model.AuthAccount, parent model.UserParent, id model.UserId) error

	CreateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (model.UserAccess, error)
	DeleteUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) error
	GetUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId, fields []string) (model.UserAccess, error)
	ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.UserAccess, error)
	AcceptUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error)

	// UploadUserImage uploads a user image and returns the image URI
	UploadUserImage(ctx context.Context, authAccount model.AuthAccount, id model.UserId, imageReader io.Reader) (imageURI string, err error)
}
