package logger

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
)

// NewFlags -
func NewFlags(name string) *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.level, name, "info", "log level")
	return flags
}

// Flags -
type Flags struct {
	level string
}

// NewFlagParser -
func NewFlagParser(flags *Flags) *FlagParser {
	if !flag.Parsed() {
		flag.Parse()
	}

	return &FlagParser{flags: flags}
}

// FlagParser -
type FlagParser struct {
	flags *Flags
}

// Level -
func (p *FlagParser) Level() zerolog.Level {
	lvl := os.Getenv("MATRIX_LOG_LEVEL")

	if p != nil {
		lvl = p.flags.level
	}

	if lvl == "" {
		return zerolog.InfoLevel
	}

	level, err := zerolog.ParseLevel(lvl)
	if err != nil {
		return zerolog.InfoLevel
	}

	return level
}
