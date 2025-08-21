package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the circle in the database.
type circleClient interface {
	CreateCircle(ctx context.Context, circle model.Circle, fields []string) (model.Circle, error)
	DeleteCircle(ctx context.Context, id model.CircleId) (model.Circle, error)
	GetCircle(ctx context.Context, authAccount model.AuthAccount, id model.CircleId, fields []string) (model.Circle, error)
	ListCircles(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]model.Circle, error)
	UpdateCircle(ctx context.Context, authAccount model.AuthAccount, circle model.Circle, fields []string) (model.Circle, error)

	CreateCircleFavorite(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) error
	DeleteCircleFavorite(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) error

	FindStandardUserCircleAccess(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.CircleAccess, error)
	FindDelegatedUserCircleAccess(ctx context.Context, authAccount model.AuthAccount, id model.CircleId) (model.CircleAccess, model.UserAccess, error)

	CreateCircleAccess(ctx context.Context, access model.CircleAccess, fields []string) (model.CircleAccess, error)
	DeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) error
	BulkDeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent) error
	GetCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId, fields []string) (model.CircleAccess, error)
	ListCircleAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.CircleAccess, error)
	UpdateCircleAccess(ctx context.Context, access model.CircleAccess, fields []string) (model.CircleAccess, error)
	CircleHandleExists(ctx context.Context, handle string) (bool, error)
}
