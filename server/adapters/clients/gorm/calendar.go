package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CalendarMap maps the core model fields to the database model fields for the Calendar model.
var CalendarMap = map[string]string{
	"permission": gmodel.CalendarAccessFields.PermissionLevel,
	"visibility": gmodel.CalendarFields.Visibility,
	"state":      gmodel.CalendarFields.State,
}

// CreateCalendar creates a new calendar
func (repo *Client) CreateCalendar(ctx context.Context, m cmodel.Calendar) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	gm, err := convert.CalendarFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar model")
		return cmodel.Calendar{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid calendar: %v", err)}
	}

	calendarFields := masks.RemovePaths(
		gmodel.CalendarFields.Mask(),
	)

	err = repo.db.
		Select(calendarFields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Create failed")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err = convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read calendar")
		return cmodel.Calendar{}, fmt.Errorf("unable to read calendar: %v", err)
	}

	return m, nil
}

// DeleteCalendar deletes a calendar
func (repo *Client) DeleteCalendar(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CalendarId) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	gm := gmodel.Calendar{CalendarId: id.CalendarId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.CalendarFields.Mask()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Delete failed")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read calendar")
		return cmodel.Calendar{}, fmt.Errorf("unable to read calendar: %v", err)
	}

	return m, nil
}

// GetCalendar retrieves a calendar
func (repo *Client) GetCalendar(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CalendarId) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	gm := gmodel.Calendar{}

	tx := repo.db.WithContext(ctx).
		Where("calendar.calendar_id = ?", id.CalendarId)

	if authAccount.CircleId != 0 {
		tx = tx.Select("calendar.*, ca.permission_level, ca.state, ca.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_circle_id = ?", authAccount.CircleId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_circle_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.CircleId)
	} else if authAccount.UserId != 0 {
		tx = tx.Select("calendar.*, ca.permission_level, ca.state, ca.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId)
	} else {
		tx = tx.Select("calendar.*, calendar_access.permission_level, calendar_access.state, calendar_access.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.First failed")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read calendar")
		return cmodel.Calendar{}, fmt.Errorf("unable to read calendar: %v", err)
	}

	return m, nil
}

// ListCalendars lists calendars
func (repo *Client) ListCalendars(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	dbCalendars := []gmodel.Calendar{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "calendar.calendar_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	if authAccount.CircleId != 0 {
		tx = tx.Select("calendar.*, ca.permission_level, ca.state, ca.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_circle_id = ?", authAccount.CircleId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_circle_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.CircleId)
	} else if authAccount.UserId != 0 {
		tx = tx.Select("calendar.*, ca.permission_level, ca.state, ca.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId)
	} else {
		tx = tx.Select("calendar.*, calendar_access.permission_level, calendar_access.state, calendar_access.calendar_access_id").
			Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	conversion, err := repo.calendarSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbCalendars).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Calendar, len(dbCalendars))
	for i, m := range dbCalendars {
		res[i], err = convert.CalendarToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("unable to read calendar")
			return nil, fmt.Errorf("unable to read calendar: %v", err)
		}
	}

	return res, nil
}

// UpdateCalendar updates a calendar
func (repo *Client) UpdateCalendar(ctx context.Context, authAccount cmodel.AuthAccount, calendar cmodel.Calendar, fields []string) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	gm, err := convert.CalendarFromCoreModel(calendar)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar model")
		return cmodel.Calendar{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid calendar: %v", err)}
	}

	mask := masks.Map(fields, gmodel.CalendarMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read calendar")
		return cmodel.Calendar{}, fmt.Errorf("unable to read calendar: %v", err)
	}

	return m, nil
}
