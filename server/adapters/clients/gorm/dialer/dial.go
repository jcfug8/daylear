package dialer

import (
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// DialParams -
type DialParams struct {
	fx.In

	Config  *Config
	Dialect Dialect
	Dialer  Dialer
	Opts    []gorm.Option `group:"gormOptions"`
}

// Dial -
func Dial(p DialParams) (*gorm.DB, error) {
	return p.Dialer.Dial(p.Dialect, p.Config, p.Opts...)
}
