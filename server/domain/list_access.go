package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// List Access Methods
func (d *Domain) CreateListAccess(ctx context.Context, authAccount model.AuthAccount, access model.ListAccess) (listAccess model.ListAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.ListAccessParent.ListId.ListId == 0 {
		log.Warn().Msg("list id is required when creating a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	determinedListAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, access)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access ownership details when creating a list access")
		return model.ListAccess{}, err
	}

	access.Requester = model.ListRecipientOrRequester{
		UserId: authAccount.AuthUserId,
	}
	access.AcceptTarget = determinedListAccessOwnershipDetails.acceptTarget
	access.State = determinedListAccessOwnershipDetails.accessState

	if access.PermissionLevel > determinedListAccessOwnershipDetails.maximumPermissionLevel {
		log.Warn().Msg("unable to create list access with the given permission level")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "cannot create access level higher than your own level"}
	}

	// create access
	access, err = d.repo.CreateListAccess(ctx, access, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to create list access")
		return model.ListAccess{}, err
	}

	return access, nil
}

func (d *Domain) DeleteListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify list is set
	if parent.ListId.ListId == 0 {
		log.Warn().Msg("list id is required when deleting a list access")
		return domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	// verify access id is set
	if id.ListAccessId == 0 {
		log.Warn().Msg("access id is required when deleting a list access")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetListAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get list access")
		return err
	}

	determinedListAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access ownership details when deleting a list access")
		return domain.ErrInternal{Msg: "unable to determine list access ownership details"}
	}

	if !determinedListAccessOwnershipDetails.isRecipientOwner && !determinedListAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when deleting a list access")
		return domain.ErrPermissionDenied{Msg: "access denied"}
	}

	err = d.repo.DeleteListAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list access")
		return err
	}

	return nil
}

func (d *Domain) GetListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId, fields []string) (model.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.ListId.ListId == 0 {
		log.Warn().Msg("list id is required when getting a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	// verify access id is set
	if id.ListAccessId == 0 {
		log.Warn().Msg("access id is required when getting a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetListAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get list access when getting a list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to get list access"}
	}

	determinedListAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access ownership details when getting a list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to determine list access ownership details"}
	}

	if !determinedListAccessOwnershipDetails.isRecipientOwner && !determinedListAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when getting a list access")
		return model.ListAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	return dbAccess, nil
}

func (d *Domain) ListListAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) (listAccesses []model.ListAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Warn().Msg("requester is required when listing list accesses")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.ListId.ListId != 0 {
		_, err := d.determineListAccess(ctx, authAccount, parent.ListId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine list access")
			return nil, err
		}
	}

	listAccesses, err = d.repo.ListListAccesses(ctx, authAccount, parent, int32(pageSize), int64(pageOffset), filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list list accesses")
		return nil, domain.ErrInternal{Msg: "unable to list list accesses"}
	}

	return listAccesses, nil
}

func (d *Domain) UpdateListAccess(ctx context.Context, authAccount model.AuthAccount, access model.ListAccess, fields []string) (model.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.ListAccessParent.ListId.ListId == 0 {
		log.Warn().Msg("list id is required when updating a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	// verify access id is set
	if access.ListAccessId.ListAccessId == 0 {
		log.Warn().Msg("access id is required when updating a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetListAccess(ctx, access.ListAccessParent, access.ListAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to get list access"}
	}

	determinedListAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access ownership details when updating a list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to determine list access ownership details"}
	}

	if !determinedListAccessOwnershipDetails.isRecipientOwner && !determinedListAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when updating a list access")
		return model.ListAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	if slices.Contains(fields, model.ListAccessField_PermissionLevel) && determinedListAccessOwnershipDetails.maximumPermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update list access permission level to a higher level than your own")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "cannot update list access permission level to a higher level than your own"}
	}

	// update access
	updatedAccess, err := d.repo.UpdateListAccess(ctx, access, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to update list access"}
	}

	return updatedAccess, nil
}

// AcceptListAccess accepts a pending list access.
func (d *Domain) AcceptListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId) (model.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.ListId.ListId == 0 {
		log.Warn().Msg("list id is required when accepting a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	// verify access id is set
	if id.ListAccessId == 0 {
		log.Warn().Msg("access id is required when accepting a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	dbAccess, err := d.repo.GetListAccess(ctx, parent, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get list access")
		return model.ListAccess{}, err
	}

	// verify the access is in pending state
	if dbAccess.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	determinedListAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(
		ctx, authAccount, dbAccess,
		withAllowAutoOmitAccessChecks(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine list access ownership details when accepting a list access")
		return model.ListAccess{}, domain.ErrInternal{Msg: "unable to determine list access ownership details"}
	}

	if determinedListAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RESOURCE && !determinedListAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("must be resource owner to accept resourse targeted list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "must be resource owner to accept resourse targeted list access"}
	} else if determinedListAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT && !determinedListAccessOwnershipDetails.isRecipientOwner {
		log.Warn().Msg("must be recipient owner to accept recipient targeted list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "must be recipient owner to accept recipient targeted list access"}
	} else if determinedListAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED {
		log.Warn().Msg("unspecified accept target when accepting a list access")
		return model.ListAccess{}, domain.ErrInvalidArgument{Msg: "unspecified accept target"}
	}

	// update the access state to accepted
	dbAccess.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateListAccess(ctx, dbAccess, []string{model.ListAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update list access")
		return model.ListAccess{}, err
	}

	return updatedAccess, nil
}
