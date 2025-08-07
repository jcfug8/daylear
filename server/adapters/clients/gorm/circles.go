package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateCircle creates a new circle.
func (repo *Client) CreateCircle(ctx context.Context, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle when creating circle row")
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	err = repo.db.
		Select(fields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create circle row")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle row when creating circle")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// DeleteCircle deletes a circle.
func (repo *Client) DeleteCircle(ctx context.Context, id cmodel.CircleId) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", id.CircleId).
		Logger()

	gm := gmodel.Circle{CircleId: id.CircleId}
	err := repo.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete circle row")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle row when deleting circle")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// GetCircle gets a circle.
func (repo *Client) GetCircle(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId, fields []string) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", id.CircleId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.Circle{}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.CircleFieldMasker.Convert(fields)).
		Where("circle.circle_id = ?", id.CircleId)

	if authAccount.UserId != 0 {
		tx = tx.Joins("LEFT JOIN circle_access as ca ON circle.circle_id = ca.circle_id AND ca.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(circle.visibility_level = ? OR ca.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId)
	} else {
		tx = tx.Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(circle.visibility_level = ? OR circle_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get circle row")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle row when getting circle")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// UpdateCircle updates a circle.
func (repo *Client) UpdateCircle(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("circleId", m.Id.CircleId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle when updating circle row")
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.UpdateCircleFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update circle row")
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid circle row when updating circle")
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// ListCircles lists circles.
func (repo *Client) ListCircles(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Circle, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("filter", filter).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("offset", offset).
		Logger()

	dbCircles := []gmodel.Circle{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "circle_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.CircleFieldMasker.Convert(fields)).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset))

	if authAccount.UserId != 0 {
		tx = tx.Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN circle_access as ca ON circle.circle_id = ca.circle_id AND ca.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(circle_access.recipient_user_id = ? AND (circle.visibility_level = ? OR ca.state = ?))",
				authAccount.UserId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else {
		tx = tx.Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.AuthUserId).
			Where("(circle.visibility_level = ? OR circle_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	conversion, err := gmodel.CircleSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing circle rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbCircles).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list circle rows")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Circle, len(dbCircles))
	for i, m := range dbCircles {
		res[i], err = convert.CircleToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("invalid circle row when listing circles")
			return nil, fmt.Errorf("unable to read circle: %v", err)
		}
	}

	return res, nil
}

// CircleHandleExists checks if a circle handle already exists.
func (repo *Client) CircleHandleExists(ctx context.Context, handle string) (bool, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("handle", handle).
		Logger()

	var count int64
	err := repo.db.WithContext(ctx).Model(&gmodel.Circle{}).Where("handle = ?", handle).Count(&count).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to check if circle handle exists")
		return false, err
	}
	return count > 0, nil
}
