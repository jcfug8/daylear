package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// Recipe Access Methods

func (d *Domain) CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (recipeAccess model.RecipeAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateRecipeAccess called")
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// based on recipient, set state and verify that the requester has access to recipient
	if access.Recipient.UserId != 0 {
		access.State = types.AccessState_ACCESS_STATE_PENDING
	} else if access.Recipient.CircleId != 0 {
		access.State = types.AccessState_ACCESS_STATE_ACCEPTED
	} else {
		log.Warn().Msg("recipient is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkRecipeAccess(ctx, authAccount, access.RecipeAccessParent.RecipeId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
	if err != nil {
		log.Error().Err(err).Msg("checkRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	if authAccount.PermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot create access with higher level than the requester's level")
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "cannot create access with higher level than the requester's level"}
	}

	// create access
	access, err = d.repo.CreateRecipeAccess(ctx, access)
	if err != nil {
		log.Error().Err(err).Msg("repo.CreateRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	log.Info().Msg("Domain CreateRecipeAccess returning successfully")
	return access, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain DeleteRecipeAccess called")
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required")
		return domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipeAccess failed")
		return err
	}

	// verify access is for the given recipe
	if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
		log.Warn().Msg("access is not for the given recipe")
		return domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// check if the access is for the given user (they can delete their own access)
	isRecipient := (access.Recipient.UserId != 0 && access.Recipient.UserId == authAccount.AuthUserId) ||
		(access.Recipient.CircleId != 0 && access.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkRecipeAccess(ctx, authAccount, parent.RecipeId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
		if err != nil {
			log.Error().Err(err).Msg("checkRecipeAccess failed")
			return err
		}
	}

	err = d.repo.DeleteRecipeAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.DeleteRecipeAccess failed")
		return err
	}

	log.Info().Msg("Domain DeleteRecipeAccess returning successfully")
	return nil
}

func (d *Domain) GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain GetRecipeAccess called")
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access record
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	// verify access is for the given recipe
	if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
		log.Warn().Msg("access is not for the given recipe")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// check if the access is for the given user (they can view their own access)
	isRecipient := (access.Recipient.UserId != 0 && access.Recipient.UserId == authAccount.AuthUserId) ||
		(access.Recipient.CircleId != 0 && access.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkRecipeAccess(ctx, authAccount, parent.RecipeId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
		if err != nil {
			log.Error().Err(err).Msg("checkRecipeAccess failed")
			return model.RecipeAccess{}, err
		}
	}

	log.Info().Msg("Domain GetRecipeAccess returning successfully")
	return access, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string) (recipeAccesses []model.RecipeAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain ListRecipeAccesses called")
	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Warn().Msg("requester is required")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.RecipeId.RecipeId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkRecipeAccess(ctx, authAccount, parent.RecipeId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
		if err != nil {
			log.Error().Err(err).Msg("checkRecipeAccess failed")
			return nil, err
		}
	}

	recipeAccesses, err = d.repo.ListRecipeAccesses(ctx, authAccount, parent, int32(pageSize), int64(pageOffset), filter)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListRecipeAccesses failed")
		return nil, err
	}

	log.Info().Msg("Domain ListRecipeAccesses returning successfully")
	return recipeAccesses, nil
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess, updateMask []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UpdateRecipeAccess called")
	// verify recipe is set
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if access.RecipeAccessId.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetRecipeAccess(ctx, access.RecipeAccessParent, access.RecipeAccessId)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	// verify access is for the given recipe
	if dbAccess.RecipeAccessParent.RecipeId.RecipeId != access.RecipeAccessParent.RecipeId.RecipeId {
		log.Warn().Msg("access is not for the given recipe")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkRecipeAccess(ctx, authAccount, access.RecipeAccessParent.RecipeId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
	if err != nil {
		log.Error().Err(err).Msg("checkRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	// if updating permission level, ensure it doesn't exceed the requester's level
	if authAccount.PermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update access level to higher than your own level")
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "cannot update access level to higher than your own level"}
	}

	// update access
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, access, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	log.Info().Msg("Domain UpdateRecipeAccess returning successfully")
	return updatedAccess, nil
}

// AcceptRecipeAccess accepts a pending recipe access
func (d *Domain) AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain AcceptRecipeAccess called")
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	// verify the user is the recipient of this access
	isRecipient := (access.Recipient.UserId != 0 && access.Recipient.UserId == authAccount.AuthUserId) ||
		(access.Recipient.CircleId != 0 && access.Recipient.CircleId == authAccount.CircleId)

	if !isRecipient {
		log.Warn().Msg("only the recipient can accept this access")
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, access, []string{model.RecipeAccessFields.State})
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipeAccess failed")
		return model.RecipeAccess{}, err
	}

	log.Info().Msg("Domain AcceptRecipeAccess returning successfully")
	return updatedAccess, nil
}
