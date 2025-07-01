package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// Recipe Access Methods

func (d *Domain) CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (recipeAccess model.RecipeAccess, err error) {
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// based on recipient, set state and verify that the requester has access to recipient
	if access.RecipeAccessParent.Recipient.UserId != 0 {
		access.State = types.AccessState_ACCESS_STATE_PENDING
	} else if access.RecipeAccessParent.Recipient.CircleId != 0 {
		access.State = types.AccessState_ACCESS_STATE_ACCEPTED
	} else {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, access.RecipeAccessParent.RecipeId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, access.RecipeAccessParent.RecipeId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	if authAccount.PermissionLevel < access.Level {
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "cannot create access with higher level than the requester's level"}
	}

	// create access
	access, err = d.repo.CreateRecipeAccess(ctx, access)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	return access, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	// verify access is for the given recipe
	if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
		return domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// check if the access is for the given user (they can delete their own access)
	isRecipient := (access.RecipeAccessParent.Recipient.UserId != 0 && access.RecipeAccessParent.Recipient.UserId == authAccount.UserId) ||
		(access.RecipeAccessParent.Recipient.CircleId != 0 && access.RecipeAccessParent.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		// user is not the recipient, so they need management permissions
		// get permission levels using same pattern as CreateRecipeAccess
		if authAccount.CircleId != 0 {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return err
			}
		} else {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return err
			}
		}

		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return domain.ErrPermissionDenied{Msg: "user does not have access to delete this recipe access"}
		}
	}

	return d.repo.DeleteRecipeAccess(ctx, parent, id)
}

func (d *Domain) GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access record
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	// verify access is for the given recipe
	if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// check if the access is for the given user (they can view their own access)
	isRecipient := (access.RecipeAccessParent.Recipient.UserId != 0 && access.RecipeAccessParent.Recipient.UserId == authAccount.UserId) ||
		(access.RecipeAccessParent.Recipient.CircleId != 0 && access.RecipeAccessParent.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		// user is not the recipient, so they need management permissions to view
		// get permission levels using same pattern as CreateRecipeAccess
		if authAccount.CircleId != 0 {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return model.RecipeAccess{}, err
			}
		} else {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return model.RecipeAccess{}, err
			}
		}

		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to view this recipe access"}
		}
	}

	return access, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) (recipeAccesses []model.RecipeAccess, err error) {
	if authAccount.UserId == 0 && authAccount.CircleId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.RecipeId.RecipeId != 0 {
		if authAccount.CircleId != 0 {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return nil, err
			}
		} else {
			authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, parent.RecipeId)
			if err != nil {
				return nil, err
			}
		}

		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return nil, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	return d.repo.ListRecipeAccesses(ctx, authAccount, parent, int32(pageSize), int64(pageOffset), filter)
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error) {
	// verify recipe is set
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if access.RecipeAccessId.RecipeAccessId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetRecipeAccess(ctx, access.RecipeAccessParent, access.RecipeAccessId)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	// verify access is for the given recipe
	if dbAccess.RecipeAccessParent.RecipeId.RecipeId != access.RecipeAccessParent.RecipeId.RecipeId {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// get requester's permission levels using same pattern as CreateRecipeAccess
	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, access.RecipeAccessParent.RecipeId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, access.RecipeAccessParent.RecipeId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
	}

	// verify requester has WRITE access to manage recipe access
	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to update recipe access"}
	}

	// if updating permission level, ensure it doesn't exceed the requester's level
	if access.Level > authAccount.PermissionLevel {
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "cannot update access level to higher than your own level"}
	}

	// update access
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, access)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	return updatedAccess, nil
}

func (d *Domain) AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	// verify the user is the recipient of this access
	isRecipient := (access.RecipeAccessParent.Recipient.UserId != 0 && access.RecipeAccessParent.Recipient.UserId == authAccount.UserId) ||
		(access.RecipeAccessParent.Recipient.CircleId != 0 && access.RecipeAccessParent.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, access)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	return updatedAccess, nil
}
