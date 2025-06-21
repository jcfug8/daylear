package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// Client defines how to interact with the circle in the database.
type circleClient interface {
	CreateCircle(context.Context, model.Circle) (model.Circle, error)
	DeleteCircle(context.Context, model.Circle) (model.Circle, error)
	GetCircle(context.Context, model.Circle, []string) (model.Circle, error)
	UpdateCircle(context.Context, model.Circle, []string) (model.Circle, error)
	ListCircles(context.Context, *model.PageToken[model.Circle], string, []string) ([]model.Circle, error)

	BulkCreateCircleUsers(context.Context, model.CircleId, []int64, permPb.PermissionLevel) error
	BulkDeleteCircleUsers(context.Context, string) error
	GetCircleUserPermission(ctx context.Context, userId int64, circleId int64) (permPb.PermissionLevel, error)
}
