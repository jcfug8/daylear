package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserAccessMap maps the core model fields to the database model fields for the unified UserAccess model.
var UserAccessMap = map[string]string{
	model.UserAccessFields.Level:         dbModel.UserAccessFields.PermissionLevel,
	model.UserAccessFields.State:         dbModel.UserAccessFields.State,
	model.UserAccessFields.RecipientUser: dbModel.UserAccessFields.RecipientUserId,
}

func (repo *Client) CreateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	db := repo.db.WithContext(ctx)

	userAccess := convert.CoreUserAccessToUserAccess(access)
	res := db.Create(&userAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("duplicate key error")
			return model.UserAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("db.Create failed")
		return model.UserAccess{}, res.Error
	}

	access.UserAccessId.UserAccessId = userAccess.UserAccessId
	return access, nil
}

func (repo *Client) DeleteUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	db := repo.db.WithContext(ctx)

	res := db.Where("user_id = ? AND user_access_id = ?", parent.UserId.UserId, id.UserAccessId).Delete(&dbModel.UserAccess{})
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

func (repo *Client) GetUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	db := repo.db.WithContext(ctx)

	var userAccess dbModel.UserAccess
	res := db.Table("user_access").
		Select("user_access.*, daylear_user.username as recipient_username").
		Joins("LEFT JOIN daylear_user ON user_access.recipient_user_id = daylear_user.user_id").
		Where("user_access.user_id = ? AND user_access.user_access_id = ?", parent.UserId.UserId, id.UserAccessId).
		First(&userAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("record not found")
			return model.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("db.First failed")
		return model.UserAccess{}, res.Error
	}

	return convert.UserAccessToCoreUserAccess(userAccess), nil
}

func (repo *Client) ListUserAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_access.user_access_id"},
		Desc:   true,
	}}

	var userAccesses []dbModel.UserAccess
	db := repo.db.WithContext(ctx).
		Table("user_access").
		Select("user_access.*, daylear_user.username as recipient_username").
		Joins("LEFT JOIN daylear_user ON user_access.recipient_user_id = daylear_user.user_id").
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	// Filter by user ID if provided
	if parent.UserId.UserId != 0 {
		db = db.Where("user_access.user_id = ?", parent.UserId.UserId)
	}

	err := db.Find(&userAccesses).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.UserAccess, len(userAccesses))
	for i, access := range userAccesses {
		accesses[i] = convert.UserAccessToCoreUserAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	dbAccess := convert.CoreUserAccessToUserAccess(access)
	db := repo.db.WithContext(ctx).Clauses(&clause.Returning{})

	err := db.Model(&dbModel.UserAccess{}).
		Where("user_access_id = ?", access.UserAccessId.UserAccessId).
		Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return model.UserAccess{}, ConvertGormError(err)
	}

	return convert.UserAccessToCoreUserAccess(dbAccess), nil
}

func (repo *Client) FindStandardUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from user_access where recipient_user_id = ? and user_id = ?

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
		return model.UserAccess{}, res.Error
	}
	return convert.UserAccessToCoreUserAccess(userAccess), nil
}

func (repo *Client) FindDelegatedCircleUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from user_access
	// 	JOIN circle_access ON circle_access.circle_id = user_access.recipient_circle_id
	// WHERE user_access.user_id = 1 AND circle_access.recipient_user_id = 1 LIMIT 1;

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
		return model.UserAccess{}, model.CircleAccess{}, res.Error
	}
	return convert.UserAccessToCoreUserAccess(result.UserAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId) (model.UserAccess, model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from user_access
	// 	JOIN user_access ON user_access.user_id = user_access.recipient_user_id
	// WHERE user_access.user_id = ? AND user_access.recipient_user_id = ? LIMIT 1;

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
		return model.UserAccess{}, model.UserAccess{}, res.Error
	}
	return convert.UserAccessToCoreUserAccess(result.UserAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
