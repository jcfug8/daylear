package imagemagick

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/ports/image"
)

var _ image.Image = (*magickImage)(nil)

type magickImage struct {
	path   string
	Format string
	closer io.Closer
}

func (m *magickImage) Resize(ctx context.Context, width int, height int) error {
	cmd := exec.CommandContext(ctx, "magick", m.path, "-resize", fmt.Sprintf("%dx%d>", width, height), m.path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (m *magickImage) GetFile() (file.File, error) {
	f, err := os.Open(m.path)
	if err != nil {
		return file.File{}, err
	}
	stat, err := f.Stat()
	if err != nil {
		f.Close()
		return file.File{}, err
	}
	m.closer = f
	return file.File{
		ContentType:    "image/" + m.Format,
		Extension:      "." + m.Format,
		ReadSeekCloser: f,
		ContentLength:  stat.Size(),
	}, nil
}

func (m *magickImage) Remove(ctx context.Context) error {
	if m.closer != nil {
		m.closer.Close()
	}
	return os.Remove(m.path)
}

func (m *magickImage) GetDimensions(ctx context.Context) (width int, height int, err error) {
	cmd := exec.CommandContext(ctx, "magick", "identify", "-ping", "-format", "%w %h", m.path)
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Split(string(out), " ")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unable to get dimensions: %s", string(out))
	}
	width, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	height, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func (m *magickImage) Convert(ctx context.Context, format string) error {
	if m.Format == format {
		return nil
	}

	outPath := m.path + ".converted." + format
	cmd := exec.CommandContext(ctx, "magick", m.path, outPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	m.Format = "." + format
	m.path = outPath
	return nil
}
