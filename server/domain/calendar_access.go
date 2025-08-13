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

	determinedCalendarAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, access)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access ownership details when creating a calendar access")
		return model.CalendarAccess{}, err
	}

	access.Requester = model.CalendarRecipientOrRequester{
		UserId: authAccount.AuthUserId,
	}
	access.AcceptTarget = determinedCalendarAccessOwnershipDetails.acceptTarget
	access.State = determinedCalendarAccessOwnershipDetails.accessState

	if access.PermissionLevel > determinedCalendarAccessOwnershipDetails.maximumPermissionLevel {
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

	determinedCalendarAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access ownership details when deleting a calendar access")
		return domain.ErrInternal{Msg: "unable to determine calendar access ownership details"}
	}

	if !determinedCalendarAccessOwnershipDetails.isRecipientOwner && !determinedCalendarAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when deleting a calendar access")
		return domain.ErrPermissionDenied{Msg: "access denied"}
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

	determinedCalendarAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access ownership details when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to determine calendar access ownership details"}
	}

	if !determinedCalendarAccessOwnershipDetails.isRecipientOwner && !determinedCalendarAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when getting a calendar access")
		return model.CalendarAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
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

	dbAccess, err := d.repo.GetCalendarAccess(ctx, access.CalendarAccessParent, access.CalendarAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar access")
		return model.CalendarAccess{}, err
	}

	determinedCalendarAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(ctx, authAccount, dbAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access ownership details when updating a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to determine calendar access ownership details"}
	}

	if !determinedCalendarAccessOwnershipDetails.isRecipientOwner && !determinedCalendarAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("access denied when updating a calendar access")
		return model.CalendarAccess{}, domain.ErrPermissionDenied{Msg: "access denied"}
	}

	if slices.Contains(fields, model.CalendarAccessField_PermissionLevel) && determinedCalendarAccessOwnershipDetails.maximumPermissionLevel < access.PermissionLevel {
		log.Warn().Msg("cannot update calendar access permission level to a higher level than your own")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "cannot update calendar access permission level to a higher level than your own"}
	}

	updatedAccess, err := d.repo.UpdateCalendarAccess(ctx, access, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update calendar access when updating a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to update calendar access"}
	}

	return updatedAccess, nil
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

	determinedCalendarAccessOwnershipDetails, err := d.determineAccessOwnershipDetails(
		ctx, authAccount, access,
		withAllowAutoOmitAccessChecks(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access ownership details when accepting a calendar access")
		return model.CalendarAccess{}, domain.ErrInternal{Msg: "unable to determine calendar access ownership details"}
	}

	if determinedCalendarAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RESOURCE && !determinedCalendarAccessOwnershipDetails.isResourceOwner {
		log.Warn().Msg("must be resource owner to accept resourse targeted calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "must be resource owner to accept resourse targeted calendar access"}
	} else if determinedCalendarAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT && !determinedCalendarAccessOwnershipDetails.isRecipientOwner {
		log.Warn().Msg("must be recipient owner to accept recipient targeted calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "must be recipient owner to accept recipient targeted calendar access"}
	} else if determinedCalendarAccessOwnershipDetails.acceptTarget == types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED {
		log.Warn().Msg("unspecified accept target when accepting a calendar access")
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "unspecified accept target"}
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
