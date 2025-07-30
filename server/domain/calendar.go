package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCalendar creates a new calendar
func (d *Domain) CreateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateCalendar called")
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if calendar.Parent.CircleId != 0 && authAccount.CircleId != 0 && calendar.Parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if calendar.Parent.CircleId != 0 {
		authAccount.CircleId = calendar.Parent.CircleId
	} else if calendar.Parent.UserId != 0 {
		authAccount.UserId = calendar.Parent.UserId
	}

	calendar.CalendarId.CalendarId = 0

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
			return model.Calendar{}, err
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Calendar{}, err
	}
	defer tx.Rollback()

	dbCalendar, err = tx.CreateCalendar(ctx, calendar)
	if err != nil {
		log.Error().Err(err).Msg("repo.CreateCalendar failed")
		return model.Calendar{}, err
	}

	calendarAccess := model.CalendarAccess{
		CalendarAccessParent: model.CalendarAccessParent{
			CalendarId: dbCalendar.CalendarId.CalendarId,
		},
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	if authAccount.CircleId != 0 {
		calendarAccess.Recipient = model.CalendarRecipientOrRequester{
			CircleId: authAccount.CircleId,
		}
	} else {
		calendarAccess.Recipient = model.CalendarRecipientOrRequester{
			UserId: authAccount.AuthUserId,
		}
	}

	dbCalendarAccess, err := tx.CreateCalendarAccess(ctx, calendarAccess)
	if err != nil {
		log.Error().Err(err).Msg("tx.CreateCalendarAccess failed")
		return model.Calendar{}, err
	}

	dbCalendar.CalendarAccess = dbCalendarAccess

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Calendar{}, err
	}

	dbCalendar.Parent = calendar.Parent

	return dbCalendar, nil
}

// DeleteCalendar deletes a calendar
func (d *Domain) DeleteCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain DeleteCalendar called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if id.CalendarId == 0 {
		log.Warn().Msg("id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if parent.CircleId != 0 && authAccount.CircleId != 0 && parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCalendarAccessLevels(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("getCalendarAccessLevels failed")
		return model.Calendar{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Calendar{}, err
	}
	defer tx.Rollback()

	dbCalendar, err = tx.DeleteCalendar(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("tx.DeleteCalendar failed")
		return model.Calendar{}, err
	}

	err = tx.BulkDeleteCalendarAccess(ctx, model.CalendarAccessParent{CalendarId: id.CalendarId})
	if err != nil {
		log.Error().Err(err).Msg("tx.BulkDeleteCalendarAccess failed")
		return model.Calendar{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Calendar{}, err
	}

	dbCalendar.Parent = parent

	return dbCalendar, nil
}

// GetCalendar retrieves a calendar
func (d *Domain) GetCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain GetCalendar called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if id.CalendarId == 0 {
		log.Warn().Msg("id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if parent.CircleId != 0 && authAccount.CircleId != 0 && parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCalendarAccessLevels(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("getCalendarAccessLevels failed")
		return model.Calendar{}, err
	}

	if authAccount.VisibilityLevel != types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC && authAccount.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		log.Warn().Msg("user does not have access")
		return model.Calendar{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	dbCalendar, err = d.repo.GetCalendar(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetCalendar failed")
		return model.Calendar{}, err
	}

	dbCalendar.Parent = parent
	if authAccount.PermissionLevel < dbCalendar.CalendarAccess.PermissionLevel {
		dbCalendar.CalendarAccess.PermissionLevel = authAccount.PermissionLevel
	}

	return dbCalendar, nil
}

// ListCalendars lists calendars
func (d *Domain) ListCalendars(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, pageSize int32, offset int64, filter string, fieldMask []string) (recipe []model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain ListCalendars called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return nil, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if parent.CircleId != 0 && authAccount.CircleId != 0 && parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return nil, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
			return nil, err
		}
	} else if authAccount.UserId != 0 {
		authAccount.PermissionLevel, err = d.getUserAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getUserAccessLevels failed")
			return nil, err
		}
	} else {
		authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	}

	calendars, err := d.repo.ListCalendars(ctx, authAccount, pageSize, offset, filter, fieldMask)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListCalendars failed")
		return nil, err
	}

	for i, calendar := range calendars {
		calendars[i].Parent = parent
		if calendar.CalendarAccess.PermissionLevel > authAccount.PermissionLevel {
			calendars[i].CalendarAccess.PermissionLevel = authAccount.PermissionLevel
		}
	}

	return calendars, nil
}

// UpdateCalendar updates a calendar
func (d *Domain) UpdateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, updateMask []string) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain UpdateCalendar called")

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if calendar.CalendarId.CalendarId == 0 {
		log.Warn().Msg("id required")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if calendar.Parent.CircleId != 0 && authAccount.CircleId != 0 && calendar.Parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if calendar.Parent.CircleId != 0 {
		authAccount.CircleId = calendar.Parent.CircleId
	} else if calendar.Parent.UserId != 0 {
		authAccount.UserId = calendar.Parent.UserId
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.checkCalendarAccess(ctx, authAccount, calendar.CalendarId, types.PermissionLevel_PERMISSION_LEVEL_WRITE)
	if err != nil {
		log.Error().Err(err).Msg("checkCalendarAccess failed")
		return model.Calendar{}, err
	}

	dbCalendar, err = d.repo.UpdateCalendar(ctx, authAccount, calendar, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateCalendar failed")
		return model.Calendar{}, err
	}

	dbCalendar.Parent = calendar.Parent

	return dbCalendar, nil
}
