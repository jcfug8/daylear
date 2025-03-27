package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
)

// CreateRecipe creates a new recipe.
func (repo *Client) CreateRecipe(ctx context.Context, m cmodel.Recipe) (cmodel.Recipe, error) {
	errz := errz.Context("repository.create_recipe")

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("invalid recipe: %v", err)
	}

	recipeFields := masks.RemovePaths(
		gmodel.RecipeFields.Mask(),
	)

	err = repo.db.
		Select(recipeFields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ErrzError(errz, "", err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("unable to read recipe: %v", err)
	}

	return m, nil
}
