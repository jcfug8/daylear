package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// GetRecipeUser gets a recipe.
func (repo *Client) GetRecipeUserPermission(ctx context.Context, userId int64, recipeId int64) (pb.ShareRecipeRequest_ResourcePermission, error) {
	errz := errz.Context("repository.get_recipe")

	gm := gmodel.RecipeUser{
		RecipeId: recipeId,
		UserId:   userId,
	}

	err := repo.db.WithContext(ctx).
		Where(gm).Take(&gm).Error
	if err != nil {
		return pb.ShareRecipeRequest_RESOURCE_PERMISSION_UNSPECIFIED, ErrzError(errz, "", err)
	}

	return gm.PermissionLevel, nil
}
