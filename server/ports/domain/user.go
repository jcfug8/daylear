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

	GetUser(ctx context.Context, id model.UserId, fieldMask []string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User, updateMask []string) (model.User, error)
	ListUsers(ctx context.Context, page *model.PageToken[model.User], filter string, fieldMask []string) ([]model.User, error)
}
