package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type tokenDomain interface {
	// CreateToken creates a token for the given user and return a key
	// that can be used to retrieve the token later.
	CreateToken(ctx context.Context, token model.User) (string, error)
	// ParseToken parses the token and returns the user associated with it.
	ParseToken(ctx context.Context, token string) (model.User, error)
	// RetrieveToken retrieves the token associated with the given key then
	// deletes the token from the store.
	RetrieveToken(ctx context.Context, key string) (string, error)
}
