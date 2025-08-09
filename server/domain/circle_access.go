package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircleAccess creates a new circle access
func (d *Domain) CreateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (dbAccess model.CircleAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.CircleId.CircleId == 0 {
		log.Warn().Msg("circle id required when creating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	determinedCircleAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, access)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access ownership details when creating a circle access")
		return model.CircleAccess{}, err
	}

	access.Requester = model.CircleRequester{
		UserId: authAccount.AuthUserId,
	}
	access.AcceptTarget = determinedCircleAccessOwnershipDetails.acceptTarget
	access.State = determinedCircleAccessOwnershipDetails.accessState

	if access.PermissionLevel > determinedCircleAccessOwnershipDetails.maximumPermissionLevel {
		log.Warn().Msg("unable to create circle access with the given permission level")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "cannot create access level higher than your own level"}
	}

	// create access
	dbAccess, err = d.repo.CreateCircleAccess(ctx, access, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to create circle access"}
	}

	return dbAccess, nil
}

// DeleteCircleAccess deletes a circle access
func (d *Domain) DeleteCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		log.Warn().Msg("circle id required when deleting a circle access")
		return domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		log.Warn().Msg("access id required when deleting a circle access")
		return domain.ErrInvalidArgument{Msg: "access id required"}
	}

	// get the existing access record
	dbAccess, err := d.repo.GetCircleAccess(ctx, parent, id, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access when deleting a circle access")
		return domain.ErrInternal{Msg: "unable to get circle access"}
	}

	determinedCircleAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access ownership details when deleting a circle access")
		return domain.ErrInternal{Msg: "unable to determine circle access ownership details"}
	}

	if !determinedCircleAccessOwnershipDetails.isRecipientOwner && !determinedCircleAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when deleting a circle access")
		return domain.ErrPermissionDenied{Msg: "access denied"}
	}

	err = d.repo.DeleteCircleAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete circle access when deleting a circle access")
		return domain.ErrInternal{Msg: "unable to delete circle access"}
	}

	return nil
}

// GetCircleAccess gets a circle access
func (d *Domain) GetCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId, fields []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		log.Warn().Msg("circle id required when getting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		log.Warn().Msg("access id required when getting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the dbAccess record
	dbAccess, err := d.repo.GetCircleAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access when getting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to get circle access"}
	}

	determinedCircleAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access ownership details when getting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine circle access ownership details"}
	}

	if !determinedCircleAccessOwnershipDetails.isRecipientOwner && !determinedCircleAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when getting a circle access")
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	return dbAccess, nil
}

// ListCircleAccesses lists circle accesses
func (d *Domain) ListCircleAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("requester is required when listing circle accesses")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.CircleId.CircleId != 0 {
		_, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: parent.CircleId.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user's circle access when listing circle accesses")
			return nil, err
		}
	}

	accesses, err := d.repo.ListCircleAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter, fields)
	if err != nil {
		return nil, err
	}

	return accesses, nil
}

// UpdateCircleAccess updates a circle access
func (d *Domain) UpdateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess, fields []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify circle is set
	if access.CircleAccessParent.CircleId.CircleId == 0 {
		log.Warn().Msg("circle id required when updating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	// verify access id is set
	if access.CircleAccessId.CircleAccessId == 0 {
		log.Warn().Msg("access id required when updating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetCircleAccess(ctx, access.CircleAccessParent, access.CircleAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access")
		return model.CircleAccess{}, err
	}

	determinedCircleAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access ownership details when updating a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine circle access ownership details"}
	}

	if !determinedCircleAccessOwnershipDetails.isRecipientOwner && !determinedCircleAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when updating a circle access")
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	if slices.Contains(fields, model.CircleAccessField_PermissionLevel) && determinedCircleAccessOwnershipDetails.maximumPermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update circle access permission level to a higher level than your own")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "cannot update circle access permission level to a higher level than your own"}
	}

	// update access
	updatedAccess, err := d.repo.UpdateCircleAccess(ctx, access, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to update circle access"}
	}

	return updatedAccess, nil
}

// AcceptCircleAccess accepts a pending circle access
func (d *Domain) AcceptCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		log.Warn().Msg("circle id required when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		log.Warn().Msg("access id required when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id required"}
	}

	// get the current access
	access, err := d.repo.GetCircleAccess(ctx, parent, id, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to get circle access"}
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("circle access must be in pending state to be accepted")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	determinedCircleAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(
		ctx, authAccount, access,
		withAllowAutoOmitAccessChecks(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access ownership details when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine circle access ownership details"}
	}

	if determinedCircleAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RESOURCE && !determinedCircleAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("must be resource owner to accept resourse targeted circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "must be resource owner to accept resourse targeted circle access"}
	} else if determinedCircleAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT && !determinedCircleAccessOwnershipDetails.isRecipientOwner {
		log.Warn().Msg("must be recipient owner to accept recipient targeted circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "must be recipient owner to accept recipient targeted circle access"}
	} else if determinedCircleAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED {
		log.Warn().Msg("unspecified accept target when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "unspecified accept target"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateCircleAccess(ctx, access, []string{model.CircleAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle access when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to update circle access"}
	}

	return updatedAccess, nil
}
