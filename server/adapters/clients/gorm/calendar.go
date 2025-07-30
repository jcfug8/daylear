package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// CreateCalendar creates a new calendar
func (c *Client) CreateCalendar(ctx context.Context, calendar cmodel.Calendar) (cmodel.Calendar, error) {
	gormCalendar := convert.CalendarToGorm(calendar)

	err := c.db.WithContext(ctx).Create(&gormCalendar).Error
	if err != nil {
		return cmodel.Calendar{}, err
	}

	return convert.CalendarFromGorm(gormCalendar), nil
}

// DeleteCalendar deletes a calendar
func (c *Client) DeleteCalendar(ctx context.Context, authAccount cmodel.AuthAccount, id int64) (cmodel.Calendar, error) {
	var gormCalendar gmodel.Calendar

	err := c.db.WithContext(ctx).Where("calendar_id = ? AND user_id = ? AND circle_id = ?",
		id, authAccount.UserId, authAccount.CircleId).First(&gormCalendar).Error
	if err != nil {
		return cmodel.Calendar{}, err
	}

	err = c.db.WithContext(ctx).Delete(&gormCalendar).Error
	if err != nil {
		return cmodel.Calendar{}, err
	}

	return convert.CalendarFromGorm(gormCalendar), nil
}

// GetCalendar retrieves a calendar
func (c *Client) GetCalendar(ctx context.Context, authAccount cmodel.AuthAccount, id int64) (cmodel.Calendar, error) {
	var gormCalendar gmodel.Calendar

	err := c.db.WithContext(ctx).Where("calendar_id = ? AND user_id = ? AND circle_id = ?",
		id, authAccount.UserId, authAccount.CircleId).First(&gormCalendar).Error
	if err != nil {
		return cmodel.Calendar{}, err
	}

	return convert.CalendarFromGorm(gormCalendar), nil
}

// ListCalendars lists calendars
func (c *Client) ListCalendars(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]cmodel.Calendar, error) {
	var gormCalendars []gmodel.Calendar

	query := c.db.WithContext(ctx).Where("user_id = ? AND circle_id = ?", authAccount.UserId, authAccount.CircleId)

	// Apply pagination
	query = query.Limit(int(pageSize)).Offset(int(offset))

	// Apply filtering if provided
	if filter != "" {
		// TODO: Implement filtering logic
	}

	err := query.Find(&gormCalendars).Error
	if err != nil {
		return nil, err
	}

	calendars := make([]cmodel.Calendar, len(gormCalendars))
	for i, gormCalendar := range gormCalendars {
		calendars[i] = convert.CalendarFromGorm(gormCalendar)
	}

	return calendars, nil
}

// UpdateCalendar updates a calendar
func (c *Client) UpdateCalendar(ctx context.Context, authAccount cmodel.AuthAccount, calendar cmodel.Calendar, updateMask []string) (cmodel.Calendar, error) {
	gormCalendar := convert.CalendarToGorm(calendar)

	// Apply field mask if provided
	if len(updateMask) > 0 {
		// TODO: Implement field mask logic
	}

	err := c.db.WithContext(ctx).Save(&gormCalendar).Error
	if err != nil {
		return cmodel.Calendar{}, err
	}

	return convert.CalendarFromGorm(gormCalendar), nil
}
