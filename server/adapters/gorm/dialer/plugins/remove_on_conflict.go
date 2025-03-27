package plugins

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RemoveOnConflict struct{}

func (*RemoveOnConflict) Name() string {
	return "remove_on_conflict"
}

func (*RemoveOnConflict) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("gorm:create").Register("remove_on_conflict",
		func(tx *gorm.DB) {
			delete(tx.Statement.Clauses, new(clause.OnConflict).Name())
		})
	return nil
}
