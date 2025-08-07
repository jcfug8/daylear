package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCalendarAccess creates calendar access
func (d *Domain) CreateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess) (dbAccess model.CalendarAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.CalendarAccessParent.CalendarId == 0 {
		log.Warn().Msg("calendar id required when creating a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "calendar id required"}
	}

	var recipeientOwner bool
	var resourceOwner bool

	access.Requester = model.CalendarRecipientOrRequester{
		UserId: authAccount.AuthUserId,
	}

	dbCalendar, err := d.repo.GetCalendar(ctx, model.CalendarId{CalendarId: access.CalendarAccessParent.CalendarId}, []string{model.CalendarField_Visibility})
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar when creating a calendar access")
		return model.CalendarAccess{}, err
	}
	determinedCalendarAccess, err := d.determineCalendarAccess(
		ctx, authAccount, model.CalendarId{CalendarId: access.CalendarAccessParent.CalendarId},
		withResourceVisibilityLevel(dbCalendar.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access when creating a calendar access")
		return model.CalendarAccess{}, err
	}
	resourceOwner = determinedCalendarAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE

	if access.Recipient.CircleId != 0 { // recipient is a circle
		dbCircle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: access.Recipient.CircleId}, []string{model.CircleField_Visibility})
		if err != nil {
			log.Error().Err(err).Msg("unable to get circle when creating a calendar access")
			return model.CalendarAccess{}, err
		}
		determinedCircleAccess, err := d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: access.Recipient.CircleId}, withResourceVisibilityLevel(dbCircle.VisibilityLevel))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access when creating a calendar access")
			return model.CalendarAccess{}, err
		}
		recipeientOwner = determinedCircleAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_WRITE
	} else if access.Recipient.UserId != 0 { // recipient is a different user
		determinedUserAccess, err := d.determineUserAccess(ctx, authAccount, model.UserId{UserId: access.Recipient.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access when creating a calendar access")
			return model.CalendarAccess{}, err
		}
		recipeientOwner = determinedUserAccess.PermissionLevel >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	} else {
		log.Warn().Msg("recipient is required when creating a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	if !resourceOwner && recipeientOwner {
		access.State = types.AccessState_ACCESS_STATE_PENDING
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_RESOURCE
	} else if resourceOwner && !recipeientOwner {
		access.State = types.AccessState_ACCESS_STATE_PENDING
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_RECIPIENT
	} else if resourceOwner && recipeientOwner {
		access.State = types.AccessState_ACCESS_STATE_ACCEPTED
		access.AcceptTarget = types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED
	} else {
		return model.CalendarAccess{}, domain.ErrPermissionDenied{Msg: "incorrect access"}
	}

	if access.PermissionLevel > max(types.PermissionLevel_PERMISSION_LEVEL_READ, determinedCalendarAccess.PermissionLevel) {
		log.Warn().Msg("unable to create calendar access with the given permission level")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "cannot create access level higher than your own level"}
	}

	// create access
	dbAccess, err = d.repo.CreateCalendarAccess(ctx, access, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create calendar access when creating a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to create calendar access"}
	}

	return dbAccess, nil
}

// DeleteCalendarAccess deletes calendar access
func (d *Domain) DeleteCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.CalendarId == 0 {
		log.Warn().Msg("calendar id is required when deleting a calendar access")
		return domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	// verify access id is set
	if id.CalendarAccessId == 0 {
		log.Warn().Msg("access id is required when deleting a calendar access")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	dbAccess, err := d.repo.GetCalendarAccess(ctx, parent, id, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar access when deleting a calendar access")
		return domain.ErrInternal{Msg: "unable to get calendar access"}
	}

	if dbAccess.Recipient.CircleId != 0 {
		_, err = d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: dbAccess.Recipient.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access when deleting a calendar access")
			return err
		}
	} else if dbAccess.Recipient.UserId != 0 {
		_, err = d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: dbAccess.Recipient.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access when deleting a calendar access")
			return err
		}
	} else {
		_, err := d.determineCalendarAccess(
			ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine calendar access when deleting a calendar access")
			return err
		}
	}

	err = d.repo.DeleteCalendarAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete calendar access when deleting a calendar access")
		return domain.ErrInternal{Msg: "unable to delete calendar access"}
	}

	return nil
}

