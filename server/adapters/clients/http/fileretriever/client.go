package fileretriever

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"github.com/rs/zerolog"

	"go.uber.org/fx"
)

const maxDownloadSize = 2 * 1024 * 1024 // 2 MB

type Client struct {
	log zerolog.Logger
}

type NewClientParams struct {
	fx.In
}

func NewClient(params NewClientParams) (*Client, error) {
	return &Client{}, nil
}

func (c *Client) GetFileContents(ctx context.Context, location string) (io.ReadCloser, error) {
	resp, err := http.Get(location)
	if err != nil {
		return file.File{}, err
	}

	if resp.ContentLength > maxDownloadSize {
		return file.File{}, fileretriever.ErrInvalidArgument{
			Msg: fmt.Sprintf("the file at %s was too large", location),
		}
	}

	return resp.Body, nil
}
