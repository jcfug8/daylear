package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/domain"
)

// ShareRecipe deletes a recipe.
func (d *Domain) ShareRecipe(ctx context.Context, parent model.RecipeParent, recipients []model.RecipeParent, id model.RecipeId, permission permPb.PermissionLevel) error {
	if id.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	if len(recipients) == 0 {
		return domain.ErrInvalidArgument{Msg: "recipients required"}
	}

	// Check if the user has permission to share the recipe
	recipient, err := d.repo.GetRecipeRecipient(ctx, parent, id)
	if err != nil {
		return err
	}
	if recipient.PermissionLevel != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
		return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}
	if parent.CircleId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}

	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// Remove existing shares for users
	err = tx.BulkDeleteRecipeRecipients(ctx, recipients, id)
	if err != nil {
		return err
	}

	// Create new shares
	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED {
		// TODO: check if the user or circle exists
		err = tx.BulkCreateRecipeRecipients(ctx, recipients, id, permission)
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

// UnshareRecipe removes sharing permissions for a recipe.
func (d *Domain) UnshareRecipe(ctx context.Context, parent model.RecipeParent, parents []model.RecipeParent, id model.RecipeId) error {
	if id.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "id required"}
	}

	if len(parents) == 0 {
		return domain.ErrInvalidArgument{Msg: "recipients required"}
	}

	// Allow self-unsharing regardless of permissions
	isSelfUnshare := false
	for _, p := range parents {
		if p.UserId == parent.UserId && p.CircleId == 0 {
			isSelfUnshare = true
			break
		}
	}

	// If not self-unsharing, check permissions
	if !isSelfUnshare {
		recipient, err := d.repo.GetRecipeRecipient(ctx, parent, id)
		if err != nil {
			return err
		}
		if recipient.PermissionLevel != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}
		if parent.CircleId != 0 {
			permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
			if err != nil {
				return err
			}
			if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
				return domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
			}
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Remove existing shares for users
	err = tx.BulkDeleteRecipeRecipients(ctx, parents, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
