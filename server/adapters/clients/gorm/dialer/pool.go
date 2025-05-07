package dialer

import (
	"time"

	"gorm.io/gorm"
)

var _ gorm.Option = (*PoolConfigOption)(nil)

// NewPoolConfigOption -
func NewPoolConfigOption(maxIdleConns, maxOpenConns int, maxConnLifetime, maxConnIdleTime time.Duration) *PoolConfigOption {
	return &PoolConfigOption{
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		MaxConnLifetime: maxConnLifetime,
		MaxConnIdleTime: maxConnIdleTime,
	}
}

// PoolConfigOption -
type PoolConfigOption struct {
	MaxIdleConns    int
	MaxOpenConns    int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// AfterInitialize implements gorm.Option
func (c *PoolConfigOption) AfterInitialize(db *gorm.DB) error {
	d, err := db.DB()
	if err != nil {
		return err
	}

	d.SetMaxIdleConns(c.MaxIdleConns)
	d.SetMaxOpenConns(c.MaxOpenConns)
	d.SetConnMaxLifetime(c.MaxConnLifetime)
	d.SetConnMaxIdleTime(c.MaxConnIdleTime)
	return nil
}

// Apply implements gorm.Option
func (c *PoolConfigOption) Apply(*gorm.Config) error {
	return nil
}
