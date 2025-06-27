package domain

import (
	"context"
	// "fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircle creates a new circle.
func (d *Domain) CreateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle) (model.Circle, error) {
	// TODO: Implement CreateCircle
	// Implementation commented out for refactoring
	/*
		if circle.Parent.UserId == 0 {
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
		dbCircle.Parent = circle.Parent

		err = tx.BulkCreateCircleUsers(ctx, dbCircle.Id, []int64{dbCircle.Parent.UserId}, permPb.PermissionLevel_PERMISSION_LEVEL_WRITE)
		if err != nil {
			return model.Circle{}, err
		}

		err = tx.Commit()
		if err != nil {
			return model.Circle{}, err
		}

		return dbCircle, nil
	*/
	return model.Circle{}, domain.ErrInternal{Msg: "CreateCircle method not implemented"}
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
func (d *Domain) GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId, fieldMask []string) (model.Circle, error) {
	// TODO: Implement GetCircle
	// Implementation commented out for refactoring
	/*
		if id.CircleId == 0 {
			return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
		}

		if parent.UserId != 0 {
			permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, id.CircleId)
			if err != nil {
				return model.Circle{}, err
			}
			if permission < permPb.PermissionLevel_PERMISSION_LEVEL_READ {
				return model.Circle{}, domain.ErrPermissionDenied{Msg: "user does not have read permission"}
			}
		}

		circle := model.Circle{Id: id, Parent: parent}
		found, err := d.repo.GetCircle(ctx, circle, fieldMask)
		if err != nil {
			return model.Circle{}, err
		}
		found.Parent = parent
		return found, nil
	*/
	return model.Circle{}, domain.ErrInternal{Msg: "GetCircle method not implemented"}
}

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string, fieldMask []string) ([]model.Circle, error) {
	// TODO: Implement ListCircles
	// Implementation commented out for refactoring
	/*
		if parent.UserId != 0 {
			filter = fmt.Sprintf("%s user_id = %d", filter, parent.UserId)
		} else {
			filter = fmt.Sprintf("%s is_public = true", filter)
		}

		circles, err := d.repo.ListCircles(ctx, page, filter, fieldMask)
		if err != nil {
			return nil, err
		}
		for i := range circles {
			circles[i].Parent = parent
		}
		return circles, nil
	*/
	return nil, domain.ErrInternal{Msg: "ListCircles method not implemented"}
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
