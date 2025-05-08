package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

func (repo *Client) BulkCreateRecipeUsers(ctx context.Context, recipeId cmodel.RecipeId, userIds []int64, permission permPb.PermissionLevel) error {
	if len(userIds) == 0 {
		return nil
	}

	recipeUsers := []gmodel.RecipeUser{}

	for _, userId := range userIds {
		recipeUsers = append(recipeUsers, gmodel.RecipeUser{
			RecipeId:        recipeId.RecipeId,
			UserId:          userId,
			PermissionLevel: permission,
		})
	}

	return repo.db.WithContext(ctx).Create(&recipeUsers).Error
}
