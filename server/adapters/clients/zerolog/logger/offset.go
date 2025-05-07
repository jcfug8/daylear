package logger

import (
	"os"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

const (
	minLevel = int(zerolog.TraceLevel)
	maxLevel = int(zerolog.Disabled)
)

// NewOffset creates a new offset for TraceLevel.
func NewOffset() *Offset {
	const base = minLevel
	return &Offset{
		baseLevel: base,
		offsetMin: minLevel - base,
		offsetMax: maxLevel - base,
	}
}

// Offset -
type Offset struct {
	m sync.RWMutex
	// offset needs to be in range [o.offsetMin, o.offsetMax]
	offset int
	// Needs to be in range [minLevel, maxLevel]
	baseLevel int
	// These are set based off of baseLevel
	offsetMin int
	offsetMax int
}

// Add adds the provided value to the offset. It returns the new offset and a
// boolean indicating whether the offset was actually changed.
func (o *Offset) Add(value int) (int, bool) {
	o.m.Lock()
	defer o.m.Unlock()

	prev := o.offset
	next := limit(prev+value, o.offsetMin, o.offsetMax)
	o.offset = next

	return next, next != prev
}

// Get returns the offset amount.
func (o *Offset) Get() int {
	o.m.RLock()
	defer o.m.RUnlock()
	return o.offset
}

// GetLevel returns the offset level.
func (o *Offset) GetLevel() zerolog.Level {
	o.m.RLock()
	defer o.m.RUnlock()
	return zerolog.Level(o.baseLevel + o.offset)
}

// Set sets the offset value to the one provided.
func (o *Offset) Set(value int) bool {
	o.m.Lock()
	defer o.m.Unlock()

	prev := o.offset
	next := limit(value, o.offsetMin, o.offsetMax)
	o.offset = next

	return next != prev
}

// Level shifts the offset to the provided level. It return a boolean if the
// offset was actually changed.
func (o *Offset) Level(level zerolog.Level) bool {
	o.m.Lock()
	defer o.m.Unlock()

	next := limit(int(level), minLevel, maxLevel)
	if next == o.baseLevel {
		return false
	}

	o.baseLevel = next
	o.offset = 0

	o.offsetMax = maxLevel - o.baseLevel
	o.offsetMin = minLevel - o.baseLevel

	return true
}

func limit(value, min, max int) int {
	switch {
	case value > max:
		return max
	case value < min:
		return min
	default:
		return value
	}
}

// Signal increases the offset on SIGUSR1 and decreases it on SIGUSR2. It
// expects os signals USR1 and USR2 to already be relayed to the provided
// channel using os/signal.Notify.
func (o *Offset) Signal(signals <-chan os.Signal) {
	log := zerolog.New(os.Stderr).
		Level(zerolog.InfoLevel).
		With().Timestamp().Logger()

	var name string
	var offset int

	for signal := range signals {
		switch signal {
		case syscall.SIGUSR1:
			offset, _ = o.Add(1)
			name = "USR1"
		case syscall.SIGUSR2:
			offset, _ = o.Add(-1)
			name = "USR2"
		}

		_, _ = name, offset

		log.Info().
			Str("signal", name).
			Msgf("received %s: offset=%d (%s)", name, offset, o.GetLevel())
	}
}
