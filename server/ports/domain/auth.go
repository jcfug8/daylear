package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type authDomain interface {
	AuthorizeRecipeParent(ctx context.Context, token string, parent model.RecipeParent) error
	AuthorizeCircleParent(ctx context.Context, token string, parent model.CircleParent) error
}
