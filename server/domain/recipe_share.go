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
	if id.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	// Check if the user has permission to share the recipe
	if parent.CircleId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}
	} else {
		permission, err := d.repo.GetRecipeUserPermission(ctx, parent.UserId, id.RecipeId)
		if err != nil {
			return err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}
	}

	// Separate users and circles
	userIds := []int64{}
	circleIds := []int64{}
	for _, p := range parents {
		if p.CircleId != 0 {
			circleIds = append(circleIds, p.CircleId)
		} else if p.UserId != 0 {
			userIds = append(userIds, p.UserId)
		}
	}

	if len(userIds) == 0 && len(circleIds) == 0 {
		return nil
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove existing shares for users
	if len(userIds) > 0 {
		userIdsString := []string{}
		for _, userId := range userIds {
			userIdsString = append(userIdsString, strconv.FormatInt(userId, 10))
		}
		filter := fmt.Sprintf("user_id = any(%s)", strings.Join(userIdsString, ","))
		err = tx.BulkDeleteRecipeUsers(ctx, filter)
		if err != nil {
			return err
		}
	}

	// Remove existing shares for circles
	if len(circleIds) > 0 {
		circleIdsString := []string{}
		for _, circleId := range circleIds {
			circleIdsString = append(circleIdsString, strconv.FormatInt(circleId, 10))
		}
		filter := fmt.Sprintf("circle_id = any(%s)", strings.Join(circleIdsString, ","))
		err = tx.BulkDeleteRecipeCircles(ctx, filter)
		if err != nil {
			return err
		}
	}

	// Create new shares
	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED {
		if len(userIds) > 0 {
			err = tx.BulkCreateRecipeUsers(ctx, id, userIds, permission)
			if err != nil {
				return err
			}
		}
		if len(circleIds) > 0 {
			err = tx.BulkCreateRecipeCircles(ctx, id, circleIds, permission)
			if err != nil {
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UnshareRecipe removes sharing permissions for a recipe.
func (d *Domain) UnshareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId) error {
	if parent.UserId == 0 {
		return domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	if len(parents) == 0 {
		return domain.ErrInvalidArgument{Msg: "recipients required"}
	}

	// Check if the user has permission to share the recipe
	if parent.CircleId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}
	} else {
		permission, err := d.repo.GetRecipeUserPermission(ctx, parent.UserId, id.RecipeId)
		if err != nil {
			return err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Collect all user IDs and circle IDs
	userIds := []int64{}
	circleIds := []int64{}
	for _, recipient := range parents {
		if recipient.CircleId != 0 {
			circleIds = append(circleIds, recipient.CircleId)
		} else if recipient.UserId != 0 {
			userIds = append(userIds, recipient.UserId)
		}
	}

	// Remove sharing permissions for users
	if len(userIds) > 0 {
		userIdsString := []string{}
		for _, userId := range userIds {
			userIdsString = append(userIdsString, strconv.FormatInt(userId, 10))
		}
		filter := fmt.Sprintf("recipe_id = %d AND user_id = any(%s)", id.RecipeId, strings.Join(userIdsString, ","))
		err = tx.BulkDeleteRecipeUsers(ctx, filter)
		if err != nil {
			return err
		}
	}

	// Remove sharing permissions for circles
	if len(circleIds) > 0 {
		circleIdsString := []string{}
		for _, circleId := range circleIds {
			circleIdsString = append(circleIdsString, strconv.FormatInt(circleId, 10))
		}
		filter := fmt.Sprintf("recipe_id = %d AND circle_id = any(%s)", id.RecipeId, strings.Join(circleIdsString, ","))
		err = tx.BulkDeleteRecipeCircles(ctx, filter)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
