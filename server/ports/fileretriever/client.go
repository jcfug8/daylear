package fileretriever

import (
	"context"
	"io"
)

type Client interface {
	GetFileContents(ctx context.Context, location string) (fileContents io.ReadCloser, err error)
}
