package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// GetUser gets a user.
func (repo *Client) GetUser(ctx context.Context, m cmodel.User, fields []string) (cmodel.User, error) {
	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	mask := masks.Map(fields, model.UserMap)
	if len(mask) == 0 {
		mask = gmodel.UserFields.Mask()
	}

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		First(&gm).Error
	if err != nil {
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}
