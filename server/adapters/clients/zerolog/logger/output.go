package logger

import (
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// OutputParams -
type OutputParams struct {
	fx.In

	Offset *Offset
	Writer io.Writer `name:"loggerWriter" optional:"true"`
}

// NewOutput -
func NewOutput(p OutputParams) *Output {
	var w io.Writer = os.Stderr
	if p.Writer != nil {
		w = p.Writer
	}

	return &Output{
		w:      w,
		offset: p.Offset,
	}
}

// Output -
type Output struct {
	m      sync.Mutex
	w      io.Writer
	offset *Offset
}

// GetOffset returns the offset for the output.
func (o *Output) GetOffset() *Offset {
	return o.offset
}

// SetOutput sets the output destination for the logger.
func (o *Output) SetOutput(w io.Writer) {
	o.m.Lock()
	defer o.m.Unlock()

	o.w = w
}

// Write implements io.Write
func (o *Output) Write(p []byte) (n int, err error) {
	o.m.Lock()
	defer o.m.Unlock()

	return o.w.Write(p)
}

// WriteLevel implements zerolog.LevelWriter
func (o *Output) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < o.offset.GetLevel() {
		return len(p), nil
	}

	return o.Write(p)
}
