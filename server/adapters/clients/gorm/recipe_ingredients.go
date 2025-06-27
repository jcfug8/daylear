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

func (repo *Client) SetRecipeIngredients(ctx context.Context, recipeId cmodel.RecipeId, ingredientGroups []cmodel.IngredientGroup) error {
	ingredientCount := 0
	for _, group := range ingredientGroups {
		ingredientCount += len(group.RecipeIngredients)
	}
	if ingredientCount == 0 {
		return nil
	}

	dbRecipeIngredients := convert.IngredientGroupsFromCoreModel(recipeId, ingredientGroups)

	txRepo, err := repo.beginTransaction()
	if err != nil {
		return err
	}
	defer txRepo.Rollback()

	err = txRepo.db.WithContext(ctx).
		Where(gmodel.RecipeIngredient{RecipeId: recipeId.RecipeId}).
		Delete(&gmodel.RecipeIngredient{}).Error
	if err != nil {
		return err
	}

	err = txRepo.db.WithContext(ctx).
		Create(&dbRecipeIngredients).Error
	if err != nil {
		return err
	}

	err = txRepo.Commit()
	if err != nil {
		return err
	}

	return nil
}

// ListRecipeIngredients lists recipe ingredients.
func (repo *Client) ListRecipeIngredients(ctx context.Context, pageSize int32, pageOffset int64, filter string, fields []string) ([]cmodel.RecipeIngredient, error) {
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
