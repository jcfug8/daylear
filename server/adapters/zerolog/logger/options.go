package logger

import "github.com/rs/zerolog"

// Option -
type Option interface {
	apply(*Options)
}

// NewOptions -
func NewOptions(opts ...Option) *Options {
	options := &Options{}

	for _, opt := range opts {
		opt.apply(options)
	}

	return options
}

// Options -
type Options struct {
	level *zerolog.Level
}

// ----------------------------------------------------------------------------

func newFuncOption(f func(*Options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

type funcOption struct {
	f func(*Options)
}

func (o *funcOption) apply(options *Options) {
	o.f(options)
}

// ----------------------------------------------------------------------------

// WithLevel -
func WithLevel(level zerolog.Level) Option {
	return newFuncOption(func(options *Options) {
		options.level = &level
	})
}
