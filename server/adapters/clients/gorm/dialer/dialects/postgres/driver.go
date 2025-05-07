package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Driver -
type Driver interface {
	Initialize(*gorm.DB) error
	Name() string
}

type pgDriver struct {
	driver Driver
	*postgres.Dialector
}

// Initialize -
func (pg *pgDriver) Initialize(db *gorm.DB) error {
	if err := pg.Dialector.Initialize(db); err != nil {
		return err
	}

	return pg.driver.Initialize(db)
}
