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
	ListRecipes(context.Context, *model.PageToken[model.Recipe], string, []string) ([]model.Recipe, error)

	BulkCreateRecipeUsers(context.Context, model.RecipeId, []int64, permPb.PermissionLevel) error
	BulkDeleteRecipeUsers(context.Context, string) error
	GetRecipeUserPermission(ctx context.Context, userId int64, recipeId int64) (permPb.PermissionLevel, error)

	SetRecipeIngredients(context.Context, model.RecipeId, []model.IngredientGroup) error
	ListRecipeIngredients(context.Context, *model.PageToken[model.RecipeIngredient], string, []string) ([]model.RecipeIngredient, error)
	BulkDeleteRecipeIngredients(context.Context, string) ([]model.RecipeIngredient, error)
}
