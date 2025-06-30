package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the recipe in the database.
type recipeClient interface {
	CreateRecipe(ctx context.Context, recipe model.Recipe) (model.Recipe, error)
	DeleteRecipe(ctx context.Context, id model.RecipeId) (model.Recipe, error)
	GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.Recipe, error)
	ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64) ([]model.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe model.Recipe, updateMask []string) (model.Recipe, error)

	CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)
}
