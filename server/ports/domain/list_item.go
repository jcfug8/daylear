package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type listItemDomain interface {
	CreateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem) (model.ListItem, error)
	DeleteListItem(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, id model.ListItemId) (model.ListItem, error)
	GetListItem(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, id model.ListItemId, fields []string) (model.ListItem, error)
	ListListItems(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, pageSize int32, pageOffset int32, filter string, fields []string) ([]model.ListItem, error)
	UpdateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem, fields []string) (model.ListItem, error)
}
