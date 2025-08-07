package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircleAccess creates a new circle access
func (d *Domain) CreateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (dbAccess model.CircleAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.CircleAccessParent.CircleId.CircleId == 0 {
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

	maxPermissionLevel := types.PermissionLevel_PERMISSION_LEVEL_READ
	access.State = types.AccessState_ACCESS_STATE_PENDING

	// if requesting access for a user other than yourself, ensure they have write access to the circle
	if authAccount.AuthUserId != access.Recipient.UserId {
		determinedCircleAccess, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: access.Requester.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access when creating a circle access")
			return model.CircleAccess{}, err
		}
		maxPermissionLevel = determinedCircleAccess.PermissionLevel
		// if the user has admin access to the user, set the access state to accepted since they could just accept the access themselves
		determinedUserAccess, err := d.repo.FindStandardUserUserAccess(ctx, authAccount, model.UserId{UserId: access.Recipient.UserId})
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access when creating a circle access")
			return model.CircleAccess{}, err
		}
		if determinedUserAccess.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
			access.State = types.AccessState_ACCESS_STATE_ACCEPTED
		}
	}

	if maxPermissionLevel < access.PermissionLevel {
		log.Warn().Msg("unable to create circle access with the given permission level")
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "unable to create circle access with the given permission level"}
	}

	// create access
	dbAccess, err = d.repo.CreateCircleAccess(ctx, access, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to create circle access"}
	}

	return access, nil
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

	// verify access is for the given circle
	if dbAccess.CircleAccessParent.CircleId.CircleId != parent.CircleId.CircleId {
		log.Warn().Msg("access is not for the given circle when deleting a circle access")
		return domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// check if the access is for the given user (they can delete their own access) or if they have write access to the circle
	if dbAccess.Recipient.UserId != authAccount.AuthUserId {
		_, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: parent.CircleId.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user'scircle access when deleting a circle access")
			return err
		}
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

	// get the access record
	access, err := d.repo.GetCircleAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access when getting a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to get circle access"}
	}

	// verify access is for the given circle
	if access.CircleAccessParent.CircleId.CircleId != parent.CircleId.CircleId {
		log.Warn().Msg("access is not for the given circle when getting a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// check if the access is for the given user (they can delete their own access) or if they have write access to the circle
	if access.Recipient.UserId != authAccount.AuthUserId {
		_, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: parent.CircleId.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user's circle access when getting a circle access")
			return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine user's circle access"}
		}
	}

	return access, nil
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

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetCircleAccess(ctx, access.CircleAccessParent, access.CircleAccessId, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle access when updating a circle access")
		return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to get circle access"}
	}

	maxPermissionLevel := dbAccess.PermissionLevel

	// verify access is for the given circle
	if dbAccess.CircleAccessParent.CircleId.CircleId != access.CircleAccessParent.CircleId.CircleId {
		log.Warn().Msg("access is not for the given circle when updating a circle access")
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// check if the access is for the given user (they can delete their own access) or if they have write access to the circle
	if dbAccess.Recipient.UserId != authAccount.AuthUserId {
		determinedAccess, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: access.CircleAccessParent.CircleId.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user's circle access when updating a circle access")
			return model.CircleAccess{}, domain.ErrInternal{Msg: "unable to determine user's circle access"}
		}
		maxPermissionLevel = determinedAccess.PermissionLevel
	}

	// check if the user is trying to increase the permission level beyond what they can grant
	if access.PermissionLevel > maxPermissionLevel {
		log.Warn().Msg("unable to update circle access with the given permission level")
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "unable to update circle access with the given permission level"}
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

	// verify the user is the recipient of this access
	if access.Recipient.UserId != authAccount.AuthUserId {
		log.Warn().Msg("only the recipient can accept a circle access")
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
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
