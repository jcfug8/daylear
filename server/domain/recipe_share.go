package domain

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/domain"
)

// ShareRecipe deletes a recipe.
func (d *Domain) ShareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error {
	for _, parent := range parents {
		if parent.UserId == 0 {
			return domain.ErrInvalidArgument{Msg: "parent required"}
		}
	}

	if id.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	permission, err := d.repo.GetRecipeUserPermission(ctx, parent.UserId, id.RecipeId)
	if err != nil {
		return err
	}
	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
		return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	// filter out the the current parent
	userIds := []int64{}
	for _, p := range parents {
		if p.UserId == 0 || p.UserId == parent.UserId {
			continue
		}
		userIds = append(userIds, p.UserId)
	}
	if len(userIds) == 0 {
		return nil
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userIdsString := []string{}
	for _, parent := range parents {
		userIdsString = append(userIdsString, strconv.FormatInt(parent.UserId, 10))
	}
	filter := fmt.Sprintf("user_id = any(%s)", strings.Join(userIdsString, ","))
	err = tx.BulkDeleteRecipeUsers(ctx, filter)
	if err != nil {
		return err
	}

	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED {
		err = tx.BulkCreateRecipeUsers(ctx, id, userIds, permission)
		if err != nil {
			return err
		}
	}
	// TODO: could make create recipe user more of a set and make it
	// so that it doesn't error if the recipe user already exists

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
