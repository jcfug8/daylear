package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type eventClient interface {
	CreateEvent(ctx context.Context, event model.Event, fields []string) (model.Event, error)
	CreateEventClones(ctx context.Context, event []model.Event) ([]model.Event, error)
	DeleteEvent(ctx context.Context, id model.EventId) (model.Event, error)
	BulkDeleteEvents(ctx context.Context, ids []model.EventId) error
	DeleteChildEvents(ctx context.Context, id model.EventId) error
	GetEvent(ctx context.Context, authAccount model.AuthAccount, id model.EventId, fields []string) (model.Event, error)
	ListEvents(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, pageSize int32, offset int64, filter string, fields []string) ([]model.Event, error)
	UpdateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event, fields []string) (model.Event, error)
}
