package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
)

// ListIngredients lists recipes.
func (repo *Client) BulkDeleteIngredients(ctx context.Context, filter string) ([]cmodel.Ingredient, error) {
	errz := errz.Context("repository.list_recipes")

	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"ingredient_id": filtering.NewSQLField[int64]("ingredient_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid filter: %v", err)
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbIngredients []gmodel.Ingredient
	if err = tx.Clauses(clause.Returning{}).Delete(&dbIngredients).Error; err != nil {
		return nil, ErrzError(errz, "", err)
	}

	ingredients := convert.IngredientListToCoreModel(dbIngredients)

	return ingredients, nil
}
