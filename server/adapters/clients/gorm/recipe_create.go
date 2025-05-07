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

// CreateRecipe creates a new recipe.
func (repo *Client) CreateRecipe(ctx context.Context, m cmodel.Recipe) (cmodel.Recipe, error) {
	gm, err := convert.RecipeFromCoreModel(m)
	if err != nil {
		return cmodel.Recipe{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid recipe: %v", err)}
	}

	recipeFields := masks.RemovePaths(
		gmodel.RecipeFields.Mask(),
	)

	err = repo.db.
		Select(recipeFields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.Recipe{}, ConvertGormError(err)
	}

	m, err = convert.RecipeToCoreModel(gm)
	if err != nil {
		return cmodel.Recipe{}, fmt.Errorf("unable to read recipe: %v", err)
	}

	return m, nil
}
