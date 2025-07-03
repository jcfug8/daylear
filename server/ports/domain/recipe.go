package domain

import (
	"context"
	"io"

	model "github.com/jcfug8/daylear/server/core/model"
)

type recipeDomain interface {
	CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (model.Recipe, error)
	DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.Recipe, error)
	GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.Recipe, error)
	ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string) ([]model.Recipe, error)
	UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, updateMask []string) (model.Recipe, error)
	AcceptRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) error
	ScrapeRecipe(ctx context.Context, authAccount model.AuthAccount, uri string) (model.Recipe, error)

	UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, imageReader io.Reader) (imageURI string, err error)

	CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error)
}
