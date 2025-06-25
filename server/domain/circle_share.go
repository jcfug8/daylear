package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// ShareCircle shares a circle with multiple users.
func (d *Domain) ShareCircle(ctx context.Context, parent model.CircleParent, parents []model.CircleParent, id model.CircleId, permission permPb.PermissionLevel) error {
	if parent.UserId == 0 || id.CircleId == 0 {
		return domain.ErrInvalidArgument{Msg: "parent and id required"}
	}
	if len(parents) == 0 {
		return domain.ErrInvalidArgument{Msg: "no recipients provided"}
	}

	permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, id.CircleId)
	if err != nil {
		return err
	}
	if permission != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	userIds := make([]int64, 0, len(parents))
	for _, p := range parents {
		if p.UserId != 0 {
			userIds = append(userIds, p.UserId)
		}
	}
	if len(userIds) == 0 {
		return domain.ErrInvalidArgument{Msg: "no valid recipient user ids"}
	}

	err = tx.BulkCreateCircleUsers(ctx, id, userIds, permission)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
