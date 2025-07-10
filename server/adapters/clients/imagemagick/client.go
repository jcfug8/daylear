package imagemagick

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/ports/image"
	"github.com/rs/zerolog"
)

type Client struct {
	log zerolog.Logger
}

func NewClient(log zerolog.Logger) (*Client, error) {
	// call magick identify to check if it is installed
	cmd := exec.Command("magick", "identify", "--version")
	out, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msgf("magick identify is not installed: %s: %s", err, string(out))
		return nil, fmt.Errorf("magick identify is not installed: %s: %s", err, string(out))
	}
	return &Client{log: log}, nil
}

func (c *Client) CreateImage(ctx context.Context, imageReader io.Reader) (image.Image, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	tmp, err := os.CreateTemp("", "magick-img-*.img")
	if err != nil {
		log.Error().Err(err).Msg("failed to create temp file")
		return nil, err
	}
	defer tmp.Close()
	_, err = io.Copy(tmp, imageReader)
	if err != nil {
		log.Error().Err(err).Msg("failed to copy image data")
		os.Remove(tmp.Name())
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "magick", "identify", "-ping", "-format", "%m", tmp.Name())
	out, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msgf("failed to identify image format: %s: %s", err, string(out))
		return nil, err
	}
	format := strings.ToLower(strings.TrimSpace(string(out)))
	return &magickImage{path: tmp.Name(), Format: format}, nil
}
