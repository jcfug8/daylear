package domain

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/core/errz"
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (model.Recipe, error) {
	if parent.UserId == 0 {
		return model.Recipe{}, errz.NewInvalidArgument("parent required")
	}

	if id.RecipeId == 0 {
		return model.Recipe{}, errz.NewInvalidArgument("id required")
	}

	permission, err := d.repo.GetRecipeUserPermission(ctx, parent.UserId, id.RecipeId)
	if err != nil {
		return model.Recipe{}, err
	}
	if permission != pb.ShareRecipeRequest_RESOURCE_PERMISSION_WRITE {
		return model.Recipe{}, errz.NewPermissionDenied("user does not have write permission")
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
