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

// BulkDeleteRecipeIngredients deletes recipe ingredients in bulk.
func (repo *Client) BulkDeleteRecipeIngredients(ctx context.Context, filter string) ([]cmodel.RecipeIngredient, error) {
	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("recipe_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	var dbRecipeIngredients []gmodel.RecipeIngredient
	if err = tx.Clauses(clause.Returning{}).Delete(&dbRecipeIngredients).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	recipeIngredients := convert.RecipeIngredientListToCoreModel(dbRecipeIngredients)

	ingredientIds := []int64{}
	for _, recipeIngredient := range recipeIngredients {
		ingredientIds = append(ingredientIds, recipeIngredient.IngredientId)
	}

	var dbIngredients []gmodel.Ingredient
	err = repo.db.WithContext(ctx).
		Where("ingredient_id IN ?", ingredientIds).
		Clauses(clause.Returning{}).
		Delete(&dbIngredients).Error
	if err != nil {
		return nil, fmt.Errorf("unable to list ingredients: %v", err)
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
