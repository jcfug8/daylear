package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateCalendarAccess creates calendar access
func (d *Domain) CreateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess) (model.CalendarAccess, error) {
	// TODO: Add authorization logic
	return d.repo.CreateCalendarAccess(ctx, access)
}

// DeleteCalendarAccess deletes calendar access
func (d *Domain) DeleteCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) error {
	// TODO: Add authorization logic
	return d.repo.DeleteCalendarAccess(ctx, parent, id)
}

// GetCalendarAccess retrieves calendar access
func (d *Domain) GetCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) (model.CalendarAccess, error) {
	// TODO: Add authorization logic
	return d.repo.GetCalendarAccess(ctx, parent, id)
}

// ListCalendarAccesses lists calendar accesses
func (d *Domain) ListCalendarAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.CalendarAccess, error) {
	// TODO: Add authorization logic
	return d.repo.ListCalendarAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter)
}

// UpdateCalendarAccess updates calendar access
func (d *Domain) UpdateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess, updateMask []string) (model.CalendarAccess, error) {
	// TODO: Add authorization logic
	return d.repo.UpdateCalendarAccess(ctx, access, updateMask)
}

// AcceptCalendarAccess accepts calendar access
func (d *Domain) AcceptCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) (model.CalendarAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain AcceptCalendarAccess called")

	// verify calendar is set
	if parent.CalendarId == 0 {
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "calendar id is required"}
	}

	// verify access id is set
	if id.CalendarAccessId == 0 {
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// get the current access
	access, err := d.repo.GetCalendarAccess(ctx, parent, id)
	if err != nil {
		return model.CalendarAccess{}, err
	}

	// verify the access is in pending state
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		return model.CalendarAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}

	// verify the user is the recipient of this access
	if access.Recipient.UserId != authAccount.AuthUserId {
		return model.CalendarAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
	}

	// update the access state to accepted
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED

	// update access using the repository
	updatedAccess, err := d.repo.UpdateCalendarAccess(ctx, access, []string{model.CalendarAccessFields.State})
	if err != nil {
		return model.CalendarAccess{}, err
	}

	log.Info().Msg("Domain AcceptCalendarAccess returning successfully")
	return updatedAccess, nil
}
