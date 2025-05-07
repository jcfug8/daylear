package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START recipeDeleteImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// DeleteRecipe deletes a recipe.
func (repo *Client) DeleteRecipe(ctx context.Context, m cmodel.Recipe) (cmodel.Recipe, error) {
	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid recipe: %v", err)}
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipeBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipe
	err = repo.db.WithContext(ctx).
		Select(gmodel.RecipeFields.Mask()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}
	// IRIOMO:CUSTOM_CODE_SLOT_END

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START deleteRecipeAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return m, nil
}
