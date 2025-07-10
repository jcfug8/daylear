package imagemagick

import (
	"context"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/jcfug8/daylear/server/ports/image"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateImage(ctx context.Context, imageReader io.Reader) (image.Image, error) {
	tmp, err := os.CreateTemp("", "magick-img-*.img")
	if err != nil {
		return nil, err
	}
	defer tmp.Close()
	_, err = io.Copy(tmp, imageReader)
	if err != nil {
		os.Remove(tmp.Name())
		return nil, err
	}
	// Detect format using file extension or magick identify
	cmd := exec.CommandContext(ctx, "magick", "identify", "-ping", "-format", "%m", tmp.Name())
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	format := strings.ToLower(strings.TrimSpace(string(out)))
	return &magickImage{path: tmp.Name(), Format: format}, nil
}
