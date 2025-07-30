package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
)

// CreateCalendarAccess creates calendar access
func (c *Client) CreateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess) (cmodel.CalendarAccess, error) {
	gormAccess := convert.CalendarAccessToGorm(access)

	err := c.db.WithContext(ctx).Create(&gormAccess).Error
	if err != nil {
		return cmodel.CalendarAccess{}, err
	}

	return convert.CalendarAccessFromGorm(gormAccess), nil
}

// DeleteCalendarAccess deletes calendar access
func (c *Client) DeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent, id cmodel.CalendarAccessId) error {
	var gormAccess gmodel.CalendarAccess

	err := c.db.WithContext(ctx).Where("calendar_access_id = ? AND calendar_id = ?",
		id.CalendarAccessId, parent.CalendarId).First(&gormAccess).Error
	if err != nil {
		return err
	}

	return c.db.WithContext(ctx).Delete(&gormAccess).Error
}

// BulkDeleteCalendarAccesses bulk deletes calendar accesses
func (c *Client) BulkDeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)

	db := c.db.WithContext(ctx)

	res := db.Where("calendar_id = ?", parent.CalendarId).Delete(&gmodel.CalendarAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("db.Delete failed")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no rows affected (not found)")
		return repository.ErrNotFound{}
	}

	return nil
}

// GetCalendarAccess retrieves calendar access
func (c *Client) GetCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent, id cmodel.CalendarAccessId) (cmodel.CalendarAccess, error) {
	var gormAccess gmodel.CalendarAccess

	err := c.db.WithContext(ctx).Where("calendar_access_id = ? AND calendar_id = ?",
		id.CalendarAccessId, parent.CalendarId).First(&gormAccess).Error
	if err != nil {
		return cmodel.CalendarAccess{}, err
	}

	return convert.CalendarAccessFromGorm(gormAccess), nil
}

// ListCalendarAccesses lists calendar accesses
func (c *Client) ListCalendarAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.CalendarAccessParent, pageSize int32, pageOffset int64, filter string) ([]cmodel.CalendarAccess, error) {
	var gormAccesses []gmodel.CalendarAccess

	query := c.db.WithContext(ctx).Where("calendar_id = ?", parent.CalendarId)

	// Apply pagination
	query = query.Limit(int(pageSize)).Offset(int(pageOffset))

	// Apply filtering if provided
	if filter != "" {
		// TODO: Implement filtering logic
	}

	err := query.Find(&gormAccesses).Error
	if err != nil {
		return nil, err
	}

	accesses := make([]cmodel.CalendarAccess, len(gormAccesses))
	for i, gormAccess := range gormAccesses {
		accesses[i] = convert.CalendarAccessFromGorm(gormAccess)
	}

	return accesses, nil
}

// UpdateCalendarAccess updates calendar access
func (c *Client) UpdateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess, updateMask []string) (cmodel.CalendarAccess, error) {
	gormAccess := convert.CalendarAccessToGorm(access)

	// Apply field mask if provided
	if len(updateMask) > 0 {
		// TODO: Implement field mask logic
	}

	err := c.db.WithContext(ctx).Save(&gormAccess).Error
	if err != nil {
		return cmodel.CalendarAccess{}, err
	}

	return convert.CalendarAccessFromGorm(gormAccess), nil
}
