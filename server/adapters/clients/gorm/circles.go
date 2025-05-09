package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateCircle creates a new circle.
func (repo *Client) CreateCircle(ctx context.Context, m cmodel.Circle) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	circleFields := []string{
		gmodel.CircleFields.CircleId,
		gmodel.CircleFields.Title,
	}

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
func (repo *Client) DeleteCircle(ctx context.Context, m cmodel.Circle) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
	if err != nil {
		return cmodel.Circle{}, fmt.Errorf("unable to read circle: %v", err)
	}

	return m, nil
}

// GetCircle gets a circle.
func (repo *Client) GetCircle(ctx context.Context, m cmodel.Circle, fields []string) (cmodel.Circle, error) {
	gm, err := convert.CircleFromCoreModel(m)
	if err != nil {
		return cmodel.Circle{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid circle: %v", err)}
	}

	mask := fields
	if len(mask) == 0 {
		mask = gmodel.CircleFields.Mask()
	}

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		First(&gm).Error
	if err != nil {
		return cmodel.Circle{}, ConvertGormError(err)
	}

	m, err = convert.CircleToCoreModel(gm)
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
func (repo *Client) ListCircles(ctx context.Context, page *cmodel.PageToken[cmodel.Circle], filter string, fields []string) ([]cmodel.Circle, error) {
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

	if page != nil {
		orders := []clause.OrderByColumn{{
			Column: clause.Column{Name: "circle_id"},
			Desc:   true,
		}}

		tx = tx.Order(clause.OrderBy{Columns: orders}).
			Limit(int(page.PageSize)).
			Offset(int(page.Skip))

		if page.Tail != nil {
			tail, err := convert.CircleFromCoreModel(*page.Tail)
			if err != nil {
				return nil, fmt.Errorf("unable to read tail: %v", err)
			}

			tx = tx.Where(
				Seek(orders, gmodel.CircleFields.Map(tail)))
		}
	}

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

// BulkCreateCircleUsers creates many circle users at once.
func (repo *Client) BulkCreateCircleUsers(ctx context.Context, circleId cmodel.CircleId, userIds []int64, permission permPb.PermissionLevel) error {
	if len(userIds) == 0 {
		return nil
	}

	circleUsers := []gmodel.CircleUser{}

	for _, userId := range userIds {
		circleUsers = append(circleUsers, gmodel.CircleUser{
			CircleId:        circleId.CircleId,
			UserId:          userId,
			PermissionLevel: permission,
		})
	}

	return repo.db.WithContext(ctx).Create(&circleUsers).Error
}

// BulkDeleteCircleUsers deletes circle users by filter.
func (repo *Client) BulkDeleteCircleUsers(ctx context.Context, filter string) error {
	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"circle_id": filtering.NewSQLField[int64]("circle_id", "="),
			"user_id":   filtering.NewSQLField[int64]("user_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbCircleUsers []gmodel.CircleUser
	if err = tx.Delete(&dbCircleUsers).Error; err != nil {
		return repository.ErrInternal{Msg: fmt.Sprintf("failed to delete circle users: %v", err)}
	}

	return nil
}

// GetCircleUserPermission gets a user's permission for a circle.
func (repo *Client) GetCircleUserPermission(ctx context.Context, userId int64, circleId int64) (permPb.PermissionLevel, error) {
	var circleUser gmodel.CircleUser
	err := repo.db.WithContext(ctx).
		Where("user_id = ? AND circle_id = ?", userId, circleId).
		First(&circleUser).Error
	if err != nil {
		return permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED, ConvertGormError(err)
	}

	return circleUser.PermissionLevel, nil
}
