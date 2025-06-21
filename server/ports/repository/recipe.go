package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// Client defines how to interact with the recipe in the database.
type recipeClient interface {
	CreateRecipe(context.Context, model.Recipe) (model.Recipe, error)
	DeleteRecipe(context.Context, model.Recipe) (model.Recipe, error)
	GetRecipe(context.Context, model.Recipe, []string) (model.Recipe, error)
	UpdateRecipe(context.Context, model.Recipe, []string) (model.Recipe, error)
	ListRecipes(context.Context, *model.PageToken[model.Recipe], model.RecipeParent, string, []string) ([]model.Recipe, error)

	// depricated in favor of recipe access
	GetRecipeRecipient(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.RecipeRecipient, error)
	// depricated in favor of recipe access
	ListRecipeRecipients(ctx context.Context, id model.RecipeId) ([]model.RecipeRecipient, error)
	// depricated in favor of recipe access
	BulkCreateRecipeRecipients(ctx context.Context, parents []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error
	// depricated in favor of recipe access
	BulkDeleteRecipeRecipients(ctx context.Context, parents []model.RecipeParent, id model.RecipeId) error

	CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, parent model.RecipeAccessParent, pageSize int64, pageOffset int64, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)

	SetRecipeIngredients(context.Context, model.RecipeId, []model.IngredientGroup) error
	ListRecipeIngredients(context.Context, *model.PageToken[model.RecipeIngredient], string, []string) ([]model.RecipeIngredient, error)
	BulkDeleteRecipeIngredients(context.Context, string) ([]model.RecipeIngredient, error)
}
