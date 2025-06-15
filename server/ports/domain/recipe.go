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

	// depricated in favor of recipe access
	ShareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error
	// depricated in favor of recipe access
	UnshareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId) error
	// depricated in favor of recipe access
	ListRecipeRecipients(ctx context.Context, parent model.RecipeParent, id model.RecipeId) ([]model.RecipeRecipient, error)

	// Recipe Access methods
	CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) error
	GetRecipeAccess(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, parent model.RecipeParent, pageSize int64, pageOffset int64, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error)

	UploadRecipeImage(ctx context.Context, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error)
}
