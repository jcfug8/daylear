package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START userCreateImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// CreateUser creates a new user.
func (repo *Client) CreateUser(ctx context.Context, m cmodel.User) (cmodel.User, error) {
	errz := errz.Context("repository.create_user")

	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, errz.Wrapf("invalid user: %v", err)
	}

	fields := masks.RemovePaths(
		gmodel.UserFields.Mask(),
	// IRIOMO:CUSTOM_CODE_SLOT_START createUserOmmitedFields
	// IRIOMO:CUSTOM_CODE_SLOT_END
	)

	// IRIOMO:CUSTOM_CODE_SLOT_START createUserBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	err = repo.db.WithContext(ctx).
		Select(fields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.User{}, ErrzError(errz, "", err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, errz.Wrapf("unable to read user: %v", err)
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START createUserAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return m, nil
}
