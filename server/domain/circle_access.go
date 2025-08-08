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

	// verify recipient is set
	if access.Recipient.UserId == 0 {
		log.Warn().Msg("recipient is required when creating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	access.Requester = model.CircleRequester{
		UserId: authAccount.AuthUserId,
	}

	var recipientOwner bool
	var resourceOwner bool

	dbCircle, err := d.repo.GetCircle(ctx, authAccount, access.CircleId, []string{model.CircleField_Visibility})
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle when creating a circle access")
		return model.CircleAccess{}, err
	}
	determinedCircleAccess, err := d.determineCircleAccess(
		ctx, authAccount, model.CircleId{CircleId: access.CircleId.CircleId},
		withResourceVisibilityLevel(dbCircle.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access when creating a circle access")
		return model.CircleAccess{}, err
	}
	resourceOwner = determinedCircleAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE

	if access.Recipient.UserId != 0 {
		determinedUserAccess, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: access.Recipient.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access when creating a circle access")
			return model.CircleAccess{}, err
		}
		recipientOwner = determinedUserAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	} else {
		log.Warn().Msg("recipient is required when creating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	if !resourceOwner && recipientOwner {
		access.State = types.AccessState_ACCESS_STATE_PENDING
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_RESOURCE
	} else if resourceOwner && !recipientOwner {
		access.State = types.AccessState_ACCESS_STATE_PENDING
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_RECIPIENT
	} else if resourceOwner && recipientOwner {
		access.State = types.AccessState_ACCESS_STATE_ACCEPTED
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED
	} else {
		log.Warn().Msg("unable to determine access state when creating a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine access state"}
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

	isRecipientOwner := false
	isCircleOwner := false

	determinedUserAccess, err := d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: dbAccess.Recipient.UserId},
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access when deleting a circle access")
		return err
	}

	isRecipientOwner = determinedUserAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	determinedCircleAccess, err := d.determineCircleAccess(
		ctx, authAccount, model.CircleId{CircleId: dbAccess.CircleId.CircleId},
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access when deleting a circle access")
		return err
	}

	isCircleOwner = determinedCircleAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	if !isRecipientOwner && !isCircleOwner {
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

	isRecipientOwner := false
	isCircleOwner := false

	determinedUserAccess, err := d.determineUserAccess(
		ctx, authAccount, model.UserId{UserId: dbAccess.Recipient.UserId},
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine user access when deleting a circle access")
		return model.CircleAccess{}, err
	}

	isRecipientOwner = determinedUserAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	determinedCircleAccess, err := d.determineCircleAccess(
		ctx, authAccount, model.CircleId{CircleId: dbAccess.CircleId.CircleId},
		withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access when deleting a circle access")
		return model.CircleAccess{}, err
	}

	isCircleOwner = determinedCircleAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	if !isRecipientOwner && !isCircleOwner {
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

	determinedCircleAccess, err := d.determineCircleAccess(
		ctx, authAccount, model.CircleId{CircleId: access.CircleAccessParent.CircleId.CircleId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine circle access when updating a circle access")
		return model.CircleAccess{}, err
	}

	if slices.Contains(fields, model.CircleAccessField_PermissionLevel) && determinedCircleAccess.PermissionLevel < access.PermissionLevel {
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

	switch access.AcceptTarget {
	case types.AcceptTarget_ACCEPT_TARGET_RESOURCE:
		_, err := d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: access.CircleId.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access when accepting circle access")
			return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine circle access"}
		}
	case types.AcceptTarget_ACCEPT_TARGET_RECIPIENT:
		if access.Recipient.UserId != 0 {
			_, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: access.Recipient.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
			if err != nil {
				log.Error().Err(err).Msg("unable to determine user access when accepting circle access")
				return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine user access"}
			}
		} else {
			log.Warn().Msg("recipient is required when accepting a circle access")
			return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
		}
	default:
		log.Warn().Msg("invalid accept target when accepting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "invalid accept target"}
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
