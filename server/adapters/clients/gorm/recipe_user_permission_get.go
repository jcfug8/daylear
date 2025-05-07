package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// GetRecipeUserPermission gets a user's permission for a recipe.
func (repo *Client) GetRecipeUserPermission(ctx context.Context, userId int64, recipeId int64) (pb.ShareRecipeRequest_ResourcePermission, error) {
	var recipeUser gmodel.RecipeUser
	err := repo.db.WithContext(ctx).
		Where("user_id = ? AND recipe_id = ?", userId, recipeId).
		First(&recipeUser).Error
	if err != nil {
		return pb.ShareRecipeRequest_RESOURCE_PERMISSION_UNSPECIFIED, ConvertGormError(err)
	}

	return recipeUser.PermissionLevel, nil
}
