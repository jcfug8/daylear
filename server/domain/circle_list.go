package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, page *model.PageToken[model.Circle], parent model.CircleParent, filter string, fieldMask []string) ([]model.Circle, error) {
	if parent.UserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	circles, err := d.repo.ListCircles(ctx, page, filter, fieldMask)
	if err != nil {
		return nil, err
	}
	for i := range circles {
		circles[i].Parent = parent
	}
	return circles, nil
}
