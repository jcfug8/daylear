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

// ListRecipeIngredients lists recipe ingredients.
func (repo *Client) ListRecipeIngredients(ctx context.Context, page *cmodel.PageToken[cmodel.RecipeIngredient], filter string, fields []string) ([]cmodel.RecipeIngredient, error) {
	tx := repo.db.WithContext(ctx)

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			"recipe_id": filtering.NewSQLField[int64]("recipe_ingredient.recipe_id", "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	if filterClause != nil {
		tx = tx.Clauses(filterClause)
	}

	// Start from ingredient table and join with recipe_ingredient and recipe tables
	tx = tx.Table("ingredient").
		Joins("JOIN recipe_ingredient ON recipe_ingredient.ingredient_id = ingredient.ingredient_id").
		Select("ingredient.title as ingredient_title, recipe_ingredient.*")

	var dbRecipeIngredients []gmodel.RecipeIngredient
	err = tx.Find(&dbRecipeIngredients).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	recipeIngredients := convert.RecipeIngredientListToCoreModel(dbRecipeIngredients)

	return recipeIngredients, nil
}
