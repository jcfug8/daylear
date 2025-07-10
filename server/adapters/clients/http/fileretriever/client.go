package fileretriever

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/jcfug8/daylear/server/core/logutil"
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

func NewClient(log zerolog.Logger) *Client {
	return &Client{log: log}
}

func (c *Client) GetFileContents(ctx context.Context, location string) (io.ReadCloser, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	resp, err := http.Get(location)
	if err != nil {
		log.Error().Err(err).Str("location", location).Msg("failed to fetch file")
		return nil, err
	}

	if resp.ContentLength > maxDownloadSize {
		log.Warn().Str("location", location).Int64("content_length", resp.ContentLength).Msg("file too large")
		return nil, fileretriever.ErrInvalidArgument{
			Msg: fmt.Sprintf("the file at %s was too large", location),
		}
	}

	return resp.Body, nil
}
