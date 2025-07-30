package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateCalendarAccessMap maps the updatable core model fields to the database model fields for the CalendarAccess model.
var UpdateCalendarAccessMap = map[string][]string{
	model.CalendarAccessFields.Level: []string{
		dbModel.CalendarAccessFields.PermissionLevel,
	},
	model.CalendarAccessFields.State: []string{
		dbModel.CalendarAccessFields.State,
	},
}

var UpdateCalendarAccessFieldMasker = fieldmask.NewFieldMasker(UpdateCalendarAccessMap)

// CalendarAccessMap maps the core model fields to the database model fields for the unified CalendarAccess model.
var CalendarAccessMap = map[string]string{
	model.CalendarAccessFields.Level:           dbModel.CalendarAccessFields.PermissionLevel,
	model.CalendarAccessFields.State:           dbModel.CalendarAccessFields.State,
	model.CalendarAccessFields.RecipientUser:   dbModel.CalendarAccessFields.RecipientUserId,
	model.CalendarAccessFields.RecipientCircle: dbModel.CalendarAccessFields.RecipientCircleId,
}

// CreateCalendarAccess creates calendar access
func (c *Client) CreateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	db := c.db.WithContext(ctx)

	// Validate that exactly one recipient type is set
	if (access.Recipient.UserId != 0) == (access.Recipient.CircleId != 0) {
		log.Error().Msg("exactly one recipient (user or circle) is required")
		return cmodel.CalendarAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	gormAccess := convert.CalendarAccessToGorm(access)
	res := db.Clauses(clause.Returning{}).Create(&gormAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("duplicate key error")
			return cmodel.CalendarAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("db.Create failed")
		return cmodel.CalendarAccess{}, res.Error
	}

	access.CalendarAccessId.CalendarAccessId = gormAccess.CalendarAccessId
	return access, nil
}

// DeleteCalendarAccess deletes calendar access
func (c *Client) DeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent, id cmodel.CalendarAccessId) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required")
		return repository.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	if id.CalendarAccessId == 0 {
		log.Error().Msg("calendar access id is required")
		return repository.ErrInvalidArgument{Msg: "calendar access id is required"}
	}

	db := c.db.WithContext(ctx)

	res := db.Where("calendar_id = ? AND calendar_access_id = ?", parent.CalendarId, id.CalendarAccessId).Delete(&dbModel.CalendarAccess{})
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

// BulkDeleteCalendarAccesses bulk deletes calendar accesses
func (c *Client) BulkDeleteCalendarAccess(ctx context.Context, parent cmodel.CalendarAccessParent) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	if parent.CalendarId == 0 {
		log.Error().Msg("calendar id is required")
		return repository.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	db := c.db.WithContext(ctx)

	res := db.Where("calendar_id = ?", parent.CalendarId).Delete(&dbModel.CalendarAccess{})
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
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	db := c.db.WithContext(ctx)

	var calendarAccess dbModel.CalendarAccess
	res := db.Table("calendar_access").
		Select(`calendar_access.*, u.username as recipient_username, u.given_name as recipient_given_name, u.family_name as recipient_family_name, c.title as recipient_circle_title, c.handle as recipient_circle_handle`).
		Joins(`LEFT JOIN daylear_user u ON calendar_access.recipient_user_id = u.user_id`).
		Joins(`LEFT JOIN circle c ON calendar_access.recipient_circle_id = c.circle_id`).
		Where("calendar_access.calendar_id = ? AND calendar_access.calendar_access_id = ?", parent.CalendarId, id.CalendarAccessId).
		First(&calendarAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("record not found")
			return cmodel.CalendarAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("db.First failed")
		return cmodel.CalendarAccess{}, res.Error
	}

	return convert.CalendarAccessFromGorm(calendarAccess), nil
}

// ListCalendarAccesses lists calendar accesses
func (c *Client) ListCalendarAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.CalendarAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		return nil, repository.ErrInvalidArgument{Msg: "user id or circle id is required"}
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "calendar_access.calendar_access_id"},
		Desc:   true,
	}}

	var calendarAccesses []dbModel.CalendarAccess
	// Start building the query
	db := c.db.WithContext(ctx).
		Table("calendar_access").
		Select(`calendar_access.*, u.username as recipient_username, u.given_name as recipient_given_name, u.family_name as recipient_family_name, c.title as recipient_circle_title, c.handle as recipient_circle_handle`).
		Joins(`LEFT JOIN daylear_user u ON calendar_access.recipient_user_id = u.user_id`).
		Joins(`LEFT JOIN circle c ON calendar_access.recipient_circle_id = c.circle_id`).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	// Filter by calendar ID if provided
	if parent.CalendarId != 0 {
		db = db.Where("calendar_access.calendar_id = ?", parent.CalendarId)
	}

	if authAccount.CircleId != 0 {
		db = db.Where(
			"calendar_access.recipient_circle_id = ? OR calendar_access.calendar_id IN (SELECT calendar_id FROM calendar_access WHERE recipient_circle_id = ? AND permission_level >= ?)",
			authAccount.CircleId, authAccount.CircleId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	} else if authAccount.AuthUserId != 0 {
		db = db.Where(
			"calendar_access.recipient_user_id = ? OR calendar_access.calendar_id IN (SELECT calendar_id FROM calendar_access WHERE recipient_user_id = ? AND permission_level >= ?)",
			authAccount.AuthUserId, authAccount.AuthUserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	}

	conversion, err := c.calendarAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&calendarAccesses).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	accesses := make([]cmodel.CalendarAccess, len(calendarAccesses))
	for i, access := range calendarAccesses {
		accesses[i] = convert.CalendarAccessFromGorm(access)
	}

	return accesses, nil
}

// UpdateCalendarAccess updates calendar access
func (c *Client) UpdateCalendarAccess(ctx context.Context, access cmodel.CalendarAccess, updateMask []string) (cmodel.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	dbAccess := convert.CalendarAccessToGorm(access)

	columns := UpdateCalendarAccessFieldMasker.Convert(updateMask)

	db := c.db.WithContext(ctx).Select(columns).Clauses(&clause.Returning{})

	err := db.Where("calendar_access_id = ?", access.CalendarAccessId.CalendarAccessId).Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return cmodel.CalendarAccess{}, ConvertGormError(err)
	}

	return convert.CalendarAccessFromGorm(dbAccess), nil
}
