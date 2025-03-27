package plugins

import "gorm.io/gorm"

type PGCrypto struct{}

func (p *PGCrypto) Name() string {
	return "PGCrypto"
}

func (p *PGCrypto) Initialize(db *gorm.DB) error {
	tx := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto;`)
	return tx.Error
}
