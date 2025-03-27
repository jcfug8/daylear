package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// GetUser gets a user.
func (repo *Client) GetUser(ctx context.Context, m cmodel.User, fields []string) (cmodel.User, error) {
	errz := errz.Context("repository.get_user")

	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, errz.Wrapf("invalid user: %v", err)
	}

	err = repo.db.WithContext(ctx).
		Select(masks.Map(fields, gmodel.UserMap)).
		Take(&gm).Error
	if err != nil {
		return cmodel.User{}, ErrzError(errz, "", err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, errz.Wrapf("unable to read user: %v", err)
	}

	return m, nil
}
