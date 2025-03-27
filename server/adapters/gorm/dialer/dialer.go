package dialer

import (
	"errors"
	"fmt"
	"io"
	"time"

	"gorm.io/gorm"
)

var _ Dialer = (*DefaultDialer)(nil)

// Dialer -
type Dialer interface {
	Dial(Dialect, *Config, ...gorm.Option) (*gorm.DB, error)
}

// Dialect -
type Dialect interface {
	Dialector() gorm.Dialector
	Options() []gorm.Option
}

// NewDialer -
func NewDialer() *DefaultDialer {
	return &DefaultDialer{}
}

// DefaultDialer -
type DefaultDialer struct{}

// Dial -
func (d *DefaultDialer) Dial(dialect Dialect, config *Config, opts ...gorm.Option) (db *gorm.DB, err error) {
	options := append(dialect.Options(), opts...)
	options = append(options, config.Config)

	config.log.Info().Msgf("dialing gorm with %d options", len(options))

	flag := 0
	for flag < 10 {
		config.log.Info().Msg("waiting for DB connection...")
		db, err = gorm.Open(dialect.Dialector(), options...)
		if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
			config.log.Warn().Err(err).Msgf("couldn't connect to the db")
		}

		if err == nil {
			break
		}

		time.Sleep(time.Second * 1)
		flag++
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	} else if db == nil {
		return nil, fmt.Errorf("failed to create db instance")
	}

	config.log.Info().Msgf("DB connection established with %d options", len(options))

	plugins := make([]string, len(config.plugins))
	for i, plugin := range config.plugins {
		plugins[i] = plugin.Name()
		config.log.Debug().Msgf("using plugin: %s", plugin.Name())

		if err = db.Use(plugin); err != nil {
			return nil, fmt.Errorf("unable to use plugin: %w", err)
		}
	}

	config.log.Info().Strs("plugins", plugins).Msgf("DB connection configured with %d plugins", len(config.plugins))

	return db, nil
}
