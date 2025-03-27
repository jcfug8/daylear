package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
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
