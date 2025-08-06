package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateCalendar creates a new calendar
func (repo *Client) CreateCalendar(ctx context.Context, m cmodel.Calendar, fields []string) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.CalendarFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar when creating calendar row")
		return cmodel.Calendar{}, repository.ErrInvalidArgument{Msg: "invalid calendar"}
	}

	err = repo.db.
		Select(gmodel.CalendarAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create calendar row")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err = convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar row when creating calendar")
		return cmodel.Calendar{}, repository.ErrInternal{Msg: "invalid calendar row when creating calendar"}
	}

	return m, nil
}

// DeleteCalendar deletes a calendar
func (repo *Client) DeleteCalendar(ctx context.Context, id cmodel.CalendarId) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", id.CalendarId).
		Logger()

	gm := gmodel.Calendar{CalendarId: id.CalendarId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.CalendarAccessFieldMasker.GetAll()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete calendar row")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalide calendar row when deleting calendar")
		return cmodel.Calendar{}, repository.ErrInternal{Msg: "invalid calendar row when deleting calendar"}
	}

	return m, nil
}

// GetCalendar retrieves a calendar
func (repo *Client) GetCalendar(ctx context.Context, id cmodel.CalendarId, fields []string) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", id.CalendarId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.Calendar{}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.CalendarFieldMasker.Convert(fields)).
		Where("calendar.calendar_id = ?", id.CalendarId)

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar row")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar row when getting calendar")
		return cmodel.Calendar{}, repository.ErrInternal{Msg: "invalid calendar row when getting calendar"}
	}

	return m, nil
}

// ListCalendars lists calendars
func (repo *Client) ListCalendars(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("filter", filter).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("offset", offset).Logger()
	dbCalendars := []gmodel.Calendar{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "calendar.calendar_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.CalendarFieldMasker.Convert(fields)).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	if authAccount.CircleId != 0 {
		tx = tx.Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_circle_id = ?", authAccount.CircleId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar_access.recipient_circle_id = ? AND (calendar.visibility_level = ? OR ca.state = ?))",
				authAccount.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else if authAccount.UserId != 0 {
		tx = tx.Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN calendar_access as ca ON calendar.calendar_id = ca.calendar_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR (calendar_access.recipient_user_id = ? AND ca.state = ?))",
				types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else {
		tx = tx.Joins("LEFT JOIN calendar_access ON calendar.calendar_id = calendar_access.calendar_id AND calendar_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(calendar.visibility_level = ? OR calendar_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	conversion, err := gmodel.CalendarSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing calendar rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbCalendars).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list calendar rows")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Calendar, len(dbCalendars))
	for i, m := range dbCalendars {
		res[i], err = convert.CalendarToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("invalid calendar row when listing calendars")
			return nil, repository.ErrInternal{Msg: "invalid calendar row when listing calendars"}
		}
	}

	return res, nil
}

// UpdateCalendar updates a calendar
func (repo *Client) UpdateCalendar(ctx context.Context, calendar cmodel.Calendar, fields []string) (cmodel.Calendar, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", calendar.CalendarId.CalendarId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.CalendarFromCoreModel(calendar)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar when updating calendar row")
		return cmodel.Calendar{}, repository.ErrInvalidArgument{Msg: "invalid calendar"}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.CalendarFieldMasker.Convert(fields)).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update calendar row")
		return cmodel.Calendar{}, ConvertGormError(err)
	}

	m, err := convert.CalendarToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid calendar row when updating calendar")
		return cmodel.Calendar{}, repository.ErrInternal{Msg: "invalid calendar row when updating calendar"}
	}

	return m, nil
}
