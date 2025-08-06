package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CircleAccessMap maps the core model fields to the database model fields for the unified CircleAccess model.
var CircleAccessMap = map[string]string{
	model.CircleAccessFields.Level:         dbModel.CircleAccessFields.PermissionLevel,
	model.CircleAccessFields.State:         dbModel.CircleAccessFields.State,
	model.CircleAccessFields.RecipientUser: dbModel.CircleAccessFields.RecipientUserId,
}

// UpdateCircleAccessMap maps the updatable core model fields to the database model fields for the CircleAccess model.
var UpdateCircleAccessMap = map[string][]string{
	model.CircleAccessFields.Level: []string{
		dbModel.CircleAccessFields.PermissionLevel,
	},
	model.CircleAccessFields.State: []string{
		dbModel.CircleAccessFields.State,
	},
}

var UpdateCircleAccessFieldMasker = fieldmask.NewFieldMasker(UpdateCircleAccessMap)

func (repo *Client) CreateCircleAccess(ctx context.Context, access model.CircleAccess) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM CreateCircleAccess called")
	db := repo.db.WithContext(ctx)

	// Validate that exactly one recipient type is set
	if access.Recipient.UserId == 0 {
		log.Error().Msg("recipient is required")
		return model.CircleAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	circleAccess := convert.CoreCircleAccessToCircleAccess(access)
	res := db.Create(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("duplicate key error on create")
			return model.CircleAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("db.Create failed")
		return model.CircleAccess{}, res.Error
	}

	access.CircleAccessId.CircleAccessId = circleAccess.CircleAccessId
	log.Info().Msg("GORM CreateCircleAccess returning successfully")
	return access, nil
}

func (repo *Client) DeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM DeleteCircleAccess called")
	db := repo.db.WithContext(ctx)

	res := db.Delete(&dbModel.CircleAccess{}, id.CircleAccessId)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("db.Delete failed")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Error().Msg("no rows affected on delete")
		return repository.ErrNotFound{}
	}

	log.Info().Msg("GORM DeleteCircleAccess returning successfully")
	return nil
}

func (repo *Client) BulkDeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM BulkDeleteCircleAccess called")
	db := repo.db.WithContext(ctx)

	res := db.Where("circle_id = ?", parent.CircleId.CircleId).Delete(&dbModel.CircleAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("db.Delete failed")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Error().Msg("no rows affected on bulk delete")
		return repository.ErrNotFound{}
	}

	log.Info().Msg("GORM BulkDeleteCircleAccess returning successfully")
	return nil
}

func (repo *Client) GetCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM GetCircleAccess called")
	db := repo.db.WithContext(ctx)

	var circleAccess dbModel.CircleAccess
	res := db.Table("circle_access").
		Select("circle_access.*, daylear_user.username as recipient_username, daylear_user.given_name as recipient_given_name, daylear_user.family_name as recipient_family_name").
		Joins("LEFT JOIN daylear_user ON circle_access.recipient_user_id = daylear_user.user_id").
		Where("circle_access.circle_id = ? AND circle_access.circle_access_id = ?", parent.CircleId.CircleId, id.CircleAccessId).
		First(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Error().Err(res.Error).Msg("record not found on get")
			return model.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("db.First failed")
		return model.CircleAccess{}, res.Error
	}

	log.Info().Msg("GORM GetCircleAccess returning successfully")
	return convert.CircleAccessToCoreCircleAccess(circleAccess), nil
}

func (repo *Client) ListCircleAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM ListCircleAccesses called")
	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required")
		return nil, repository.ErrInvalidArgument{Msg: "user id is required"}
	}

	conversion, err := repo.circleAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("circleAccessSQLConverter.Convert failed")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	var circleAccesses []dbModel.CircleAccess
	db := repo.db.WithContext(ctx).Table("circle_access").
		Select("circle_access.*, daylear_user.username as recipient_username, daylear_user.given_name as recipient_given_name, daylear_user.family_name as recipient_family_name").
		Joins("LEFT JOIN daylear_user ON circle_access.recipient_user_id = daylear_user.user_id")

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
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.CircleAccess, len(circleAccesses))
	for i, access := range circleAccesses {
		accesses[i] = convert.CircleAccessToCoreCircleAccess(access)
	}

	log.Info().Msg("GORM ListCircleAccesses returning successfully")
	return accesses, nil
}

func (repo *Client) UpdateCircleAccess(ctx context.Context, access model.CircleAccess, updateMask []string) (model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM UpdateCircleAccess called")
	dbAccess := convert.CoreCircleAccessToCircleAccess(access)

	columns := UpdateCircleAccessFieldMasker.Convert(updateMask)

	db := repo.db.WithContext(ctx).Select(columns).Clauses(clause.Returning{})

	err := db.Where("circle_access_id = ?", access.CircleAccessId.CircleAccessId).Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return model.CircleAccess{}, ConvertGormError(err)
	}

	log.Info().Msg("GORM UpdateCircleAccess returning successfully")
	return convert.CircleAccessToCoreCircleAccess(dbAccess), nil
}

func (repo *Client) FindStandardUserCircleAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from circle_access where recipient_user_id = ? and circle_id = ?

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

func (repo *Client) FindDelegatedCircleCircleAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.CircleAccess, cmodel.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from circle_access
	// 	JOIN circle_access ON circle_access.circle_id = circle_access.recipient_circle_id
	// WHERE circle_access.circle_id = 1 AND circle_access.recipient_user_id = 1 LIMIT 1;

	type Result struct {
		dbModel.CircleAccess
		dbModel.CA
	}
	var result Result
	res := repo.db.WithContext(ctx).
		Select("circle_access.*, ca.*").
		Table("circle_access").
		Joins("JOIN circle_access AS ca ON circle_access.circle_id = ca.recipient_circle_id").
		Where("ca.circle_id = ? AND circle_access.recipient_user_id = ?", id.CircleId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated circle circle access not found")
			return cmodel.CircleAccess{}, cmodel.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated circle circle access")
		return cmodel.CircleAccess{}, cmodel.CircleAccess{}, res.Error
	}
	return convert.CircleAccessToCoreCircleAccess(result.CircleAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserCircleAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.CircleAccess, cmodel.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	// SELECT * from circle_access
	// 	JOIN user_access ON user_access.user_id = circle_access.recipient_user_id
	// WHERE circle_access.circle_id = ? AND user_access.recipient_user_id = ? LIMIT 1;

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
