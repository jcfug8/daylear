package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
)

func (repo *Client) BulkCreateIngredients(ctx context.Context, ingredients []cmodel.Ingredient) ([]cmodel.Ingredient, error) {
	if len(ingredients) == 0 {
		return ingredients, nil
	}

	dbIngredients := convert.IngredientListFromCoreModel(ingredients)

	ingredientFields := masks.RemovePaths(
		gmodel.IngredientFields.Mask(),
	)

	err := repo.db.WithContext(ctx).
		Select(ingredientFields).
		Clauses(clause.Returning{}).
		Create(&dbIngredients).Error
	if err != nil {
		return nil, err
	}

	ingredients = convert.IngredientListToCoreModel(dbIngredients)

	return ingredients, nil
}
