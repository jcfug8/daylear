package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// UpdateCircle updates a circle.
func (d *Domain) UpdateCircle(ctx context.Context, circle model.Circle, updateMask []string) (model.Circle, error) {
	if circle.Parent.UserId == 0 || circle.Id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent and id required"}
	}

	permission, err := d.repo.GetCircleUserPermission(ctx, circle.Parent.UserId, circle.Id.CircleId)
	if err != nil {
		return model.Circle{}, err
	}
	if permission < permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
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
}
