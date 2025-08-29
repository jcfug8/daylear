package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type accessKeyDomain interface {
	CreateAccessKey(ctx context.Context, authAccount model.AuthAccount, accessKey model.AccessKey) (model.AccessKey, error)
	DeleteAccessKey(ctx context.Context, authAccount model.AuthAccount, id model.AccessKeyId) (model.AccessKey, error)
	GetAccessKey(ctx context.Context, authAccount model.AuthAccount, parent model.AccessKeyParent, id model.AccessKeyId, fields []string) (model.AccessKey, error)
	ListAccessKeys(ctx context.Context, authAccount model.AuthAccount, userId int64, pageSize int32, offset int64, filter string, fields []string) ([]model.AccessKey, error)
	UpdateAccessKey(ctx context.Context, authAccount model.AuthAccount, accessKey model.AccessKey, fields []string) (model.AccessKey, error)
}
