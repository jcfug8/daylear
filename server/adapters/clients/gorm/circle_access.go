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
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Client) CreateCircleAccess(ctx context.Context, access model.CircleAccess, fields []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", access.CircleId.CircleId).
		Int64("recipientUserId", access.Recipient.UserId).
		Strs("fields", fields).
		Logger()

	// Validate that exactly one recipient type is set
	if access.Recipient.UserId == 0 {
		log.Error().Msg("recipient is required to create circle access row")
		return model.CircleAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	circleAccess := convert.CoreCircleAccessToCircleAccess(access)
	res := repo.db.WithContext(ctx).
		Select(dbModel.CircleAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("unable to create circle access row: already exists")
			return model.CircleAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("unable to create circle access row")
		return model.CircleAccess{}, res.Error
	}

	access.CircleAccessId.CircleAccessId = circleAccess.CircleAccessId
	return access, nil
}

func (repo *Client) DeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", parent.CircleId.CircleId).
		Int64("circleAccessId", id.CircleAccessId).
		Logger()

	if parent.CircleId.CircleId == 0 {
		log.Error().Msg("circle id is required to delete circle access row")
		return repository.ErrInvalidArgument{Msg: "circle id is required"}
	}

	if id.CircleAccessId == 0 {
		log.Error().Msg("circle access id is required to delete circle access row")
		return repository.ErrInvalidArgument{Msg: "circle access id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Delete(&dbModel.CircleAccess{}, id.CircleAccessId)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to delete circle access row")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no circle access row deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", parent.CircleId.CircleId).
		Logger()

	if parent.CircleId.CircleId == 0 {
		log.Error().Msg("circle id is required to bulk delete circle access rows")
		return repository.ErrInvalidArgument{Msg: "circle id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("circle_id = ?", parent.CircleId.CircleId).Delete(&dbModel.CircleAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to bulk delete circle access rows")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no circle access rows deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId, fields []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", parent.CircleId.CircleId).
		Int64("circleAccessId", id.CircleAccessId).
		Strs("fields", fields).
		Logger()

	fields = dbModel.CircleAccessFieldMasker.Convert(fields)

	var circleAccess dbModel.CircleAccess
	tx := repo.db.WithContext(ctx).Select(fields)
	if slices.Contains(fields, dbModel.UserColumn_Username) || slices.Contains(fields, dbModel.UserColumn_GivenName) || slices.Contains(fields, dbModel.UserColumn_FamilyName) {
		tx = tx.Joins("LEFT JOIN daylear_user ON circle_access.recipient_user_id = daylear_user.user_id")
	}

	res := tx.Where("circle_access.circle_id = ? AND circle_access.circle_access_id = ?", parent.CircleId.CircleId, id.CircleAccessId).
		First(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("circle access row not found")
			return model.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to get circle access row")
		return model.CircleAccess{}, res.Error
	}

	return convert.CircleAccessToCoreCircleAccess(circleAccess), nil
}

func (repo *Client) ListCircleAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filterStr string, fields []string) ([]model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", parent.CircleId.CircleId).
		Str("filter", filterStr).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).
		Logger()

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required to list circle access rows")
		return nil, repository.ErrInvalidArgument{Msg: "user id is required"}
	}

	fields = dbModel.CircleAccessFieldMasker.Convert(fields)

	conversion, err := dbModel.CircleAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing circle access rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	var circleAccesses []dbModel.CircleAccess
	db := repo.db.WithContext(ctx).Select(fields)

	if slices.Contains(fields, dbModel.UserColumn_Username) || slices.Contains(fields, dbModel.UserColumn_GivenName) || slices.Contains(fields, dbModel.UserColumn_FamilyName) {
		db = db.Joins("LEFT JOIN daylear_user ON circle_access.recipient_user_id = daylear_user.user_id")
	}

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	// Filter by circle ID if provided
	if parent.CircleId.CircleId != 0 {
		db = db.Where("circle_access.circle_id = ?", parent.CircleId.CircleId)
	}

	db = db.Where(
		"circle_access.recipient_user_id = ? OR circle_access.circle_id IN (SELECT circle_id FROM circle_access WHERE recipient_user_id = ? AND permission_level >= ?)",
		authAccount.AuthUserId, authAccount.AuthUserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
	)

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&circleAccesses).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list circle access rows")
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.CircleAccess, len(circleAccesses))
	for i, access := range circleAccesses {
		accesses[i] = convert.CircleAccessToCoreCircleAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateCircleAccess(ctx context.Context, access model.CircleAccess, fields []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleAccessId", access.CircleAccessId.CircleAccessId).
		Strs("fields", fields).
		Logger()

	dbAccess := convert.CoreCircleAccessToCircleAccess(access)

	db := repo.db.WithContext(ctx).
		Select(dbModel.UpdateCircleAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{})

	err := db.Where("circle_access_id = ?", access.CircleAccessId.CircleAccessId).Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle access row")
		return model.CircleAccess{}, ConvertGormError(err)
	}

	return convert.CircleAccessToCoreCircleAccess(dbAccess), nil
}

func (repo *Client) FindStandardUserCircleAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", id.CircleId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	var circleAccess dbModel.CircleAccess
	res := repo.db.WithContext(ctx).
		Where("circle_id = ? AND recipient_user_id = ?", id.CircleId, authAccount.AuthUserId).
		First(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("standard user circle access not found")
			return cmodel.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find standard user circle access")
		return cmodel.CircleAccess{}, res.Error
	}
	return convert.CircleAccessToCoreCircleAccess(circleAccess), nil
}

func (repo *Client) FindDelegatedUserCircleAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.CircleAccess, cmodel.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", id.CircleId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.CircleAccess
		dbModel.UserAccess
	}
	var result Result

	res := repo.db.WithContext(ctx).
		Select("circle_access.*, user_access.*").
		Table("circle_access").
		Joins("JOIN user_access ON user_access.user_id = circle_access.recipient_user_id").
		Where("circle_access.circle_id = ? AND user_access.recipient_user_id = ?", id.CircleId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated user circle access not found")
			return cmodel.CircleAccess{}, cmodel.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated user circle access")
		return cmodel.CircleAccess{}, cmodel.UserAccess{}, res.Error
	}
	return convert.CircleAccessToCoreCircleAccess(result.CircleAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
