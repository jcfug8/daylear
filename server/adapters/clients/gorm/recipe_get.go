package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// GetRecipe gets a recipe.
func (repo *Client) GetRecipe(ctx context.Context, m cmodel.Recipe, fields []string) (cmodel.Recipe, error) {
	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid recipe: %v", err)}
	}

	mask := masks.Map(fields, gmodel.RecipeMap)
	if len(mask) == 0 {
		mask = gmodel.RecipeFields.Mask()
	}

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		First(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}
