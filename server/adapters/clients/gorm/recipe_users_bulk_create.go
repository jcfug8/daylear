package gorm

import (
	"context"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

func (repo *Client) BulkCreateRecipeUsers(ctx context.Context, recipeId cmodel.RecipeId, userIds []int64, permission pb.ShareRecipeRequest_ResourcePermission) error {
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
