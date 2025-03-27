package file

import "io"

type File struct {
	ContentType string
	Extension   string
	io.ReadSeekCloser
	ContentLength int64
}
