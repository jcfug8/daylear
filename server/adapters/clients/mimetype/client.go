package mimetype

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jcfug8/daylear/server/core/file"
	"github.com/rs/zerolog"
)

type Client struct {
	log zerolog.Logger
}

func NewClient(log zerolog.Logger) *Client {
	return &Client{
		log: log,
	}
}

func (c *Client) Inspect(ctx context.Context, fileReader io.Reader) (file.File, error) {
	tempFile, err := os.CreateTemp("/tmp", "inspect-file-*.tmp")
	if err != nil {
		return file.File{}, err
	}
	defer func() {
		go func() {
			time.Sleep(30 * time.Second)
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()
	}()

	// Write the contents of the fileReader to the temporary file
	contentLength, err := io.Copy(tempFile, fileReader)
	if err != nil {
		tempFile.Close()
		return file.File{}, err
	}

	// Seek to the beginning of the temporary file
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		tempFile.Close()
		return file.File{}, err
	}

	// Detect the MIME type using the temporary file
	mtype, err := mimetype.DetectReader(tempFile)
	if err != nil {
		tempFile.Close()
		return file.File{}, err
	}

	// Seek to the beginning of the temporary file
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		tempFile.Close()
		return file.File{}, err
	}

	return file.File{
		ContentType:    mtype.String(),
		Extension:      mtype.Extension(),
		ReadSeekCloser: tempFile,
		ContentLength:  contentLength,
	}, nil
}
