package fileinspector

import (
	"context"
	"io"

	"github.com/jcfug8/daylear/server/core/file"
)

type Client interface {
	Inspect(ctx context.Context, fileReader io.Reader) (file.File, error)
}
