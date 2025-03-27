package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	// IRIOMO:CUSTOM_CODE_SLOT_START domainUserCreateImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// CreateUser creates a new user.
func (d *Domain) CreateUser(ctx context.Context, user model.User) (model.User, error) {

	user.Id.UserId = 0

	// IRIOMO:CUSTOM_CODE_SLOT_START domainCreateResourceBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	dbUser, err := d.repo.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START domainCreateResourceAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return dbUser, nil
}
