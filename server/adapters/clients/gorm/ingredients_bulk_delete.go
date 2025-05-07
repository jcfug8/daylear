package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// BulkDeleteIngredients deletes ingredients in bulk.
func (repo *Client) BulkDeleteIngredients(ctx context.Context, filter string) ([]cmodel.Ingredient,error) {
	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"ingredient_id": filtering.NewSQLField[int64]("ingredient_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbIngredients []gmodel.Ingredient
	if err = tx.Clauses(clause.Returning{}).Delete(&dbIngredients).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	ingredients := convert.IngredientListToCoreModel(dbIngredients)

	return ingredients, nil
}
