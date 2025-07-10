package image

import (
	"context"
	"io"

	"github.com/jcfug8/daylear/server/core/file"
)

type Client interface {
	CreateImage(ctx context.Context, imageReader io.Reader) (Image, error)
}

type Image interface {
	Convert(ctx context.Context, format string) error
	Resize(ctx context.Context, width int, height int) (err error)
	GetDimensions(ctx context.Context) (width int, height int, err error)
	Remove(ctx context.Context) error
	GetFile() (file.File, error)
}
