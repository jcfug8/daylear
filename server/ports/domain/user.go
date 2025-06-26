package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	// IRIOMO:CUSTOM_CODE_SLOT_START userImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

type userDomain interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
	IdentifyUser(ctx context.Context, user model.User) (model.User, error)

	DeleteUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.User, error)
	GetUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fieldMask []string) (model.User, error)
	ListUsers(ctx context.Context, authAccount model.AuthAccount, pageSize int, offset int, filter string, fieldMask []string) ([]model.User, error)
	UpdateUser(ctx context.Context, authAccount model.AuthAccount, user model.User, updateMask []string) (model.User, error)

	CreateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (model.UserAccess, error)
	DeleteUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) error
	GetUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error)
	ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int32, filter string) ([]model.UserAccess, error)
	UpdateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (model.UserAccess, error)
	AcceptUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error)
}
