package domain

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

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
		dbUser, err := d.repo.ListUsers(ctx, nil, filter, nil)
		if err != nil {
			d.log.Error().Err(err).Msg("failed to list users")
			return "", domain.ErrInternal{Msg: "failed to list users"}
		}
		if len(dbUser) == 0 {
			break
		}

		username = fmt.Sprintf("%s_%d", baseUsername, rand.Intn(10000))
	}

	if username == "" {
		return "", domain.ErrInternal{Msg: "failed to determine username"}
	}

	return username, nil
}
