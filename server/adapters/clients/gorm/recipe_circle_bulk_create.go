package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

func (repo *Client) BulkCreateRecipeCircles(ctx context.Context, recipeId cmodel.RecipeId, circleIds []int64, permission permPb.PermissionLevel) error {
	if len(circleIds) == 0 {
		return nil
	}

	recipeCircles := []gmodel.RecipeCircle{}

	for _, circleId := range circleIds {
		recipeCircles = append(recipeCircles, gmodel.RecipeCircle{
			RecipeId:        recipeId.RecipeId,
			CircleId:        circleId,
			PermissionLevel: permission,
		})
	}

	return repo.db.WithContext(ctx).Create(&recipeCircles).Error
}
