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

// ListRecipes lists recipes.
func (repo *Client) ListRecipes(ctx context.Context, page *cmodel.PageToken[cmodel.Recipe], parent cmodel.RecipeParent, filter string, fields []string) ([]cmodel.Recipe, error) {
	queryModel := gmodel.Recipe{}

	args := make([]any, 0, 1)

	fields = masks.Map(fields, gmodel.RecipeMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("r.%s", field)
		}
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("r.recipe_id", "="),
		})

	filterClause, _, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if parent.CircleId != 0 {
		tx.Joins("JOIN recipe_circle ON recipe_circle.recipe_id = r.recipe_id AND recipe_circle.circle_id = ?", parent.CircleId)
	} else if parent.UserId != 0 {
		tx.Joins("JOIN recipe_user ON recipe_user.recipe_id = r.recipe_id AND recipe_user.user_id = ?", parent.UserId)
	}

	if parent.UserId != 0 && parent.CircleId != 0 {
		tx.Joins("JOIN circle_user ON circle_user.circle_id = recipe_circle.circle_id AND circle_user.user_id = ?", parent.UserId)
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}
	if len(args) > 0 {
		tx = tx.Where(queryModel, args...)
	}

	if page != nil {
		orders := []clause.OrderByColumn{{
			Column: clause.Column{Name: "r.recipe_id"},
			Desc:   true,
		}}

		tx = tx.Order(clause.OrderBy{Columns: orders}).
			Limit(int(page.PageSize)).
			Offset(int(page.Skip))

		if page.Tail != nil {
			tail, err := convert.RecipeFromCoreModel(*page.Tail)
			if err != nil {
				return nil, fmt.Errorf("unable to read tail: %v", err)
			}

			tx = tx.Where(
				Seek(orders, gmodel.RecipeFields.Map(tail)))
		}
	}

	var mods []gmodel.Recipe
	tableAlias := fmt.Sprintf("%s AS r", gmodel.Recipe{}.TableName())
	if err = tx.Table(tableAlias).Find(&mods).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.Recipe, len(mods))
	for i, m := range mods {
		res[i], err = convert.RecipeToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read recipe: %v", err)
		}
	}

	return res, nil
}
