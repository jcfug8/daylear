package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type authDomain interface {
	AuthorizeRecipeParent(ctx context.Context, token string, parent model.RecipeParent) error
	AuthorizeCircleParent(ctx context.Context, tokenUser model.User, parent model.CircleParent) error
}
