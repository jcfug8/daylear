package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type authDomain interface {
	AuthorizeParent(ctx context.Context, token string, parent model.RecipeParent) error
}
