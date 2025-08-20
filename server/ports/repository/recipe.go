package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the recipe in the database.
type recipeClient interface {
	CreateRecipe(ctx context.Context, recipe model.Recipe, fields []string) (model.Recipe, error)
	DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.Recipe, error)
	GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, fields []string) (model.Recipe, error)
	ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]model.Recipe, error)
	UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, fields []string) (model.Recipe, error)

	FindStandardUserRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, error)
	FindDelegatedCircleRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, model.CircleAccess, error)
	FindDelegatedUserRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, model.UserAccess, error)

	CreateRecipeAccess(ctx context.Context, access model.RecipeAccess, fields []string) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	BulkDeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent) error
	GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId, fields []string) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess, fields []string) (model.RecipeAccess, error)

	// Favoriting methods
	CreateRecipeFavorite(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) error
	DeleteRecipeFavorite(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) error
	BulkDeleteRecipeFavorites(ctx context.Context, id model.RecipeId) error
}
