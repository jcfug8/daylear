package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/domain"
)

// GetRecipe gets a recipe.
func (d *Domain) GetRecipe(ctx context.Context, parent model.RecipeParent, id model.RecipeId, fieldMask []string) (model.Recipe, error) {
	if parent.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err := d.repo.GetRecipeRecipient(ctx, parent, id)
	if err != nil {
		return model.Recipe{}, err
	}
	if parent.CircleId != 0 {
		_, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
		if err != nil {
			return model.Recipe{}, err
		}
	}

	dbRecipe, err := d.repo.GetRecipe(ctx, model.Recipe{
		Id:     id,
		Parent: parent,
	}, fieldMask)
	if err != nil {
		return model.Recipe{}, err
	}
	dbRecipe.Parent = parent

	getIngredients := false
	for _, fieldMaskField := range fieldMask {
		if fieldMaskField == model.RecipeFields.IngredientGroups {
			getIngredients = true
		}
	}

	if getIngredients {
		filter := fmt.Sprintf("recipe_id = %d", dbRecipe.Id.RecipeId)
		recipeIngredients, err := d.repo.ListRecipeIngredients(ctx, nil, filter, nil)
		if err != nil {
			return model.Recipe{}, err
		}
		dbRecipe.SetRecipeIngredients(recipeIngredients)
	}

	return dbRecipe, nil
}
