package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) IdentifyUser(ctx context.Context, user model.User) (model.User, error) {
	filter := ""
	switch {
	case user.AmazonId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.AmazonId, user.AmazonId)
	case user.GoogleId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.GoogleId, user.GoogleId)
	case user.FacebookId != "":
		filter = fmt.Sprintf(`%s = "%s"`, model.UserFields.FacebookId, user.FacebookId)
	}

	users, err := d.repo.ListUsers(ctx, nil, filter, model.UserFields.Mask())
	if err != nil {
		return model.User{}, err
	}

	if len(users) > 1 {
		return model.User{}, domain.ErrInternal{Msg: "Multiple users found"}
	}

	if len(users) == 0 {
		return model.User{}, domain.ErrNotFound{Msg: "User not found"}
	}

	return users[0], nil
}
