package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateCalendarAccess creates calendar access
func (repo *Client) CreateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess, fields []string) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	// Validate that exactly one recipient type is set
	if (access.Recipient.UserId != 0) == (access.Recipient.CircleId != 0) {
		log.Error().Msg("exactly one recipient (user or circle) is required to create calendar access row")
		return cmodel.CalendarAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	gormAccess := convert.CalendarAccessToGorm(access)
	res := repo.db.WithContext(ctx).Select(dbModel.CalendarAccessFieldMasker.Convert(fields)).Clauses(clause.Returning{}).Create(&gormAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("unable to create calendar access row: already exists")
			return cmodel.CalendarAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("unable to create calendar access row")
		return cmodel.CalendarAccess{}, ConvertGormError(res.Error)
	}

	access.CalendarAccessId.CalendarAccessId = gormAccess.CalendarAccessId
	return access, nil
}

// DeleteCalendarAccess deletes calendar access
func (repo *Client) DeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent, id cmodel.CalendarAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", parent.CalendarId).
		Int64("calendarAccessId", id.CalendarAccessId).
		Logger()

	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required to delete calendar access row")
		return repository.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	if id.CalendarAccessId == 0 {
		log.Error().Msg("calendar access id is required to delete calendar access row")
		return repository.ErrInvalidArgument{Msg: "calendar access id is required"}
	}

	res := repo.db.WithContext(ctx).Where("calendar_id = ? AND calendar_access_id = ?", parent.CalendarId, id.CalendarAccessId).Delete(&dbModel.CalendarAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to delete calendar access row")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no calendar access row deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

// BulkDeleteCalendarAccesses bulk deletes calendar accesses
func (repo *Client) BulkDeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", parent.CalendarId).
		Logger()

	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required to bulk delete calendar access rows")
		return repository.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	res := repo.db.WithContext(ctx).Where("calendar_id = ?", parent.CalendarId).Delete(&dbModel.CalendarAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to bulk delete calendar access rows")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no calendar access rows deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

// GetCalendarAccess retrieves calendar access
func (repo *Client) GetCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent, id cmodel.CalendarAccessId, fields []string) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", parent.CalendarId).
		Int64("calendarAccessId", id.CalendarAccessId).
		Strs("fields", fields).
		Logger()

	dbFields := dbModel.CalendarAccessFieldMasker.Convert(fields)

	var calendarAccess dbModel.CalendarAccess
	tx := repo.db.WithContext(ctx).Select(dbFields)
	if fieldmask.ContainsAny(dbFields, dbModel.UserColumn_Username, dbModel.UserColumn_GivenName, dbModel.UserColumn_FamilyName) {
		tx = tx.Joins(`LEFT JOIN daylear_user ON calendar_access.recipient_user_id = daylear_user.user_id`)
	}
	if fieldmask.ContainsAny(dbFields, dbModel.CircleColumn_Title, dbModel.CircleColumn_Handle) {
		tx = tx.Joins(`LEFT JOIN circle ON calendar_access.recipient_circle_id = circle.circle_id`)
	}

	res := tx.Where("calendar_access.calendar_id = ? AND calendar_access.calendar_access_id = ?", parent.CalendarId, id.CalendarAccessId).
		First(&calendarAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("calendar access row not found")
			return cmodel.CalendarAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to get calendar access row")
		return cmodel.CalendarAccess{}, ConvertGormError(res.Error)
	}

	return convert.CalendarAccessFromGorm(calendarAccess), nil
}

// ListCalendarAccesses lists calendar accesses
func (repo *Client) ListCalendarAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.CalendarAccessParent, pageSize int32, pageOffset int64, filterStr string, fields []string) ([]cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", parent.CalendarId).
		Str("filter", filterStr).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).Logger()

	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Error().Msg("user id or circle id is required to list calendar access rows")
		return nil, repository.ErrInvalidArgument{Msg: "user id or circle id is required"}
	}

	dbFields := dbModel.CalendarAccessFieldMasker.Convert(fields)

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "calendar_access.calendar_access_id"},
		Desc:   true,
	}}

	var calendarAccesses []dbModel.CalendarAccess
	// Start building the query
	tx := repo.db.WithContext(ctx).
		Select(dbFields).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	if fieldmask.ContainsAny(dbFields, dbModel.UserColumn_Username, dbModel.UserColumn_GivenName, dbModel.UserColumn_FamilyName) {
		tx = tx.Joins(`LEFT JOIN daylear_user ON calendar_access.recipient_user_id = daylear_user.user_id`)
	}
	if fieldmask.ContainsAny(dbFields, dbModel.CircleColumn_Title, dbModel.CircleColumn_Handle) {
		tx = tx.Joins(`LEFT JOIN circle ON calendar_access.recipient_circle_id = circle.circle_id`)
	}

	// Filter by calendar ID if provided
	if parent.CalendarId != 0 {
		tx = tx.Where("calendar_access.calendar_id = ?", parent.CalendarId)
	}

	if authAccount.CircleId != 0 {
		tx = tx.Where(
			"calendar_access.recipient_circle_id = ? OR calendar_access.calendar_id IN (SELECT calendar_id FROM calendar_access WHERE recipient_circle_id = ? AND permission_level >= ?)",
			authAccount.CircleId, authAccount.CircleId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	} else if authAccount.AuthUserId != 0 {
		tx = tx.Where(
			"calendar_access.recipient_user_id = ? OR calendar_access.calendar_id IN (SELECT calendar_id FROM calendar_access WHERE recipient_user_id = ? AND permission_level >= ?)",
			authAccount.AuthUserId, authAccount.AuthUserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	}

	conversion, err := dbModel.CalendarAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing calendar access rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	res := tx.Find(&calendarAccesses)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to list calendar access rows")
		return nil, ConvertGormError(res.Error)
	}

	accesses := make([]cmodel.CalendarAccess, len(calendarAccesses))
	for i, access := range calendarAccesses {
		accesses[i] = convert.CalendarAccessFromGorm(access)
	}

	return accesses, nil
}

