package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type circleDomain interface {
	CreateCircle(ctx context.Context, authAccount model.AuthAccount, recipe model.Circle) (model.Circle, error)
	DeleteCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.Circle, error)
	GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.Circle, error)
	ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]model.Circle, error)
	UpdateCircle(ctx context.Context, authAccount model.AuthAccount, recipe model.Circle, updateMask []string) (model.Circle, error)

	CreateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (model.CircleAccess, error)
	DeleteCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) error
	GetCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error)
	ListCircleAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.CircleAccess, error)
	UpdateCircleAccess(ctx context.Context, authAccount model.AuthAccount, access model.CircleAccess) (model.CircleAccess, error)
	AcceptCircleAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error)
}
