package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START recipeDeleteImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// DeleteRecipe deletes a recipe.
func (repo *Client) DeleteRecipe(ctx context.Context, m cmodel.Recipe) (cmodel.Recipe, error) {
	errz := errz.Context("repository.delete_recipe")

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("invalid recipe: %v", err)
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipeBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipe
	err = repo.db.WithContext(ctx).
		Select(gmodel.RecipeFields.Mask()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ErrzError(errz, "", err)
	}
	// IRIOMO:CUSTOM_CODE_SLOT_END

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("unable to read recipe: %v", err)
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipeAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return m, nil
}
