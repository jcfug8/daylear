package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// GetRecipeUserPermission gets a user's permission for a recipe.
func (repo *Client) GetRecipeUserPermission(ctx context.Context, userId int64, recipeId int64) (permPb.PermissionLevel, error) {
	var recipeUser gmodel.RecipeUser
	err := repo.db.WithContext(ctx).
		Where("user_id = ? AND recipe_id = ?", userId, recipeId).
		First(&recipeUser).Error
	if err != nil {
		return permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED, ConvertGormError(err)
	}

	return recipeUser.PermissionLevel, nil
}
