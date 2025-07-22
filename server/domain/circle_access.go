package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircleAccess creates a new circle access
func (d *Domain) CreateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateCircleAccess called")
	if access.CircleAccessParent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify recipient is set
	if access.Recipient == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	// Set the requester
	if authAccount.CircleId != 0 {
		access.Requester = model.CircleRequester{
			CircleId: authAccount.CircleId,
		}
	} else {
		access.Requester = model.CircleRequester{
			UserId: authAccount.AuthUserId,
		}
	}

	// For user recipients, access should be pending; for circle recipients (not applicable here), it would be accepted
	// Since CircleAccess.Recipient is an int64 (user ID), we assume it's always a user and set to pending
	access.State = types.AccessState_ACCESS_STATE_PENDING

	// Get requester's permission to the circle
	tempAuthAccount := authAccount
	tempAuthAccount.CircleId = access.CircleAccessParent.CircleId.CircleId
	permissionLevel, _, err := d.getCircleAccessLevels(ctx, tempAuthAccount)
	if err != nil {
		return model.CircleAccess{}, err
	}

	if permissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	if permissionLevel < access.PermissionLevel {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "cannot create access with higher level than the requester's level"}
	}

	// create access
	access, err = d.repo.CreateCircleAccess(ctx, access)
	if err != nil {
		return model.CircleAccess{}, err
	}

	log.Info().Msg("Domain CreateCircleAccess returning successfully")
	return access, nil
}

// DeleteCircleAccess deletes a circle access
func (d *Domain) DeleteCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain DeleteCircleAccess called")
	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record
	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	// verify access is for the given circle
	if access.CircleAccessParent.CircleId.CircleId != parent.CircleId.CircleId {
		return domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// check if the access is for the given user (they can delete their own access)
	if access.Recipient != authAccount.AuthUserId {
		// user is not the recipient, so they need management permissions
		// get permission levels for the circle
		tempAuthAccount := authAccount
		tempAuthAccount.CircleId = parent.CircleId.CircleId
		permissionLevel, _, err := d.getCircleAccessLevels(ctx, tempAuthAccount)
		if err != nil {
			return err
		}

		if permissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return domain.ErrPermissionDenied{Msg: "user does not have access to delete this circle access"}
		}
	}

	err = d.repo.DeleteCircleAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	log.Info().Msg("Domain DeleteCircleAccess returning successfully")
	return nil
}

// GetCircleAccess gets a circle access
func (d *Domain) GetCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain GetCircleAccess called")
	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access record
	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify access is for the given circle
	if access.CircleAccessParent.CircleId.CircleId != parent.CircleId.CircleId {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	if access.Recipient != authAccount.AuthUserId {
		// user is not the recipient, so they need management permissions to view
		// get permission levels for the circle
		tempAuthAccount := authAccount
		tempAuthAccount.CircleId = parent.CircleId.CircleId
		permissionLevel, _, err := d.getCircleAccessLevels(ctx, tempAuthAccount)
		if err != nil {
			return model.CircleAccess{}, err
		}

		if permissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to view this circle access"}
		}
	}

	log.Info().Msg("Domain GetCircleAccess returning successfully")
	return access, nil
}

// ListCircleAccesses lists circle accesses
func (d *Domain) ListCircleAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain ListCircleAccesses called")
	if authAccount.AuthUserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.CircleId.CircleId != 0 {
		// Check if user has permission to list accesses for this circle
		tempAuthAccount := authAccount
		tempAuthAccount.CircleId = parent.CircleId.CircleId
		permissionLevel, _, err := d.getCircleAccessLevels(ctx, tempAuthAccount)
		if err != nil {
			return nil, err
		}

		if permissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return nil, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	accesses, err := d.repo.ListCircleAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter)
	if err != nil {
		return nil, err
	}

	log.Info().Msg("Domain ListCircleAccesses returning successfully")
	return accesses, nil
}

// UpdateCircleAccess updates a circle access
func (d *Domain) UpdateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess, updateMask []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UpdateCircleAccess called")
	// verify circle is set
	if access.CircleAccessParent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if access.CircleAccessId.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the existing access record to verify it exists
	dbAccess, err := d.repo.GetCircleAccess(ctx, access.CircleAccessParent, access.CircleAccessId)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify access is for the given circle
	if dbAccess.CircleAccessParent.CircleId.CircleId != access.CircleAccessParent.CircleId.CircleId {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// get requester's permission levels for the circle
	tempAuthAccount := authAccount
	tempAuthAccount.CircleId = access.CircleAccessParent.CircleId.CircleId
	permissionLevel, _, err := d.getCircleAccessLevels(ctx, tempAuthAccount)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify requester has WRITE access to manage circle access
	if permissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to update circle access"}
	}

	// if updating permission level, ensure it doesn't exceed the requester's level
	if access.PermissionLevel > permissionLevel {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "cannot update access level to higher than your own level"}
	}

	// preserve read-only fields
	access.Requester = dbAccess.Requester
	access.Recipient = dbAccess.Recipient

	// update access
	updatedAccess, err := d.repo.UpdateCircleAccess(ctx, access, updateMask)
	if err != nil {
		return model.CircleAccess{}, err
	}

	log.Info().Msg("Domain UpdateCircleAccess returning successfully")
	return updatedAccess, nil
}

// AcceptCircleAccess accepts a pending circle access
func (d *Domain) AcceptCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain AcceptCircleAccess called")
	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	// verify the user is the recipient of this access
	if access.Recipient != authAccount.AuthUserId {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateCircleAccess(ctx, access, []string{model.CircleAccessFields.State})
	if err != nil {
		return model.CircleAccess{}, err
	}

	log.Info().Msg("Domain AcceptCircleAccess returning successfully")
	return updatedAccess, nil
}
