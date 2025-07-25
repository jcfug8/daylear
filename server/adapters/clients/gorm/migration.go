package gorm

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/migrations"
	model "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"gorm.io/gorm"
)

var manualMigrations = []migrations.Migration{
	{
		Key:         "1",
		Description: "Create circle.handle field and backfill with unique slugified titles",
		Run: func(ctx context.Context, tx *gorm.DB) error {
			// 1. Add the column if it doesn't exist (raw SQL)
			if !tx.Migrator().HasColumn(&model.Circle{}, "handle") {
				if err := tx.Exec(`ALTER TABLE circle ADD COLUMN handle VARCHAR(64)`).Error; err != nil {
					return fmt.Errorf("failed to add handle column: %w", err)
				}
			}

			// 2. Get all circles (id, title)
			type circleRow struct {
				CircleId int64
				Title    string
			}
			var circles []circleRow
			if err := tx.Model(&model.Circle{}).Select("circle_id, title").Find(&circles).Error; err != nil {
				return fmt.Errorf("failed to fetch circles: %w", err)
			}

			// 3. Build a set of existing handles to ensure uniqueness
			existingHandles := map[string]struct{}{}
			var usedHandles []string
			tx.Model(&model.Circle{}).Select("handle").Find(&usedHandles)
			for _, h := range usedHandles {
				h = strings.ToLower(h)
				existingHandles[h] = struct{}{}
			}

			// 4. For each circle, generate a unique handle and update
			for _, c := range circles {
				base := slugify(c.Title)
				if base == "" {
					base = fmt.Sprintf("circle-%d", c.CircleId)
				}
				handle := base
				n := 1
				for {
					_, exists := existingHandles[handle]
					if !exists {
						break
					}
					handle = fmt.Sprintf("%s-%d", base, n)
					n++
				}
				existingHandles[handle] = struct{}{}
				if err := tx.Model(&model.Circle{}).Where("circle_id = ?", c.CircleId).Update("handle", handle).Error; err != nil {
					return fmt.Errorf("failed to update handle for circle %d: %w", c.CircleId, err)
				}
			}
			return nil
		},
	},
	{
		Key:         "2",
		Description: "Drop visibility column from user",
		Run: func(ctx context.Context, tx *gorm.DB) error {
			if err := tx.Exec(`ALTER TABLE daylear_user DROP COLUMN visibility`).Error; err != nil {
				return fmt.Errorf("failed to add visibility column: %w", err)
			}
			return nil
		},
	},
}

// slugify is a helper to create a URL-friendly string
func slugify(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`(^-|-$)`).ReplaceAllString(s, "")
	return s
}
