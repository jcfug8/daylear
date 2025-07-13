package migrations

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Migration represents a single manual migration.
type Migration struct {
	Key         string // unique key for the migration
	Description string
	Run         func(ctx context.Context, tx *gorm.DB) error
}

// RunMigrations runs the provided migrations, ensuring each is only run once.
func RunMigrations(ctx context.Context, db *gorm.DB, migrations []Migration) error {
	if err := ensureMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	for _, m := range migrations {
		run, err := hasMigrationRun(ctx, db, m.Key)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", m.Key, err)
		}
		if run {
			continue
		}

		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := m.Run(ctx, tx); err != nil {
				return fmt.Errorf("migration %s failed: %w", m.Key, err)
			}
			return recordMigration(ctx, tx, m.Key, m.Description)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func ensureMigrationsTable(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Exec(`
		CREATE TABLE IF NOT EXISTS manual_migrations (
			id SERIAL PRIMARY KEY,
			migration_key TEXT UNIQUE NOT NULL,
			description TEXT,
			applied_at TIMESTAMP NOT NULL
		)
	`).Error
}

func hasMigrationRun(ctx context.Context, db *gorm.DB, key string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Raw(
		"SELECT COUNT(1) FROM manual_migrations WHERE migration_key = ?", key,
	).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func recordMigration(ctx context.Context, db *gorm.DB, key, desc string) error {
	return db.WithContext(ctx).Exec(
		"INSERT INTO manual_migrations (migration_key, description, applied_at) VALUES (?, ?, ?)",
		key, desc, time.Now(),
	).Error
}
