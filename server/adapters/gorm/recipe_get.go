package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// GetRecipe gets a recipe.
func (repo *Client) GetRecipe(ctx context.Context, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	errz := errz.Context("repository.get_recipe")

	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("invalid recipe: %v", err)
	}

	err = repo.db.WithContext(ctx).
		Select(masks.Map(fields, gmodel.RecipeMap)).
		Take(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ErrzError(errz, "", err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, errz.Wrapf("unable to read recipe: %v", err)
	}

	return m, nil
}
