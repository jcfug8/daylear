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
	CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int32, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error)
	AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)

	UploadRecipeImage(ctx context.Context, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error)
}
