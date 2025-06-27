package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the user in the database.
type userClient interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	DeleteUser(ctx context.Context, id model.UserId) (model.User, error)
	GetUser(ctx context.Context, id model.UserId, fieldMask []string) (model.User, error)
	ListUsers(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]model.User, error)
	UpdateUser(ctx context.Context, user model.User, updateMask []string) (model.User, error)

	CreateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error)
	DeleteUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) error
	GetUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error)
	ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.UserAccess, error)
	UpdateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error)
}
