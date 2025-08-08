package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// Recipe Access Methods

func (d *Domain) CreateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess) (recipeAccess model.RecipeAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required when creating a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, access)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access ownership details when creating a recipe access")
		return model.RecipeAccess{}, err
	}

	access.Requester = model.RecipeRecipientOrRequester{
		UserId: authAccount.AuthUserId,
	}
	access.AcceptTarget = determinedRecipeAccessOwnershipDetails.acceptTarget
	access.State = determinedRecipeAccessOwnershipDetails.accessState

	if access.PermissionLevel > determinedRecipeAccessOwnershipDetails.maximumPermissionLevel {
		log.Warn().Msg("unable to create recipe access with the given permission level")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "cannot create access level higher than your own level"}
	}

	// create access
	access, err = d.repo.CreateRecipeAccess(ctx, access, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to create recipe access")
		return model.RecipeAccess{}, err
	}

	return access, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required when deleting a recipe access")
		return domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required when deleting a recipe access")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetRecipeAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe access")
		return err
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access ownership details when deleting a recipe access")
		return domain.ErrInternal{Msg: "unable to determine recipe access ownership details"}
	}

	if !determinedRecipeAccessOwnershipDetails.isRecipientOwner && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when deleting a recipe access")
		return domain.ErrPermissionDenied{Msg: "access denied"}
	}

	err = d.repo.DeleteRecipeAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete recipe access")
		return err
	}

	return nil
}

func (d *Domain) GetRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId, fields []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required when getting a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required when getting a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access record
	access, err := d.repo.GetRecipeAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe access")
		return model.RecipeAccess{}, err
	}

	// get the dbAccess record
	dbAccess, err := d.repo.GetRecipeAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe access when getting a recipe access")
		return model.RecipeAccess{}, domain.ErrInternal{Msg: "unable to get recipe access"}
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access ownership details when getting a recipe access")
		return model.RecipeAccess{}, domain.ErrInternal{Msg: "unable to determine recipe access ownership details"}
	}

	if !determinedRecipeAccessOwnershipDetails.isRecipientOwner && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when getting a recipe access")
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	return access, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) (recipeAccesses []model.RecipeAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Warn().Msg("requester is required when listing recipe accesses")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.RecipeId.RecipeId != 0 {
		_, err := d.determineRecipeAccess(ctx, authAccount, parent.RecipeId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine recipe access")
			return nil, err
		}
	}

	recipeAccesses, err = d.repo.ListRecipeAccesses(ctx, authAccount, parent, int32(pageSize), int64(pageOffset), filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list recipe accesses")
		return nil, err
	}

	return recipeAccesses, nil
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, authAccount model.AuthAccount, access model.RecipeAccess, fields []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required when updating a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if access.RecipeAccessId.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required when updating a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetRecipeAccess(ctx, access.RecipeAccessParent, access.RecipeAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe access")
		return model.RecipeAccess{}, err
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access ownership details when updating a recipe access")
		return model.RecipeAccess{}, domain.ErrInternal{Msg: "unable to determine recipe access ownership details"}
	}

	if !determinedRecipeAccessOwnershipDetails.isRecipientOwner && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when updating a recipe access")
		return model.RecipeAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	if slices.Contains(fields, model.RecipeAccessField_PermissionLevel) && determinedRecipeAccessOwnershipDetails.maximumPermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update recipe access permission level to a higher level than your own")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "cannot update recipe access permission level to a higher level than your own"}
	}

	// update access
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, access, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update recipe access")
		return model.RecipeAccess{}, err
	}

	return updatedAccess, nil
}

// AcceptRecipeAccess accepts a pending recipe access.
func (d *Domain) AcceptRecipeAccess(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.RecipeId.RecipeId == 0 {
		log.Warn().Msg("recipe id is required when accepting a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		log.Warn().Msg("access id is required when accepting a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	dbAccess, err := d.repo.GetRecipeAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get recipe access")
		return model.RecipeAccess{}, err
	}

	// verify the access is in pending state
	if dbAccess.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access ownership details when accepting a recipe access")
		return model.RecipeAccess{}, domain.ErrInternal{Msg: "unable to determine recipe access ownership details"}
	}

	if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RESOURCE && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("must be resource owner to accept resourse targeted recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "must be resource owner to accept resourse targeted recipe access"}
	} else if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT && !determinedRecipeAccessOwnershipDetails.isRecipientOwner {
		log.Warn().Msg("must be recipient owner to accept recipient targeted recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "must be recipient owner to accept recipient targeted recipe access"}
	} else if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED {
		log.Warn().Msg("unspecified accept target when accepting a recipe access")
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "unspecified accept target"}
	}

	// update the access state to accepted
	dbAccess.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateRecipeAccess(ctx, dbAccess, []string{model.RecipeAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update recipe access")
		return model.RecipeAccess{}, err
	}

	return updatedAccess, nil
}
