package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type eventRecipeClient interface {
	CreateEventRecipe(ctx context.Context, eventRecipe model.EventRecipe, fields []string) (model.EventRecipe, error)
	DeleteEventRecipe(ctx context.Context, id model.EventRecipeId) (model.EventRecipe, error)
	GetEventRecipe(ctx context.Context, authAccount model.AuthAccount, id model.EventRecipeId, fields []string) (model.EventRecipe, error)
	ListEventRecipes(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, pageSize int32, offset int64, filter string, fields []string) ([]model.EventRecipe, error)
}
