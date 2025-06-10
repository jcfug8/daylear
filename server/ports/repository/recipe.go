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

	GetRecipeRecipient(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.RecipeRecipient, error)
	ListRecipeRecipients(ctx context.Context, id model.RecipeId) ([]model.RecipeRecipient, error)
	BulkCreateRecipeRecipients(ctx context.Context, parents []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error
	BulkDeleteRecipeRecipients(ctx context.Context, parents []model.RecipeParent, id model.RecipeId) error

	SetRecipeIngredients(context.Context, model.RecipeId, []model.IngredientGroup) error
	ListRecipeIngredients(context.Context, *model.PageToken[model.RecipeIngredient], string, []string) ([]model.RecipeIngredient, error)
	BulkDeleteRecipeIngredients(context.Context, string) ([]model.RecipeIngredient, error)
}
