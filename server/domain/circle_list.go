package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
)

// ListCircles lists circles for a parent.
func (d *Domain) ListCircles(ctx context.Context, page *model.PageToken[model.Circle], parent model.CircleParent, filter string, fieldMask []string) ([]model.Circle, error) {
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
}
