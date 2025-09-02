package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// Client defines how to interact with the list in the database.
type listClient interface {
	CreateList(ctx context.Context, authAccount model.AuthAccount, list model.List) (model.List, error)
	DeleteList(ctx context.Context, authAccount model.AuthAccount, id model.ListId) (model.List, error)
	GetList(ctx context.Context, authAccount model.AuthAccount, id model.ListId, fields []string) (model.List, error)
	ListLists(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int32, filter string, fields []string) ([]model.List, error)
	UpdateList(ctx context.Context, authAccount model.AuthAccount, list model.List, fields []string) (model.List, error)

	FindStandardUserListAccess(ctx context.Context, authAccount model.AuthAccount, id model.ListId) (model.ListAccess, error)
	FindDelegatedCircleListAccess(ctx context.Context, authAccount model.AuthAccount, id model.ListId) (model.ListAccess, model.CircleAccess, error)
	FindDelegatedUserListAccess(ctx context.Context, authAccount model.AuthAccount, id model.ListId) (model.ListAccess, model.UserAccess, error)

	CreateListAccess(ctx context.Context, access model.ListAccess, fields []string) (model.ListAccess, error)
	DeleteListAccess(ctx context.Context, parent model.ListAccessParent, id model.ListAccessId) error
	BulkDeleteListAccess(ctx context.Context, parent model.ListAccessParent) error
	GetListAccess(ctx context.Context, parent model.ListAccessParent, id model.ListAccessId, fields []string) (model.ListAccess, error)
	ListListAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.ListAccess, error)
	UpdateListAccess(ctx context.Context, access model.ListAccess, fields []string) (model.ListAccess, error)

	// Favoriting methods
	CreateListFavorite(ctx context.Context, authAccount model.AuthAccount, id model.ListId) error
	DeleteListFavorite(ctx context.Context, authAccount model.AuthAccount, id model.ListId) error
	BulkDeleteListFavorites(ctx context.Context, id model.ListId) error
}
