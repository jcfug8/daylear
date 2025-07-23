package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateUserAccess - Only allows WRITE level
func (d *Domain) CreateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (dbUserAccess model.UserAccess, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain CreateUserAccess called")
	if access.UserAccessParent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if access.Recipient.UserId == 0 {
		log.Warn().Msg("recipient is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	if access.Level != types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("access level must be write")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access level must be write"}
	}

	authAccount.UserId = access.UserAccessParent.UserId.UserId
	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getUserAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getUserAccessLevels failed")
		return model.UserAccess{}, err
	}

	// Prepare both A and B accesses
	user1 := authAccount.AuthUserId
	user2 := authAccount.UserId

	userA := model.UserAccess{
		UserAccessParent: model.UserAccessParent{UserId: model.UserId{UserId: user2}},
		Level:            types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		State:            types.AccessState_ACCESS_STATE_PENDING,
		Requester:        model.UserId{UserId: user1},
		Recipient:        model.UserId{UserId: user1},
	}
	userB := model.UserAccess{
		UserAccessParent: model.UserAccessParent{UserId: model.UserId{UserId: user1}},
		Level:            types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		State:            types.AccessState_ACCESS_STATE_PENDING,
		Requester:        model.UserId{UserId: user1},
		Recipient:        model.UserId{UserId: user2},
	}

	// Create both accesses
	dbUserAccessA, errA := d.repo.CreateUserAccess(ctx, userA)
	if errA != nil {
		log.Error().Err(errA).Msg("repo.CreateUserAccess (A) failed")
		return model.UserAccess{}, errA
	}
	_, errB := d.repo.CreateUserAccess(ctx, userB)
	if errB != nil {
		log.Error().Err(errB).Msg("repo.CreateUserAccess (B) failed")
		// Attempt to clean up A if B fails
		_ = d.repo.DeleteUserAccess(ctx, userA.UserAccessParent, userA.UserAccessId)
		return model.UserAccess{}, errB
	}

	log.Info().Msg("Domain CreateUserAccess returning successfully (B)")
	return dbUserAccessA, nil
}

// DeleteUserAccess -
func (d *Domain) DeleteUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) error {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain DeleteUserAccess called")
	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required")
		return domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required")
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}
	access, err := d.repo.GetUserAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUserAccess failed")
		return err
	}

	if access.Recipient.UserId != authAccount.AuthUserId && access.Requester.UserId != authAccount.AuthUserId && access.UserAccessParent.UserId.UserId != authAccount.AuthUserId {
		log.Warn().Msg("user does not have access to delete this user access")
		return domain.ErrPermissionDenied{Msg: "user does not have access to delete this user access"}
	}
	// Delete the current access
	err = d.repo.DeleteUserAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.DeleteUserAccess failed")
		return err
	}
	// Delete the sister access (swap UserId and Recipient)
	sisterParent := model.UserAccessParent{UserId: access.Recipient}
	sisterFilter := fmt.Sprintf("requester_user_id=%d AND recipient_user_id=%d", access.Requester.UserId, parent.UserId.UserId)
	// Find the sister access id
	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, sisterParent, 1, 0, sisterFilter)
	if err != nil {
	}
	if len(accesses) == 0 {
		log.Warn().Msg("sister access not found on delete")
	}

	err = d.repo.DeleteUserAccess(ctx, sisterParent, accesses[0].UserAccessId)
	if err != nil {
		log.Error().Err(err).Msg("repo.DeleteUserAccess failed (sister)")
		return err
	}

	log.Info().Msg("Domain DeleteUserAccess returning successfully (and deleted sister if found)")
	return nil
}

// GetUserAccess -
func (d *Domain) GetUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain GetUserAccess called")
	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	access, err := d.repo.GetUserAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUserAccess failed")
		return model.UserAccess{}, err
	}

	if access.Recipient.UserId != authAccount.AuthUserId && access.Requester.UserId != authAccount.AuthUserId && access.UserAccessParent.UserId.UserId != authAccount.AuthUserId {
		log.Warn().Msg("user does not have access to view this user access")
		return model.UserAccess{}, domain.ErrPermissionDenied{Msg: "user does not have access to view this user access"}
	}

	log.Info().Msg("Domain GetUserAccess returning successfully")
	return access, nil
}

// ListUserAccesses -
func (d *Domain) ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain ListUserAccesses called")
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("requester is required")
		return nil, domain.ErrInvalidArgument{Msg: "requester is required"}
	}

	if parent.UserId.UserId != 0 {
		perm, _, err := d.getUserAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getUserAccessLevels failed")
			return nil, err
		}
		if perm < types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
			log.Warn().Msg("user does not have access to list accesses")
			return nil, domain.ErrPermissionDenied{Msg: "user does not have access to list accesses"}
		}
	}

	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, parent, pageSize, pageOffset, filter)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListUserAccesses failed")
		return nil, err
	}
	log.Info().Msg("Domain ListUserAccesses returning successfully")
	return accesses, nil
}

// AcceptUserAccess -
func (d *Domain) AcceptUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	log.Info().Msg("Domain AcceptUserAccess called")
	if parent.UserId.UserId == 0 {
		log.Warn().Msg("user id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}
	if id.UserAccessId == 0 {
		log.Warn().Msg("access id is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access id is required"}
	}
	access, err := d.repo.GetUserAccess(ctx, parent, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetUserAccess failed")
		return model.UserAccess{}, err
	}
	if access.State != types.AccessState_ACCESS_STATE_PENDING {
		log.Warn().Msg("access must be in pending state to be accepted")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "access must be in pending state to be accepted"}
	}
	if access.Recipient.UserId != authAccount.AuthUserId || access.UserAccessParent.UserId.UserId == access.Recipient.UserId {
		log.Warn().Msg("only the recipient (B) can accept this access")
		return model.UserAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient (B) can accept this access"}
	}
	// Accept B
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED
	updated, err := d.repo.UpdateUserAccess(ctx, access)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateUserAccess failed (B)")
		return model.UserAccess{}, err
	}
	// Accept sister A (swap UserId and Recipient)
	sisterParent := model.UserAccessParent{UserId: access.Recipient}
	sisterFilter := fmt.Sprintf("requester_user_id=%d AND recipient_user_id=%d", access.Requester.UserId, parent.UserId.UserId)
	accesses, err := d.repo.ListUserAccesses(ctx, authAccount, sisterParent, 1, 0, sisterFilter)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListUserAccesses failed (sister)")
		return model.UserAccess{}, err
	}
	if len(accesses) == 0 {
		log.Warn().Msg("sister access not found on accept")
	}

	accesses[0].State = types.AccessState_ACCESS_STATE_ACCEPTED
	_, err = d.repo.UpdateUserAccess(ctx, accesses[0])
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateUserAccess failed (sister)")
		return model.UserAccess{}, err
	}

	log.Info().Msg("Domain AcceptUserAccess returning successfully (and accepted sister if found)")
	return updated, nil
}
