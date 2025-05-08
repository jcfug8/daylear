package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.Recipe, error) {
	if parent.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	permission, err := d.repo.GetRecipeUserPermission(ctx, parent.UserId, id.RecipeId)
	if err != nil {
		return model.Recipe{}, err
	}
	if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}

	defer tx.Rollback()

	recipe, err := tx.DeleteRecipe(ctx, model.Recipe{
		Id:     id,
		Parent: parent,
	})
	if err != nil {
		return model.Recipe{}, err
	}

	filter := fmt.Sprintf("recipe_id = %d", recipe.Id.RecipeId)
	recipeIngredients, err := tx.BulkDeleteRecipeIngredients(ctx, filter)
	if err != nil {
		return model.Recipe{}, err
	}
	recipe.SetRecipeIngredients(recipeIngredients)

	err = tx.BulkDeleteRecipeUsers(ctx, filter)
	if err != nil {
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return recipe, nil
}
