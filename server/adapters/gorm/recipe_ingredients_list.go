package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
)

// ListRecipeIngredients lists recipes.
func (repo *Client) ListRecipeIngredients(ctx context.Context, page *cmodel.PageToken[cmodel.RecipeIngredient], filter string, fields []string) ([]cmodel.RecipeIngredient, error) {
	errz := errz.Context("repository.list_recipes")

	fields = masks.Map(fields, gmodel.RecipeMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("recipe_id", "="),
			"user_id":   filtering.NewSQLField[int64]("user_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid filter: %v", err)
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	if page != nil {
		orders := []clause.OrderByColumn{{
			Column: clause.Column{Name: "recipe_id"},
			Desc:   true,
		}}

		tx = tx.Order(clause.OrderBy{Columns: orders}).
			Limit(int(page.PageSize)).
			Offset(int(page.Skip))

		if page.Tail != nil {
			tail := convert.RecipeIngredientFromCoreModel(*page.Tail)

			tx = tx.Where(
				Seek(orders, gmodel.RecipeIngredientFields.Map(tail)))
		}
	}

	var dbRecipeIngredients []gmodel.RecipeIngredient
	if err = tx.Find(&dbRecipeIngredients).Error; err != nil {
		return nil, ErrzError(errz, "", err)
	}

	recipeIngredients := convert.RecipeIngredientListToCoreModel(dbRecipeIngredients)

	ingredientIds := []int64{}
	for _, recipeIngredient := range recipeIngredients {
		ingredientIds = append(ingredientIds, recipeIngredient.IngredientId)
	}

	var dbIngredients []gmodel.Ingredient
	err = repo.db.WithContext(ctx).
		Where("ingredient_id IN ?", ingredientIds).
		Find(&dbIngredients).Error
	if err != nil {
		return nil, errz.Wrapf("unable to list ingredients: %v", err)
	}

	ingredients := convert.IngredientListToCoreModel(dbIngredients)

	for i, recipeIngredient := range recipeIngredients {
		for _, ingredient := range ingredients {
			if recipeIngredient.IngredientId == ingredient.IngredientId {
				recipeIngredients[i].Ingredient = ingredient
				break
			}
		}
	}

	return recipeIngredients, nil
}
