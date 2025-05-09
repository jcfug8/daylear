package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCircle creates a new circle.
func (d *Domain) CreateCircle(ctx context.Context, circle model.Circle) (model.Circle, error) {
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

	err = tx.BulkCreateCircleUsers(ctx, dbCircle.Id, []int64{dbCircle.Parent.UserId}, permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE)
	if err != nil {
		return model.Circle{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Circle{}, err
	}

	return dbCircle, nil
}
