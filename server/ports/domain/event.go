package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type eventDomain interface {
	CreateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event) (model.Event, error)
	DeleteEvent(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, id model.EventId) (model.Event, error)
	GetEvent(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, id model.EventId, fields []string) (model.Event, error)
	ListEvents(ctx context.Context, authAccount model.AuthAccount, parent model.EventParent, pageSize int32, offset int64, filter string, fields []string) ([]model.Event, error)
	UpdateEvent(ctx context.Context, authAccount model.AuthAccount, event model.Event, fields []string) (model.Event, error)
}
