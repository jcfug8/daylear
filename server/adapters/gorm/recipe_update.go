package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START recipeUpdateImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// UpdateRecipe updates a recipe.
func (repo *Client) UpdateRecipe(ctx context.Context, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	errz := errz.Context("repository.update_recipes")

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, ErrzError(errz, "error reading recipe: %v", err)
	}

	mask := masks.Map(fields, gmodel.RecipeMap)
	for _, path := range mask {
		switch path {
		// IRIOMO:CUSTOM_CODE_SLOT_START updateRecipeMask
		// IRIOMO:CUSTOM_CODE_SLOT_END
		}
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START updateRecipeBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ErrzError(errz, "", err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("unable to read recipe: %v", err)
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START updateRecipeAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return m, nil
}
