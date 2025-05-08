package domain

import (
	"context"
	"io"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

type recipeDomain interface {
	CreateRecipe(ctx context.Context, recipe model.Recipe) (model.Recipe, error)
	DeleteRecipe(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.Recipe, error)
	GetRecipe(ctx context.Context, parent model.RecipeParent, id model.RecipeId, fieldMask []string) (model.Recipe, error)
	ListRecipes(ctx context.Context, page *model.PageToken[model.Recipe], parent model.RecipeParent, filter string, fieldMask []string) ([]model.Recipe, error)
	UpdateRecipe(ctx context.Context, recipe model.Recipe, updateMask []string) (model.Recipe, error)
	ShareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error
	UploadRecipeImage(ctx context.Context, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error)
}
