package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// Recipe Access Methods

func (d *Domain) CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error) {
	// TODO: Implement CreateRecipeAccess
	// Implementation commented out for refactoring
	/*
		// based on recipient, set state and verify that the requester has access to recipient
		if access.RecipeAccessParent.Recipient.UserId != 0 {
			access.State = pb.Access_STATE_PENDING
		} else if access.RecipeAccessParent.Recipient.CircleId != 0 {
			access.State = pb.Access_STATE_ACCEPTED
		} else {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
		}

		// verify permission is set
		if access.Level == permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "permission is required"}
		}

		accessLevel, err := d.verifyRecipeAccess(ctx, authAccount, access.RecipeAccessParent, access.RecipeAccessId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
		if accessLevel < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "not authorized to create access"}
		}
		if accessLevel < access.Level {
			return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "cannot create access with higher level than the requester's level"}
		}

		// create access
		access, err = d.repo.CreateRecipeAccess(ctx, access)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		return access, nil
	*/
	return model.RecipeAccess{}, domain.ErrInternal{Msg: "CreateRecipeAccess method not implemented"}
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	// TODO: Implement DeleteRecipeAccess
	// Implementation commented out for refactoring
	/*
		// verify recipe is set
		if parent.RecipeId.RecipeId == 0 {
			return domain.ErrInvalidArgument{Msg: "recipe id is required"}
		}

		// verify access id is set
		if id.RecipeAccessId == 0 {
			return domain.ErrInvalidArgument{Msg: "access id is required"}
		}

		access, err := d.repo.GetRecipeAccess(ctx, parent, id)
		if err != nil {
			return err
		}

		// verify access is for the given recipe
		if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
			return domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
		}

		// check if the access is for the given user
		if access.RecipeAccessParent.Recipient.UserId != authAccount.UserId && access.RecipeAccessParent.Recipient.CircleId != authAccount.CircleId {
			// verify the user has access to the recipe
			accessLevel, err := d.verifyRecipeAccess(ctx, authAccount, parent, id)
			if err != nil {
				return err
			}
			if accessLevel < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
				return domain.ErrInvalidArgument{Msg: "user does not have access to recipe"}
			}
		}

		return d.repo.DeleteRecipeAccess(ctx, parent, id)
	*/
	return domain.ErrInternal{Msg: "DeleteRecipeAccess method not implemented"}
}

func (d *Domain) GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	// TODO: Implement GetRecipeAccess
	// Implementation commented out for refactoring
	/*
		// verify recipe is set
		if parent.RecipeId.RecipeId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
		}

		// verify access id is set
		if id.RecipeAccessId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
		}

		access, err := d.repo.GetRecipeAccess(ctx, parent, id)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		// verify access is for the given recipe
		if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
		}

		// check if the access is for the given user
		if access.RecipeAccessParent.Recipient.UserId != authAccount.UserId && access.RecipeAccessParent.Recipient.CircleId != authAccount.CircleId {
			// verify the user has access to the recipe
			accessLevel, err := d.verifyRecipeAccess(ctx, authAccount, parent)
			if err != nil {
				return model.RecipeAccess{}, err
			}
			if accessLevel < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
				return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "user does not have access to recipe"}
			}
		}

		return access, nil
	*/
	return model.RecipeAccess{}, domain.ErrInternal{Msg: "GetRecipeAccess method not implemented"}
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.RecipeAccess, error) {
	// TODO: Implement ListRecipeAccesses
	// Implementation commented out for refactoring
	/*
		if parent.Requester.UserId == 0 && parent.Requester.CircleId == 0 {
			return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
		}

		return d.repo.ListRecipeAccesses(ctx, parent, int64(pageSize), int64(pageOffset), filter)
	*/
	return nil, domain.ErrInternal{Msg: "ListRecipeAccesses method not implemented"}
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (model.RecipeAccess, error) {
	// TODO: Implement UpdateRecipeAccess
	// Implementation commented out for refactoring
	/*
		// verify recipe is set
		if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
		}

		// verify access id is set
		if access.RecipeAccessId.RecipeAccessId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
		}

		// verify requester is set
		if access.RecipeAccessParent.Requester.UserId == 0 && access.RecipeAccessParent.Requester.CircleId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "requester is required"}
		}

		// verify access is for the given recipe
		dbAccess, err := d.repo.GetRecipeAccess(ctx, access.RecipeAccessParent, access.RecipeAccessId)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		if dbAccess.RecipeAccessParent.RecipeId.RecipeId != access.RecipeAccessParent.RecipeId.RecipeId {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
		}

		// verify requester has access to recipe
		var filter string
		// TODO: this looks wrong
		// based on recipient, set state and verify that the requester has access to recipient
		if access.RecipeAccessParent.Recipient.UserId != 0 {
			filter = fmt.Sprintf("user_id = %d", access.RecipeAccessParent.Recipient.UserId)
		} else if access.RecipeAccessParent.Recipient.CircleId != 0 {
			filter = fmt.Sprintf("circle_id = %d", access.RecipeAccessParent.Recipient.CircleId)
		} else {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
		}

		// verify requester has access to recipe
		recipeAccess, err := d.repo.ListRecipeAccesses(ctx, access.RecipeAccessParent, 1, 0, filter)
		if err != nil {
			return model.RecipeAccess{}, err
		}
		if len(recipeAccess) == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "requester does not have access to recipe"}
		}

		if recipeAccess[0].Level < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "requester does not have correct access to recipe"}
		}

		// currently only the state is allow to be updated to accepted
		if recipeAccess[0].Level > dbAccess.Level {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "cannot update level to higher than the requester's level"}
		}

		// update access
		access, err = d.repo.UpdateRecipeAccess(ctx, access)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		return access, nil
	*/
	return model.RecipeAccess{}, domain.ErrInternal{Msg: "UpdateRecipeAccess method not implemented"}
}

func (d *Domain) AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	// TODO: Implement AcceptRecipeAccess
	// Implementation commented out for refactoring
	/*
		// verify recipe is set
		if parent.RecipeId.RecipeId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
		}

		// verify access id is set
		if id.RecipeAccessId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
		}

		// verify requester is set (the user accepting the access)
		if parent.Requester.UserId == 0 && parent.Requester.CircleId == 0 {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "requester is required"}
		}

		// get the current access
		access, err := d.repo.GetRecipeAccess(ctx, parent, id)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		// verify the access is in pending state
		if access.State != pb.Access_STATE_PENDING {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
		}

		// verify the user is the recipient of this access
		if access.RecipeAccessParent.Recipient.UserId != parent.Requester.UserId && access.RecipeAccessParent.Recipient.CircleId != parent.Requester.CircleId {
			return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
		}

		// update the access state to accepted
		access.State = pb.Access_STATE_ACCEPTED

		// update access using the repository
		access, err = d.repo.UpdateRecipeAccess(ctx, access)
		if err != nil {
			return model.RecipeAccess{}, err
		}

		return access, nil
	*/
	return model.RecipeAccess{}, domain.ErrInternal{Msg: "AcceptRecipeAccess method not implemented"}
}
