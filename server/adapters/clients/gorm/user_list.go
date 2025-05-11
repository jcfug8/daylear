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

// ListUsers lists users.
func (repo *Client) ListUsers(ctx context.Context, page *cmodel.PageToken[cmodel.User], filter string, fields []string) ([]cmodel.User, error) {
	args := make([]any, 0, 1)

	fields = masks.Map(fields, gmodel.UserMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("u.%s", field)
		}
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			cmodel.UserFields.GoogleId:   filtering.NewSQLField[string](gmodel.UserFields.GoogleId, "="),
			cmodel.UserFields.FacebookId: filtering.NewSQLField[string](gmodel.UserFields.FacebookId, "="),
			cmodel.UserFields.AmazonId:   filtering.NewSQLField[string](gmodel.UserFields.AmazonId, "="),
			cmodel.UserFields.Username:   filtering.NewSQLField[string](gmodel.UserFields.Username, "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	tx = tx.Clauses(filterClause)
	if len(args) > 0 {
		tx = tx.Where(gmodel.User{}, args...)
	}

	if page != nil {
		orders := []clause.OrderByColumn{{
			Column: clause.Column{Name: "user_id"},
			Desc:   true,
		}}

		tx = tx.Order(clause.OrderBy{Columns: orders}).
			Limit(int(page.PageSize)).
			Offset(int(page.Skip))

		if page.Tail != nil {
			tail, err := convert.UserFromCoreModel(*page.Tail)
			if err != nil {
				return nil, fmt.Errorf("unable to read tail: %v", err)
			}

			tx = tx.Where(
				Seek(orders, gmodel.UserFields.Map(tail)))
		}
	}

	var mods []gmodel.User
	tableAlias := fmt.Sprintf("%s AS u", gmodel.User{}.TableName())
	err = tx.Table(tableAlias).Find(&mods).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.User, len(mods))
	for i, m := range mods {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read user: %v", err)
		}
	}

	return res, nil
}
