package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CircleMap maps the core model fields to the database model fields for the unified CircleAccess model.
var CircleMap = map[string]string{
	"permission": dbModel.CircleFields.Permission,
	"state":      dbModel.CircleFields.State,
	"title":      dbModel.CircleFields.Title,
	"visibility": dbModel.CircleFields.Visibility,
	"handle":     dbModel.CircleFields.Handle,
}

// CreateCircle creates a new circle.
func (repo *Client) CreateCircle(ctx context.Context, m cmodel.Circle) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM CreateCircle called")
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleFromCoreModel failed")
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	circleFields := masks.RemovePaths(
		gmodel.CircleFields.Mask(),
	)

	err = repo.db.
		Select(circleFields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Create failed")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleToCoreModel failed")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}
	log.Info().Msg("GORM CreateCircle returning successfully")
	return m, nil
}

// DeleteCircle deletes a circle.
func (repo *Client) DeleteCircle(ctx context.Context, id cmodel.CircleId) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM DeleteCircle called")
	gm := gmodel.Circle{CircleId: id.CircleId}
	err := repo.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Delete failed")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleToCoreModel failed")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}
	log.Info().Msg("GORM DeleteCircle returning successfully")
	return m, nil
}

// GetCircle gets a circle.
func (repo *Client) GetCircle(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM GetCircle called")
	gm := gmodel.Circle{}

	tx := repo.db.WithContext(ctx).
		Select("circle.*", "circle_access.permission_level", "circle_access.state", "circle_access.circle_access_id").
		Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.AuthUserId).
		Where("circle.circle_id = ? AND (circle.visibility_level = ? OR circle_access.recipient_user_id = ?)", id.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.First failed")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleToCoreModel failed")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}
	log.Info().Msg("GORM GetCircle returning successfully")
	return m, nil
}

// UpdateCircle updates a circle.
func (repo *Client) UpdateCircle(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM UpdateCircle called")
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleFromCoreModel failed")
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	mask := masks.Map(fields, gmodel.CircleMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("convert.CircleToCoreModel failed")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}
	log.Info().Msg("GORM UpdateCircle returning successfully")
	return m, nil
}

// ListCircles lists circles.
func (repo *Client) ListCircles(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM ListCircles called")
	dbCircles := []gmodel.Circle{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "circle_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select("circle.*", "circle_access.permission_level", "circle_access.state", "circle_access.circle_access_id").
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.AuthUserId).
		Where("(circle_access.recipient_user_id = ? OR circle.visibility_level = ?)", authAccount.AuthUserId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC)

	conversion, err := repo.circleSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("circleSQLConverter.Convert failed")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbCircles).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Circle, len(dbCircles))
	for i, m := range dbCircles {
		res[i], err = convert.CircleToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("convert.CircleToCoreModel failed")
			return nil, fmt.Errorf("unable to read circle: %v", err)
		}
	}
	log.Info().Msg("GORM ListCircles returning successfully")
	return res, nil
}

// CircleHandleExists checks if a circle handle already exists.
func (repo *Client) CircleHandleExists(ctx context.Context, handle string) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&gmodel.Circle{}).Where("handle = ?", handle).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
