package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type eventRecipeDomain interface {
	CreateEventRecipe(ctx context.Context, authAccount model.AuthAccount, eventRecipe model.EventRecipe) (model.EventRecipe, error)
	DeleteEventRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, id model.EventRecipeId) (model.EventRecipe, error)
	GetEventRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, id model.EventRecipeId, fields []string) (model.EventRecipe, error)
	ListEventRecipes(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, pageSize int32, offset int64, filter string, fields []string) ([]model.EventRecipe, error)
}
