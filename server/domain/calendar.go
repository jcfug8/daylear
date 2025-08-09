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
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when creating a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}
	calendar.CalendarId.CalendarId = 0

	if calendar.Parent.CircleId != 0 {
		_, err = d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: authAccount.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when creating a calendar")
			return model.Calendar{}, err
		}
	} else if calendar.Parent.UserId != 0 {
		_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when creating a calendar")
			return model.Calendar{}, err
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin creating calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to begin creating calendar"}
	}
	defer tx.Rollback()

	dbCalendar, err = tx.CreateCalendar(ctx, calendar, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to create calendar"}
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

	dbCalendarAccess, err := tx.CreateCalendarAccess(ctx, calendarAccess, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create calendar access")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to create calendar access"}
	}

	dbCalendar.CalendarAccess = dbCalendarAccess

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish creating calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to finish creating calendar"}
	}

	dbCalendar.Parent = calendar.Parent

	return dbCalendar, nil
}

// DeleteCalendar deletes a calendar
func (d *Domain) DeleteCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when deleting a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if id.CalendarId == 0 {
		log.Warn().Msg("calendar id required when deleting a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "calendar id required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, id, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when deleting a calendar")
		return model.Calendar{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin deleting calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to begin deleting calendar"}
	}
	defer tx.Rollback()

	dbCalendar, err = tx.DeleteCalendar(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to delete calendar"}
	}

	err = tx.BulkDeleteCalendarAccess(ctx, model.CalendarAccessParent{CalendarId: id.CalendarId})
	if err != nil {
		log.Error().Err(err).Msg("unable to delete calendar access")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to delete calendar access"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish deleting calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to finish deleting calendar"}
	}

	dbCalendar.Parent = parent

	return dbCalendar, nil
}

// GetCalendar retrieves a calendar
func (d *Domain) GetCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId, fields []string) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when getting a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if id.CalendarId == 0 {
		log.Warn().Msg("id required when getting a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	dbCalendar, err = d.repo.GetCalendar(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to get calendar"}
	}

	dbCalendar.CalendarAccess, err = d.determineCalendarAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(dbCalendar.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when getting a calendar")
		return model.Calendar{}, err
	}

	return dbCalendar, nil
}

// ListCalendars lists calendars
func (d *Domain) ListCalendars(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, pageSize int32, offset int64, filter string, fields []string) (dbCalendars []model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when listing calendars")
		return nil, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		dbCircle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: parent.CircleId}, []string{model.CircleField_Visibility})
		if err != nil {
			log.Error().Err(err).Msg("unable to get circle when listing calendars")
			return nil, domain.ErrInternal{Msg: "unable to get circle when listing calendars"}
		}
		_, err = d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: parent.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(dbCircle.VisibilityLevel),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing calendars")
			return nil, err
		}
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		determinedUserAccess, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: authAccount.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing recipes")
			return nil, err
		}
		authAccount.PermissionLevel = determinedUserAccess.GetPermissionLevel()
	}

	dbCalendars, err = d.repo.ListCalendars(ctx, authAccount, pageSize, offset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list calendars")
		return nil, domain.ErrInternal{Msg: "unable to list calendars"}
	}

	return dbCalendars, nil
}

// UpdateCalendar updates a calendar
func (d *Domain) UpdateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, fields []string) (dbCalendar model.Calendar, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when updating a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if calendar.CalendarId.CalendarId == 0 {
		log.Warn().Msg("id required when updating a calendar")
		return model.Calendar{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err = d.determineCalendarAccess(ctx, authAccount, calendar.CalendarId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when updating a calendar")
		return model.Calendar{}, err
	}

	dbCalendar, err = d.repo.UpdateCalendar(ctx, authAccount, calendar, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update calendar")
		return model.Calendar{}, domain.ErrInternal{Msg: "unable to update calendar"}
	}

	dbCalendar.Parent = calendar.Parent

	return dbCalendar, nil
}
