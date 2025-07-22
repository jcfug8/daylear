package domain

import (
	"context"

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
	if access.Recipient == 0 {
		log.Warn().Msg("recipient is required")
		return model.UserAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	authAccount.UserId = access.UserAccessParent.UserId.UserId
	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getUserAccessLevels(ctx, authAccount)
	if err != nil {
		log.Error().Err(err).Msg("getUserAccessLevels failed")
		return model.UserAccess{}, err
	}

	access.State = types.AccessState_ACCESS_STATE_PENDING
	access.Requester = authAccount.AuthUserId

	dbUserAccess, err = d.repo.CreateUserAccess(ctx, access)
	if err != nil {
		log.Error().Err(err).Msg("repo.CreateUserAccess failed")
		return model.UserAccess{}, err
	}

	log.Info().Msg("Domain CreateUserAccess returning successfully")
	return dbUserAccess, nil
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

	if access.Recipient == authAccount.AuthUserId || access.Requester == authAccount.AuthUserId {
		log.Warn().Msg("user does not have access to delete this user access")
		return domain.ErrPermissionDenied{Msg: "user does not have access to delete this user access"}
	}
	if err := d.repo.DeleteUserAccess(ctx, parent, id); err != nil {
		log.Error().Err(err).Msg("repo.DeleteUserAccess failed")
		return err
	}
	log.Info().Msg("Domain DeleteUserAccess returning successfully")
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

	if access.Recipient != authAccount.AuthUserId && access.Requester != authAccount.AuthUserId {
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
	if access.Recipient != authAccount.AuthUserId {
		log.Warn().Msg("only the recipient can accept this access")
		return model.UserAccess{}, domain.ErrPermissionDenied{Msg: "only the recipient can accept this access"}
	}
	access.State = types.AccessState_ACCESS_STATE_ACCEPTED
	updated, err := d.repo.UpdateUserAccess(ctx, access)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateUserAccess failed")
		return model.UserAccess{}, err
	}
	log.Info().Msg("Domain AcceptUserAccess returning successfully")
	return updated, nil
}
