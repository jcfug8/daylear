package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateCircle creates a new circle.
func (repo *Client) CreateCircle(ctx context.Context, m cmodel.Circle) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
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
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// DeleteCircle deletes a circle.
func (repo *Client) DeleteCircle(ctx context.Context, id cmodel.CircleId) (cmodel.Circle, error) {
	gm := gmodel.Circle{CircleId: id.CircleId}
	err := repo.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// GetCircle gets a circle.
func (repo *Client) GetCircle(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.CircleId) (cmodel.Circle, error) {
	gm := gmodel.Circle{}

	tx := repo.db.WithContext(ctx).
		Select("circle.*", "circle_access.permission_level").
		Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id AND circle_access.recipient_user_id = ?", authAccount.UserId).
		Where("circle.circle_id = ? AND (circle.visibility_level = ? OR circle_access.recipient_user_id = ?)", id.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId)

	err := tx.First(&gm).Error
	if err != nil {
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err := convert.CircleToCoreModel(gm)
	if err != nil {
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// UpdateCircle updates a circle.
func (repo *Client) UpdateCircle(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	mask := masks.Map(fields, gmodel.CircleMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// ListCircles lists circles.
func (repo *Client) ListCircles(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.Circle, error) {
	dbCircles := []gmodel.Circle{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "circle_id"},
		Desc:   true,
	}}

	err := repo.db.WithContext(ctx).
		Select("circle.*, circle_access.permission_level").
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(offset)).
		Joins("LEFT JOIN circle_access ON circle.circle_id = circle_access.circle_id").
		Where("(circle_access.recipient_user_id = ? OR circle.visibility_level = ?)", authAccount.UserId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC).
		Find(&dbCircles).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Circle, len(dbCircles))
	for i, m := range dbCircles {
		res[i], err = convert.CircleToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read circle: %v", err)
		}
	}

	return res, nil
}
