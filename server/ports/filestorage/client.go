package filestorage

import (
	"context"

	"github.com/jcfug8/daylear/server/core/file"
)

type Client interface {
	UploadPublicFile(ctx context.Context, path string, file file.File) (string, error)
	DeleteFile(ctx context.Context, path string) error
}