// UpdateCalendarAccess updates calendar access
func (repo *Client) UpdateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess, fields []string) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarAccessId", access.CalendarAccessId.CalendarAccessId).
		Strs("fields", fields).
		Logger()

	dbAccess := convert.CalendarAccessToGorm(access)

	res := repo.db.WithContext(ctx).
		Select(dbModel.CalendarAccessFieldMasker.Convert(fields, fieldmask.OnlyUpdatable())).
		Clauses(&clause.Returning{}).
		Where("calendar_access_id = ?", access.CalendarAccessId.CalendarAccessId).
		Updates(&dbAccess)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update calendar access row")
		return cmodel.CalendarAccess{}, ConvertGormError(res.Error)
	}

	return convert.CalendarAccessFromGorm(dbAccess), nil
}

func (repo *Client) FindStandardUserCalendarAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CalendarId) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from calendar_access where recipient_user_id = ? and calendar_id = ?

	var calendarAccess dbModel.CalendarAccess
	res := repo.db.WithContext(ctx).
		Where("calendar_id = ? AND recipient_user_id = ?", id.CalendarId, authAccount.AuthUserId).
		First(&calendarAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("standard user calendar access not found")
			return cmodel.CalendarAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find standard user calendar access")
		return cmodel.CalendarAccess{}, ConvertGormError(res.Error)
	}
	return convert.CalendarAccessFromGorm(calendarAccess), nil
}

func (repo *Client) FindDelegatedCircleCalendarAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CalendarId) (cmodel.CalendarAccess, cmodel.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from calendar_access
	// 	JOIN circle_access ON circle_access.circle_id = calendar_access.recipient_circle_id
	// WHERE calendar_access.calendar_id = 1 AND circle_access.recipient_user_id = 1 LIMIT 1;

	type Result struct {
		dbModel.CalendarAccess
		dbModel.CircleAccess
	}
	var result Result
	res := repo.db.WithContext(ctx).
		Select("calendar_access.*, circle_access.*, LEAST(calendar_access.permission_level, circle_access.permission_level) as effective_permission_level").
		Table("calendar_access").
		Joins("JOIN circle_access ON circle_access.circle_id = calendar_access.recipient_circle_id").
		Where("calendar_access.calendar_id = ? AND circle_access.recipient_user_id = ? AND calendar_access.state = ?", id.CalendarId, authAccount.AuthUserId, types.AccessState_ACCESS_STATE_ACCEPTED).
		Order("effective_permission_level DESC").
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated circle calendar access not found")
			return cmodel.CalendarAccess{}, cmodel.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated circle calendar access")
		return cmodel.CalendarAccess{}, cmodel.CircleAccess{}, ConvertGormError(res.Error)
	}
	return convert.CalendarAccessFromGorm(result.CalendarAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserCalendarAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CalendarId) (cmodel.CalendarAccess, cmodel.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from calendar_access
	// 	JOIN user_access ON user_access.user_id = calendar_access.recipient_user_id
	// WHERE calendar_access.calendar_id = ? AND user_access.recipient_user_id = ? LIMIT 1;

	type Result struct {
		dbModel.CalendarAccess
		dbModel.UserAccess
	}
	var result Result

	res := repo.db.WithContext(ctx).
		Select("calendar_access.*, user_access.*").
		Table("calendar_access").
		Joins("JOIN user_access ON user_access.user_id = calendar_access.recipient_user_id").
		Where("calendar_access.calendar_id = ? AND user_access.recipient_user_id = ? AND calendar_access.state = ?", id.CalendarId, authAccount.AuthUserId, types.AccessState_ACCESS_STATE_ACCEPTED).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated user calendar access not found")
			return cmodel.CalendarAccess{}, cmodel.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated user calendar access")
		return cmodel.CalendarAccess{}, cmodel.UserAccess{}, ConvertGormError(res.Error)
	}
	return convert.CalendarAccessFromGorm(result.CalendarAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
