package domain

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircleAccess creates a new circle access
func (d *Domain) CreateCircleAccess(ctx context.Context, access model.CircleAccess) (model.CircleAccess, error) {
	// verify requester is set
	if access.CircleAccessParent.Requester.UserId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	// verify circle is set
	if access.CircleAccessParent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify recipient is set
	if access.CircleAccessParent.Recipient != 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	// verify requester has write permission to the circle
	permission, err := d.repo.GetCircleUserPermission(ctx, access.CircleAccessParent.Requester.UserId, access.CircleAccessParent.CircleId.CircleId)
	if err != nil {
		return model.CircleAccess{}, err
	}
	if permission < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "requester does not have write permission to circle"}
	}

	// verify requester cannot grant higher permission than they have
	if access.Level > permission {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "cannot grant permission higher than your own"}
	}

	// create access using the repository
	createdAccess, err := d.repo.CreateCircleAccess(ctx, access)
	if err != nil {
		return model.CircleAccess{}, err
	}

	return createdAccess, nil
}

// DeleteCircleAccess deletes a circle access
func (d *Domain) DeleteCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) error {
	// verify auth account is set
	if authAccount.UserId == 0 {
		return domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access to verify ownership or permissions
	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	// check if the user can delete this access:
	// 1. User is the recipient of the access (can delete their own access)
	// 2. User has write permission to the circle
	canDelete := false

	// Check if user is the recipient
	if access.CircleAccessParent.Recipient.UserId == authAccount.UserId ||
		(access.CircleAccessParent.Recipient.CircleId != 0 && access.CircleAccessParent.Recipient.CircleId == authAccount.CircleId) {
		canDelete = true
	} else {
		// Check if user has write permission to the circle
		permission, err := d.repo.GetCircleUserPermission(ctx, authAccount.UserId, parent.CircleId.CircleId)
		if err != nil && !errors.Is(err, domain.ErrNotFound{}) {
			return err
		}
		if permission >= permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			canDelete = true
		}
	}

	if !canDelete {
		return domain.ErrPermissionDenied{Msg: "user does not have permission to delete this access"}
	}

	// delete access using the repository
	err = d.repo.DeleteCircleAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	return nil
}

// GetCircleAccess gets a circle access
func (d *Domain) GetCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify access is for the given circle
	if access.CircleAccessParent.CircleId.CircleId != parent.CircleId.CircleId {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not for the given circle"}
	}

	// check if the access is for the given user or if user has access to the circle
	if access.CircleAccessParent.Recipient.UserId != authAccount.UserId && access.CircleAccessParent.Recipient.CircleId != authAccount.CircleId {
		// verify the user has access to the circle
		permission, err := d.repo.GetCircleUserPermission(ctx, authAccount.UserId, parent.CircleId.CircleId)
		if err != nil {
			return model.CircleAccess{}, err
		}
		if permission < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to circle"}
		}
	}

	return access, nil
}

// ListCircleAccesses lists circle accesses
func (d *Domain) ListCircleAccesses(ctx context.Context, parent model.CircleAccessParent, pageSize int32, pageOffset int32, filter string) ([]model.CircleAccess, error) {
	if parent.Requester.UserId == 0 && parent.Requester.CircleId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	return d.repo.ListCircleAccesses(ctx, parent, int64(pageSize), int64(pageOffset), filter)
}

// UpdateCircleAccess updates a circle access
func (d *Domain) UpdateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (model.CircleAccess, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle is set
	if access.CircleAccessParent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if access.CircleAccessId.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get existing access to verify permissions
	existingAccess, err := d.repo.GetCircleAccess(ctx, access.CircleAccessParent, access.CircleAccessId)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify user has write permission to the circle
	permission, err := d.repo.GetCircleUserPermission(ctx, authAccount.UserId, access.CircleAccessParent.CircleId.CircleId)
	if err != nil {
		return model.CircleAccess{}, err
	}
	if permission < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user does not have write permission to circle"}
	}

	// verify user cannot grant higher permission than they have
	if access.Level > permission {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "cannot grant permission higher than your own"}
	}

	// preserve read-only fields
	access.CircleAccessParent.Requester = existingAccess.CircleAccessParent.Requester
	access.CircleAccessParent.Recipient = existingAccess.CircleAccessParent.Recipient

	// update access using the repository
	updatedAccess, err := d.repo.UpdateCircleAccess(ctx, access)
	if err != nil {
		return model.CircleAccess{}, err
	}

	return updatedAccess, nil
}

// AcceptCircleAccess accepts a pending circle access
func (d *Domain) AcceptCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle is set
	if parent.CircleId.CircleId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify access id is set
	if id.CircleAccessId == 0 {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the access
	access, err := d.repo.GetCircleAccess(ctx, parent, id)
	if err != nil {
		return model.CircleAccess{}, err
	}

	// verify user is the recipient of the access
	if access.CircleAccessParent.Recipient.UserId != authAccount.UserId && access.CircleAccessParent.Recipient.CircleId != authAccount.CircleId {
		return model.CircleAccess{}, domain.ErrPermissionDenied{Msg: "user can only accept their own access"}
	}

	// verify access is pending
	if access.State != permPb.Access_STATE_PENDING {
		return model.CircleAccess{}, domain.ErrInvalidArgument{Msg: "access is not pending"}
	}

	// update state to accepted
	access.State = permPb.Access_STATE_ACCEPTED

	// update access using the repository
	acceptedAccess, err := d.repo.UpdateCircleAccess(ctx, access)
	if err != nil {
		return model.CircleAccess{}, err
	}

	return acceptedAccess, nil
}

func (d *Domain) verifyCircleAccess(ctx context.Context, authAccount model.AuthAccount) (permPb.PermissionLevel, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle id is set
	if authAccount.CircleId == 0 {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify circle exists
	circle, err := d.repo.GetCircle(ctx, model.Circle{Id: model.CircleId{CircleId: authAccount.CircleId}}, []string{model.CircleFields.Id})
	if err != nil {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	// check if user is a member of the circle
	permissionLevel, err := d.repo.GetCircleUserPermission(ctx, authAccount.UserId, authAccount.CircleId)
	if errors.Is(err, domain.ErrNotFound{}) && circle.Visibility != permPb.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user not a member of the circle"}
	} else if err != nil && !errors.Is(err, domain.ErrNotFound{}) {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	return permissionLevel, nil
}
