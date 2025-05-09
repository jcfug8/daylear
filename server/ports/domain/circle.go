package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

type circleDomain interface {
	CreateCircle(ctx context.Context, recipe model.Circle) (model.Circle, error)
	DeleteCircle(ctx context.Context, parent model.CircleParent, id model.CircleId) (model.Circle, error)
	GetCircle(ctx context.Context, parent model.CircleParent, id model.CircleId, fieldMask []string) (model.Circle, error)
	ListCircles(ctx context.Context, page *model.PageToken[model.Circle], parent model.CircleParent, filter string, fieldMask []string) ([]model.Circle, error)
	UpdateCircle(ctx context.Context, recipe model.Circle, updateMask []string) (model.Circle, error)
	ShareCircle(ctx context.Context, parent model.CircleParent, parents []model.CircleParent, id model.CircleId, permission permPb.PermissionLevel) error
}
