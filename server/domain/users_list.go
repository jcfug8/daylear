package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

// ListUsers lists users.
func (d *Domain) ListUsers(ctx context.Context, page *model.PageToken[model.User], filter string, fieldMask []string) (users []model.User, err error) {
	users, err = d.repo.ListUsers(ctx, page, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	return users, nil
}
