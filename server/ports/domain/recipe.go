package domain

import (
	"context"
	"io"

	"github.com/jcfug8/daylear/server/core/file"
	model "github.com/jcfug8/daylear/server/core/model"
)

type recipeDomain interface {
	CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (model.Recipe, error)
	DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId) (model.Recipe, error)
	GetRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId, fields []string) (model.Recipe, error)
	ListRecipes(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, pageSize int32, offset int64, filter string, fields []string) ([]model.Recipe, error)
	UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, fields []string) (model.Recipe, error)

	ScrapeRecipe(ctx context.Context, authAccount model.AuthAccount, uri string) (model.Recipe, error)
	OCRRecipe(ctx context.Context, authAccount model.AuthAccount, imageReaders []io.Reader) (model.Recipe, error)

	UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error)
	GenerateRecipeImage(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId) (file.File, error)

	CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error)
	DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error
	GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
	ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.RecipeAccess, error)
	UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess, updateMask []string) (model.RecipeAccess, error)
	AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error)
}
