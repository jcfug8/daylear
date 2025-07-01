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
func (d *Domain) DeleteCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	if authAccount.UserId == 0 || id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent and id required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		return model.Circle{}, err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Circle{}, err
	}
	defer tx.Rollback()

	circle, err = tx.DeleteCircle(ctx, id)
	if err != nil {
		return model.Circle{}, err
	}

	err = tx.BulkDeleteCircleAccess(ctx, model.CircleAccessParent{CircleId: id})
	if err != nil {
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Circle{}, err
	}

	return circle, nil
}

// GetCircle gets a circle.
func (d *Domain) GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (circle model.Circle, err error) {
	if authAccount.UserId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
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

	circles, err := d.repo.ListCircles(ctx, authAccount, pageSize, pageOffset, filter, fieldMask)
	if err != nil {
		return nil, err
	}

	return circles, nil
}

// UpdateCircle updates a circle.
func (d *Domain) UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, updateMask []string) (dbCircle model.Circle, err error) {
	if authAccount.UserId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if circle.Id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.CircleId = circle.Id.CircleId

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		return model.Circle{}, err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	updated, err := d.repo.UpdateCircle(ctx, authAccount, circle, updateMask)
	if err != nil {
		return model.Circle{}, err
	}

	return updated, nil
}
