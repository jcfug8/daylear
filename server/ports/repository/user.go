package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the user in the database.
type userClient interface {
	CreateUser(ctx context.Context, user model.User, fields []string) (model.User, error)
	DeleteUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.User, error)
	GetUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fields []string) (model.User, error)
	ListUsers(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]model.User, error)
	UpdateUser(ctx context.Context, authAccount model.AuthAccount, user model.User, fields []string) (model.User, error)

	FindStandardUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, error)
	FindDelegatedUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, model.UserAccess, error)

	CreateUserAccess(ctx context.Context, access model.UserAccess, fields []string) (model.UserAccess, error)
	DeleteUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) error
	GetUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId, fields []string) (model.UserAccess, error)
	ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.UserAccess, error)
	UpdateUserAccess(ctx context.Context, access model.UserAccess, fields []string) (model.UserAccess, error)

	CreateUserFavorite(ctx context.Context, authAccount model.AuthAccount, id model.UserId) error
	DeleteUserFavorite(ctx context.Context, authAccount model.AuthAccount, id model.UserId) error
}
