package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateUserAccess - Only allows WRITE level
func (d *Domain) CreateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (dbUserAccess model.UserAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateUserAccess called")
	if access.UserAccessParent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if access.Recipient.UserId == 0 {
		log.Warn().Msg("recipient is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	if access.PermissionLevel != types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("access level must be write")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access level must be write"}
	}

	authAccount.UserId = access.UserAccessParent.UserId.UserId

	// Prepare both A and B accesses
	user1 := authAccount.AuthUserId
	user2 := authAccount.UserId

	userA := model.UserAccess{
		UserAccessParent: model.UserAccessParent{UserId: model.UserId{UserId: user2}},
		PermissionLevel:  types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		State:            types.AccessState_ACCESS_STATE_PENDING,
		Requester:        model.UserId{UserId: user1},
		Recipient:        model.UserId{UserId: user1},
		AcceptTarget:     types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED,
	}

	userB := model.UserAccess{
		UserAccessParent: model.UserAccessParent{UserId: model.UserId{UserId: user1}},
		PermissionLevel:  types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		Requester:        model.UserId{UserId: user1},
		Recipient:        model.UserId{UserId: user2},
	}

	determinedUserAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(
		ctx, authAccount, userB,
		withMinimimRecipientPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access ownership details when creating a user access")
		return model.UserAccess{}, err
	}
	userB.AcceptTarget = determinedUserAccessOwnershipDetails.acceptTarget
	userB.State = determinedUserAccessOwnershipDetails.accessState

	// Create both accesses
	dbUserAccessA, errA := d.repo.CreateUserAccess(ctx, userA, nil)
	if errA != nil {
		log.Error().Err(errA).Msg("repo.CreateUserAccess (A) failed")
		return model.UserAccess{}, errA
	}
	_, errB := d.repo.CreateUserAccess(ctx, userB, nil)
	if errB != nil {
		log.Error().Err(errB).Msg("repo.CreateUserAccess (B) failed")
		// Attempt to clean up A if B fails
		_ = d.repo.DeleteUserAccess(ctx, userA.UserAccessParent, userA.UserAccessId)
		return model.UserAccess{}, errB
	}

	log.Info().Msg("Domain CreateUserAccess returning successfully (B)")
	return dbUserAccessA, nil
}

// DeleteUserAccess -
func (d *Domain) DeleteUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required when deleting a user access")
		return domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required when deleting a user access")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}
	dbAccess, err := d.repo.GetUserAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user access")
		return domain.ErrInternal{Msg: "unable to get user access"}
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access ownership details when deleting a user access")
		return domain.ErrInternal{Msg: "unable to determine recipe access ownership details"}
	}

	if !determinedRecipeAccessOwnershipDetails.isRecipientOwner && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when deleting a user access")
		return domain.ErrPermissionDenied{Msg: "access denied"}
	}

	// Delete the current access
	err = d.repo.DeleteUserAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete user access")
		return domain.ErrInternal{Msg: "unable to delete user access"}
	}
	// Delete the sister access (swap UserId and Recipient)
	sisterParent := model.UserAccessParent{UserId: dbAccess.Recipient}
	sisterFilter := fmt.Sprintf("requester_user_id=%d AND recipient_user_id=%d", dbAccess.Requester.UserId, parent.UserId.UserId)
	// Find the sister access id
	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, sisterParent, 1, 0, sisterFilter, []string{model.UserAccessField_Id})
	if err != nil {
		log.Error().Err(err).Msg("unable to list sister user access")
		return domain.ErrInternal{Msg: "unable to list sister user access"}
	}
	if len(accesses) == 0 {
		log.Warn().Msg("sister access not found on delete")
	}

	err = d.repo.DeleteUserAccess(ctx, sisterParent, accesses[0].UserAccessId)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete sister user access")
		return domain.ErrInternal{Msg: "unable to delete sister user access"}
	}

	return nil
}

// GetUserAccess -
func (d *Domain) GetUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId, fields []string) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required when getting a user access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required when getting a user access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetUserAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to get user access"}
	}

	determinedUserAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access ownership details when getting a user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to determine user access ownership details"}
	}

	if !determinedUserAccessOwnershipDetails.isRecipientOwner && !determinedUserAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when getting a user access")
		return model.UserAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	return dbAccess, nil
}

// ListUserAccesses -
func (d *Domain) ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("requester is required when listing user accesses")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.UserId.UserId != 0 {
		authAccount.UserId = parent.UserId.UserId
		_, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing user accesses")
			return nil, err
		}
	}

	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list user accesses")
		return nil, domain.ErrInternal{Msg: "unable to list user accesses"}
	}

	return accesses, nil
}

// AcceptUserAccess -
func (d *Domain) AcceptUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required when accepting a user access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required when accepting a user access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}
	dbAccess, err := d.repo.GetUserAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to get user access"}
	}
	if dbAccess.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	determinedRecipeAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access ownership details when accepting a user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to determine user access ownership details"}
	}

	if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RESOURCE && !determinedRecipeAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("must be resource owner to accept resourse targeted recipe access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "must be resource owner to accept resourse targeted recipe access"}
	} else if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT && !determinedRecipeAccessOwnershipDetails.isRecipientOwner {
		log.Warn().Msg("must be recipient owner to accept recipient targeted recipe access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "must be recipient owner to accept recipient targeted recipe access"}
	} else if determinedRecipeAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED {
		log.Warn().Msg("unspecified accept target when accepting a recipe access")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "unspecified accept target"}
	}

	// Accept B
	dbAccess.State = types.AccessState_ACCESS_STATE_ACCEPTED
	updated, err := d.repo.UpdateUserAccess(ctx, dbAccess, []string{model.UserAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to update user access"}
	}
	// Accept sister A (swap UserId and Recipient)
	sisterParent := model.UserAccessParent{UserId: dbAccess.Recipient}
	sisterFilter := fmt.Sprintf("requester_user_id=%d AND recipient_user_id=%d", dbAccess.Requester.UserId, parent.UserId.UserId)
	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, sisterParent, 1, 0, sisterFilter, []string{model.UserAccessField_Id})
	if err != nil {
		log.Error().Err(err).Msg("unable to list sister user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to list sister user access"}
	}
	if len(accesses) == 0 {
		log.Warn().Msg("sister access not found on accept")
	}

	accesses[0].State = types.AccessState_ACCESS_STATE_ACCEPTED
	_, err = d.repo.UpdateUserAccess(ctx, accesses[0], []string{model.UserAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update sister user access")
		return model.UserAccess{}, domain.ErrInternal{Msg: "unable to update user access"}
	}

	return updated, nil
}