// GetCalendarAccess retrieves calendar access
func (d *Domain) GetCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId, fields []string) (model.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if parent.CalendarId == 0 {
		log.Warn().Msg("calendar id is required when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	// verify access id is set
	if id.CalendarAccessId == 0 {
		log.Warn().Msg("access id is required when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	dbAccess, err := d.repo.GetCalendarAccess(ctx, parent, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar access when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to get calendar access"}
	}

	if dbAccess.Recipient.CircleId != 0 {
		_, err = d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: dbAccess.Recipient.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine circle access when getting a calendar access")
			return model.CalendarAccess{}, err
		}
	} else if dbAccess.Recipient.UserId != 0 {
		_, err = d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: dbAccess.Recipient.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine user access when getting a calendar access")
			return model.CalendarAccess{}, err
		}
	} else {
		_, err := d.determineCalendarAccess(
			ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine calendar access when getting a calendar access")
			return model.CalendarAccess{}, err
		}
	}

	return dbAccess, nil
}

// ListCalendarAccesses lists calendar accesses
func (d *Domain) ListCalendarAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) (dbAccesses []model.CalendarAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when listing calendars accesses")
		return nil, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		_, err = d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: parent.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing calendars accesses")
			return nil, err
		}
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: parent.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing calendars accesses")
			return nil, err
		}
	}

	dbAccesses, err = d.repo.ListCalendarAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list calendars accesses")
		return nil, domain.ErrInternal{Msg: "unable to list calendars accesses"}
	}

	return dbAccesses, nil
}

// UpdateCalendarAccess updates calendar access
func (d *Domain) UpdateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess, fields []string) (model.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if access.CalendarAccessParent.CalendarId == 0 {
		log.Warn().Msg("calendar id is required when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	// verify access id is set
	if access.CalendarAccessId.CalendarAccessId == 0 {
		log.Warn().Msg("access id is required when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	determinedCalendarAccess, err := d.determineCalendarAccess(
		ctx, authAccount, model.CalendarId{CalendarId: access.CalendarAccessParent.CalendarId},
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access when getting a calendar access")
		return model.CalendarAccess{}, err
	}

	if slices.Contains(fields, model.CalendarAccessField_PermissionLevel) && determinedCalendarAccess.PermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update calendar access permission level to a higher level than your own")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "cannot update calendar access permission level to a higher level than your own"}
	}

	return d.repo.UpdateCalendarAccess(ctx, access, fields)
}

// AcceptCalendarAccess accepts calendar access
func (d *Domain) AcceptCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) (model.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	// verify calendar is set
	if parent.CalendarId == 0 {
		log.Warn().Msg("calendar id is required when accepting calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	// verify access id is set
	if id.CalendarAccessId == 0 {
		log.Warn().Msg("access id is required when accepting calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	access, err := d.repo.GetCalendarAccess(ctx, parent, id, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar access when accepting calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to get calendar access when accepting calendar access"}
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	switch access.AcceptTarget {
	case types.AcceptTarget_ACCEPT_TARGET_RESOURCE:
		_, err = d.determineCalendarAccess(
			ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine calendar access when accepting calendar access")
			return model.CalendarAccess{}, err
		}
	case types.AcceptTarget_ACCEPT_TARGET_RECIPIENT:
		if access.Recipient.CircleId != 0 {
			_, err = d.determineCircleAccess(
				ctx, authAccount, model.CircleId{CircleId: access.Recipient.CircleId},
				withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
			)
			if err != nil {
				log.Error().Err(err).Msg("unable to determine circle access when accepting calendar access")
				return model.CalendarAccess{}, err
			}
		} else if access.Recipient.UserId != 0 {
			_, err = d.determineUserAccess(
				ctx, authAccount, model.UserId{UserId: access.Recipient.UserId},
				withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN),
			)
			if err != nil {
				log.Error().Err(err).Msg("unable to determine user access when accepting calendar access")
				return model.CalendarAccess{}, err
			}
		} else {
			return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "invalid recipient"}
		}
	default:
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "invalid accept target"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateCalendarAccess(ctx, access, []string{model.CalendarAccessField_State})
	if err != nil {
		log.Error().Err(err).Msg("unable to update calendar access when accepting calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to update calendar access when accepting calendar access"}
	}

	return updatedAccess, nil
}
