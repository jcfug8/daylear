package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the user in the database.
type userClient interface {
	CreateUser(context.Context, model.User) (model.User, error)
	// DeleteUser(context.Context, model.User) (model.User, error)
	GetUser(context.Context, model.User, []string) (model.User, error)
	ListUsers(context.Context, *model.PageToken[model.User], string, []string) ([]model.User, error)
	// UpdateUser(context.Context, model.User, []string) (model.User, error)
}
