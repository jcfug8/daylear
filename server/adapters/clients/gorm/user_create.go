package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START userCreateImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// CreateUser creates a new user.
func (repo *Client) CreateUser(ctx context.Context, m cmodel.User) (cmodel.User, error) {
	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	fields := masks.RemovePaths(gmodel.UserFields.Mask())

	err = repo.db.WithContext(ctx).
		Select(fields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}
