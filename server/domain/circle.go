package domain

import (
	"context"
	// "fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircle creates a new circle.
func (d *Domain) CreateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle) (model.Circle, error) {
	if authAccount.UserId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	circle.Id.CircleId = 0

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Circle{}, err
	}
	defer tx.Rollback()

	dbCircle, err := tx.CreateCircle(ctx, circle)
	if err != nil {
		return model.Circle{}, err
	}

	dbCircle.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	circleAccess := model.CircleAccess{
		CircleAccessParent: model.CircleAccessParent{
			CircleId: dbCircle.Id,
		},
		Requester: model.AuthAccount{
			UserId: authAccount.UserId,
		},
		Recipient: authAccount.UserId,
		Level:     types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:     types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	_, err = tx.CreateCircleAccess(ctx, circleAccess)
	if err != nil {
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Circle{}, err
	}

	return dbCircle, nil
}

// DeleteCircle deletes a circle.
func (d *Domain) DeleteCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.Circle, error) {
	// TODO: Implement DeleteCircle
	// Implementation commented out for refactoring
	/*
		if parent.UserId == 0 || id.CircleId == 0 {
			return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent and id required"}
		}

		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, id.CircleId)
		if err != nil {
			return model.Circle{}, err
		}
		if permission != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}

		tx, err := d.repo.Begin(ctx)
		if err != nil {
			return model.Circle{}, err
		}
		defer tx.Rollback()

		circle := model.Circle{Id: id, Parent: parent}
		deleted, err := tx.DeleteCircle(ctx, circle)
		if err != nil {
			return model.Circle{}, err
		}

		err = tx.BulkDeleteCircleUsers(ctx, fmt.Sprintf("circle_id = %d", id.CircleId))
		if err != nil {
			return model.Circle{}, err
		}

		err = tx.Commit()
		if err != nil {
			return model.Circle{}, err
		}

		return deleted, nil
	*/
	return model.Circle{}, domain.ErrInternal{Msg: "DeleteCircle method not implemented"}
}

// GetCircle gets a circle.
func (d *Domain) GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	if authAccount.UserId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	authAccount.VisibilityLevel = types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN

	authAccount.CircleId = id.CircleId

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.verifyCircleAccess(ctx, authAccount)
	if err != nil {
		return model.Circle{}, err
	}

	if authAccount.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	circle, err = d.repo.GetCircle(ctx, authAccount, id)
	if err != nil {
		return model.Circle{}, err
	}

	return circle, nil
}

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string, fieldMask []string) ([]model.Circle, error) {
	if authAccount.UserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	authAccount.VisibilityLevel = types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN

	circles, err := d.repo.ListCircles(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	return circles, nil
}

// UpdateCircle updates a circle.
func (d *Domain) UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, updateMask []string) (model.Circle, error) {
	// TODO: Implement UpdateCircle
	// Implementation commented out for refactoring
	/*
		if circle.Parent.UserId == 0 {
			return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
		}

		if circle.Id.CircleId == 0 {
			return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
		}

		permission, err := d.repo.GetCircleUserPermission(ctx, circle.Parent.UserId, circle.Id.CircleId)
		if err != nil {
			return model.Circle{}, err
		}
		if permission < permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}

		tx, err := d.repo.Begin(ctx)
		if err != nil {
			return model.Circle{}, err
		}
		defer tx.Rollback()

		updated, err := tx.UpdateCircle(ctx, circle, updateMask)
		if err != nil {
			return model.Circle{}, err
		}
		updated.Parent = circle.Parent

		err = tx.Commit()
		if err != nil {
			return model.Circle{}, err
		}

		return updated, nil
	*/
	return model.Circle{}, domain.ErrInternal{Msg: "UpdateCircle method not implemented"}
}
