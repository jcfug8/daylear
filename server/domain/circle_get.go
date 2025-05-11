package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// GetCircle gets a circle.
func (d *Domain) GetCircle(ctx context.Context, parent model.CircleParent, id model.CircleId, fieldMask []string) (model.Circle, error) {
	if id.CircleId == 0 {
		return model.Circle{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if parent.UserId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, id.CircleId)
		if err != nil {
			return model.Circle{}, err
		}
		if permission < permPb.PermissionLevel_RESOURCE_PERMISSION_READ {
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
}
