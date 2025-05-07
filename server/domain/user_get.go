package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	// IRIOMO:CUSTOM_CODE_SLOT_START domainUserGetImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// GetUser gets a user.
func (d *Domain) GetUser(ctx context.Context, id model.UserId, fieldMask []string) (model.User, error) {
	if id.UserId == 0 {
		return model.User{}, ErrInvalidArgument{Msg: "id required"}
	}

	user, err := d.repo.GetUser(ctx, model.User{
		Id: id,
	}, fieldMask)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
