package gorm

import (
	"context"
	"errors"

	"slices"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Client) CreateUserAccess(ctx context.Context, access model.UserAccess, fields []string) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", access.UserId.UserId).
		Int64("recipientUserId", access.Recipient.UserId).
		Strs("fields", fields).
		Logger()

	if access.Recipient.UserId == 0 {
		log.Error().Msg("recipient is required to create user access row")
		return model.UserAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	userAccess := convert.CoreUserAccessToUserAccess(access)
	res := repo.db.WithContext(ctx).
		Select(dbModel.UserAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&userAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("unable to create user access row: already exists")
			return model.UserAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("unable to create user access row")
		return model.UserAccess{}, ConvertGormError(res.Error)
	}

	access.UserAccessId.UserAccessId = userAccess.UserAccessId
	return access, nil
}

func (repo *Client) DeleteUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", parent.UserId.UserId).
		Int64("userAccessId", id.UserAccessId).
		Logger()

	if parent.UserId.UserId == 0 {
		log.Error().Msg("user id is required to delete user access row")
		return repository.ErrInvalidArgument{Msg: "user id is required"}
	}

	if id.UserAccessId == 0 {
		log.Error().Msg("user access id is required to delete user access row")
		return repository.ErrInvalidArgument{Msg: "user access id is required"}
	}

	res := repo.db.WithContext(ctx).
		Where("user_id = ? AND user_access_id = ?", parent.UserId.UserId, id.UserAccessId).
		Delete(&dbModel.UserAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to delete user access row")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no user access row deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId, fields []string) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", parent.UserId.UserId).
		Int64("userAccessId", id.UserAccessId).
		Strs("fields", fields).
		Logger()

	fields = dbModel.UserAccessFieldMasker.Convert(fields)

	var userAccess dbModel.UserAccess
	tx := repo.db.WithContext(ctx).
		Select(fields)

	if slices.Contains(fields, dbModel.UserColumn_Username) || slices.Contains(fields, dbModel.UserColumn_GivenName) || slices.Contains(fields, dbModel.UserColumn_FamilyName) {
		tx = tx.Joins("LEFT JOIN daylear_user ON user_access.recipient_user_id = daylear_user.user_id")
	}

	res := tx.Where("user_access.user_id = ? AND user_access.user_access_id = ?", parent.UserId.UserId, id.UserAccessId).
		First(&userAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("user access row not found")
			return model.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to get user access row")
		return model.UserAccess{}, ConvertGormError(res.Error)
	}

	return convert.UserAccessToCoreUserAccess(userAccess), nil
}

func (repo *Client) ListUserAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filterStr string, fields []string) ([]model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", parent.UserId.UserId).
		Str("filter", filterStr).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).
		Logger()

	fields = dbModel.UserAccessFieldMasker.Convert(fields)

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_access.user_access_id"},
		Desc:   true,
	}}

	var userAccesses []dbModel.UserAccess
	tx := repo.db.WithContext(ctx).
		Select(fields).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	if slices.Contains(fields, dbModel.UserColumn_Username) || slices.Contains(fields, dbModel.UserColumn_GivenName) || slices.Contains(fields, dbModel.UserColumn_FamilyName) {
		tx = tx.Joins("LEFT JOIN daylear_user ON user_access.recipient_user_id = daylear_user.user_id")
	}

	// Filter by user ID if provided
	if parent.UserId.UserId != 0 {
		tx = tx.Where("user_access.user_id = ?", parent.UserId.UserId)
	}

	res := tx.Find(&userAccesses)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to list user access rows")
		return nil, ConvertGormError(res.Error)
	}

	accesses := make([]model.UserAccess, len(userAccesses))
	for i, access := range userAccesses {
		accesses[i] = convert.UserAccessToCoreUserAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateUserAccess(ctx context.Context, access model.UserAccess, fields []string) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userAccessId", access.UserAccessId.UserAccessId).
		Strs("fields", fields).
		Logger()

	dbAccess := convert.CoreUserAccessToUserAccess(access)

	res := repo.db.WithContext(ctx).
		Select(dbModel.UpdateUserAccessFieldMasker.Convert(fields)).
		Clauses(&clause.Returning{}).
		Where("user_access_id = ?", access.UserAccessId.UserAccessId).
		Updates(&dbAccess)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update user access row")
		return model.UserAccess{}, ConvertGormError(res.Error)
	}

	return convert.UserAccessToCoreUserAccess(dbAccess), nil
}

func (repo *Client) FindStandardUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", id.UserId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	var userAccess dbModel.UserAccess
	res := repo.db.WithContext(ctx).
		Where("user_id = ? AND recipient_user_id = ?", id.UserId, authAccount.AuthUserId).
		First(&userAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("standard user user access not found")
			return model.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find standard user user access")
		return model.UserAccess{}, ConvertGormError(res.Error)
	}
	return convert.UserAccessToCoreUserAccess(userAccess), nil
}

func (repo *Client) FindDelegatedCircleUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", id.UserId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.UserAccess
		dbModel.CircleAccess
	}
	var result Result
	res := repo.db.WithContext(ctx).
		Select("user_access.*, circle_access.*").
		Table("user_access").
		Joins("JOIN circle_access ON circle_access.circle_id = user_access.recipient_circle_id").
		Where("user_access.user_id = ? AND circle_access.recipient_user_id = ?", id.UserId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated circle user access not found")
			return model.UserAccess{}, model.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated circle user access")
		return model.UserAccess{}, model.CircleAccess{}, ConvertGormError(res.Error)
	}
	return convert.UserAccessToCoreUserAccess(result.UserAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", id.UserId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.UserAccess
		dbModel.UA
	}
	var result Result

	res := repo.db.WithContext(ctx).
		Select("user_access.*, ua.*").
		Table("user_access").
		Joins("JOIN user_access AS ua ON user_access.user_id = ua.recipient_user_id").
		Where("ua.user_id = ? AND user_access.recipient_user_id = ?", id.UserId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated user user access not found")
			return model.UserAccess{}, model.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated user user access")
		return model.UserAccess{}, model.UserAccess{}, ConvertGormError(res.Error)
	}
	return convert.UserAccessToCoreUserAccess(result.UserAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
