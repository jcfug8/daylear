package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type listDomain interface {
	CreateList(ctx context.Context, authAccount model.AuthAccount, list model.List) (model.List, error)
	DeleteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) error
	GetList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId, fields []string) (model.List, error)
	ListLists(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, pageSize int32, pageOffset int32, filter string, fields []string) ([]model.List, error)
	UpdateList(ctx context.Context, authAccount model.AuthAccount, list model.List, fields []string) (model.List, error)
	FavoriteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) error
	UnfavoriteList(ctx context.Context, authAccount model.AuthAccount, parent model.ListParent, id model.ListId) error

	CreateListAccess(ctx context.Context, authAccount model.AuthAccount, access model.ListAccess) (model.ListAccess, error)
	DeleteListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId) error
	GetListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId, fields []string) (model.ListAccess, error)
	ListListAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.ListAccess, error)
	UpdateListAccess(ctx context.Context, authAccount model.AuthAccount, access model.ListAccess, fields []string) (model.ListAccess, error)
	AcceptListAccess(ctx context.Context, authAccount model.AuthAccount, parent model.ListAccessParent, id model.ListAccessId) (model.ListAccess, error)
}
