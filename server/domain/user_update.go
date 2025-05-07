package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// UpdateUser updates a user.
func (d *Domain) UpdateUser(ctx context.Context, user model.User, updateMask []string) (model.User, error) {
	if user.Id.UserId == 0 {
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	dbUser, err := d.repo.UpdateUser(ctx, user, updateMask)
	if err != nil {
		return model.User{}, err
	}

	return dbUser, nil
}
