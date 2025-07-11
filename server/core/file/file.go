package file

import (
	"bytes"
	"io"
)

type File struct {
	ContentType string
	Extension   string
	io.ReadSeekCloser
	ContentLength int64
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil
}

func NewReadSeekCloser(b []byte) io.ReadSeekCloser {
	return &readSeekCloser{
		Reader: bytes.NewReader(b),
	}
}
