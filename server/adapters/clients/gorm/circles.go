package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
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
func (repo *Client) GetCircle(ctx context.Context, id cmodel.CircleId, fields []string) (cmodel.Circle, error) {
	gm := gmodel.Circle{CircleId: id.CircleId}
	err := repo.db.WithContext(ctx).
		Select(masks.Map(fields, gmodel.CircleMap)).
		First(&gm).Error
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	mask := masks.Map(fields, gmodel.CircleMap)
	if len(mask) == 0 {
		mask = gmodel.CircleFields.Mask()
	}

	tx := repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{})

	err = tx.First(&gm).Error
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
func (repo *Client) UpdateCircle(ctx context.Context, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	updateFields := masks.Map(fields, gmodel.CircleMap)
	if len(updateFields) == 0 {
		updateFields = gmodel.CircleFields.Mask()
	}

	err = repo.db.WithContext(ctx).
		Model(&gm).
		Select(updateFields).
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
	queryModel := gmodel.Circle{}

	args := make([]any, 0, 1)

	fields = masks.Map(fields, gmodel.CircleMap)

	tx := repo.db.WithContext(ctx)

	if len(fields) > 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("c.%s", field)
		}
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"circle_id": filtering.NewSQLField[int64]("c.circle_id", "="),
			"user_id":   filtering.NewSQLField[int64]("circle_user.user_id", "="),
			"title":     filtering.NewSQLField[string]("c.title", "="),
			"is_public": filtering.NewSQLField[bool]("c.is_public", "="),
		})

	filterClause, info, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if info.HasField("user_id") {
		tx.Joins("JOIN circle_user ON circle_user.circle_id = c.circle_id")
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}
	if len(args) > 0 {
		tx = tx.Where(queryModel, args...)
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "circle_id"},
		Desc:   true,
	}}

	tx = tx.Order(clause.OrderBy{Columns: orders})

	tx = tx.Limit(int(pageSize)).
		Offset(int(offset))

	var mods []gmodel.Circle
	tableAlias := fmt.Sprintf("%s AS c", gmodel.Circle{}.TableName())
	if err = tx.Table(tableAlias).Find(&mods).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Circle, len(mods))
	for i, m := range mods {
		res[i], err = convert.CircleToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read circle: %v", err)
		}
	}

	return res, nil
}
