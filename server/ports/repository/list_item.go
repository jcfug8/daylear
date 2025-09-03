package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with list items in the database.
type listItemClient interface {
	// Basic CRUD operations
	CreateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem) (model.ListItem, error)
	DeleteListItem(ctx context.Context, authAccount model.AuthAccount, id model.ListItemId) (model.ListItem, error)
	GetListItem(ctx context.Context, authAccount model.AuthAccount, id model.ListItemId, fields []string) (model.ListItem, error)
	ListListItems(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, pageSize int32, pageOffset int32, filter string, fields []string) ([]model.ListItem, error)
	UpdateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem, fields []string) (model.ListItem, error)

	// Bulk operations
	BulkDeleteListItems(ctx context.Context, parent model.ListItemParent) error
}
