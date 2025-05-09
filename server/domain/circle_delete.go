package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// DeleteCircle deletes a circle.
func (d *Domain) DeleteCircle(ctx context.Context, parent model.CircleParent, id model.CircleId) (model.Circle, error) {
	if parent.UserId == 0 || id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "parent and id required"}
	}

	permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, id.CircleId)
	if err != nil {
		return model.Circle{}, err
	}
	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
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
}
