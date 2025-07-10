package token

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type Client interface {
	Encode(ctx context.Context, user model.User) (string, error)
	Decode(ctx context.Context, tn string) (model.User, error)
}
