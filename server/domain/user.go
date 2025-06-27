package domain

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
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

	users, err := d.repo.ListUsers(ctx, model.AuthAccount{}, 1, 0, filter, model.UserFields.Mask())
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

// ListUsers lists users.
func (d *Domain) ListUsers(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string, fieldMask []string) (users []model.User, err error) {
	users, err = d.repo.ListUsers(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user.
func (d *Domain) UpdateUser(ctx context.Context, authAccount model.AuthAccount, user model.User, updateMask []string) (model.User, error) {
	if user.Id.UserId == 0 {
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	dbUser, err := d.repo.UpdateUser(ctx, user, updateMask)
	if err != nil {
		return model.User{}, err
	}

	return dbUser, nil
}

// GetUser gets a user.
func (d *Domain) GetUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId, fieldMask []string) (model.User, error) {
	if id.UserId == 0 {
		return model.User{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	user, err := d.repo.GetUser(ctx, id, fieldMask)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// DeleteUser deletes a user.
func (d *Domain) DeleteUser(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.User, error) {
	return model.User{}, nil
}

// CreateUser creates a new user.
func (d *Domain) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	if user.Email == "" {
		return model.User{}, domain.ErrInvalidArgument{Msg: "email is required"}
	}

	// if the username is not set, determine it
	if user.Username == "" {
		username, err := d.determineUsername(ctx, user.Email)
		if err != nil {
			return model.User{}, err
		}
		user.Username = username
	}

	dbUser, err := d.repo.CreateUser(ctx, user)
	if err != nil {
		d.log.Error().Err(err).Msg("failed to create user")
		return model.User{}, domain.ErrInternal{Msg: "failed to create user"}
	}

	return dbUser, nil
}

func (d *Domain) determineUsername(ctx context.Context, email string) (string, error) {
	emailParts := strings.Split(email, "@")
	if len(emailParts) != 2 {
		return "", domain.ErrInvalidArgument{Msg: "invalid email"}
	}
	baseUsername := emailParts[0]
	username := baseUsername

	for tryCount := 0; tryCount < 5; tryCount++ {
		filter := fmt.Sprintf("username = '%s'", username)
		dbUsers, err := d.repo.ListUsers(ctx, model.AuthAccount{}, 1, 0, filter, nil)
		if err != nil {
			d.log.Error().Err(err).Msg("failed to list users")
			return "", domain.ErrInternal{Msg: "failed to list users"}
		}
		if len(dbUsers) == 0 {
			break
		}

		username = fmt.Sprintf("%s_%d", baseUsername, rand.Intn(10000))
	}

	if username == "" {
		return "", domain.ErrInternal{Msg: "failed to determine username"}
	}

	return username, nil
}
