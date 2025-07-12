package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the circle in the database.
type circleClient interface {
	CreateCircle(ctx context.Context, circle model.Circle) (model.Circle, error)
	DeleteCircle(ctx context.Context, id model.CircleId) (model.Circle, error)
	GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.Circle, error)
	ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]model.Circle, error)
	UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, updateMask []string) (model.Circle, error)

	CreateCircleAccess(ctx context.Context, access model.CircleAccess) (model.CircleAccess, error)
	DeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) error
	BulkDeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent) error
	GetCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error)
	ListCircleAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.CircleAccess, error)
	UpdateCircleAccess(ctx context.Context, access model.CircleAccess, updateMask []string) (model.CircleAccess, error)
	CircleHandleExists(ctx context.Context, handle string) (bool, error)
}
