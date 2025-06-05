package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// GetRecipeCirclePermission gets a circle's permission for a recipe.
func (repo *Client) GetRecipeCirclePermission(ctx context.Context, circleId int64, recipeId int64) (permPb.PermissionLevel, error) {
	var recipeCircle gmodel.RecipeCircle
	err := repo.db.WithContext(ctx).
		Where("circle_id = ? AND recipe_id = ?", circleId, recipeId).
		First(&recipeCircle).Error
	if err != nil {
		return permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED, ConvertGormError(err)
	}

	return recipeCircle.PermissionLevel, nil
}
