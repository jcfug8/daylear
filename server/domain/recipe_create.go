package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/repository"
)

// CreateRecipe creates a new recipe.
func (d *Domain) CreateRecipe(ctx context.Context, recipe model.Recipe) (model.Recipe, error) {
	if recipe.Parent.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	recipe.Id.RecipeId = 0

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}
	defer tx.Rollback()

	recipe.ImageURI, err = d.createImageURI(ctx, recipe)
	if err != nil {
		return model.Recipe{}, err
	}

	dbRecipe, err := tx.CreateRecipe(ctx, recipe)
	if err != nil {
		return model.Recipe{}, err
	}
	dbRecipe.Parent = recipe.Parent

	dbRecipe.IngredientGroups = recipe.IngredientGroups
	dbRecipe.IngredientGroups, err = d.createIngredientGroups(ctx, tx, dbRecipe)
	if err != nil {
		return model.Recipe{}, err
	}

	// link to user
	if dbRecipe.Parent.CircleId != 0 {
		err = tx.BulkCreateRecipeCircles(ctx, dbRecipe.Id, []int64{dbRecipe.Parent.CircleId}, permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE)
	} else {
		err = tx.BulkCreateRecipeUsers(ctx, dbRecipe.Id, []int64{dbRecipe.Parent.UserId}, permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE)
	}
	if err != nil {
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

func (d *Domain) createImageURI(ctx context.Context, recipe model.Recipe) (string, error) {
	if recipe.ImageURI == "" {
		return "", nil
	}

	fileContents, err := d.fileRetriever.GetFileContents(ctx, recipe.ImageURI)
	if err != nil {
		return "", err
	}
	defer fileContents.Close()

	imageURI, err := d.uploadRecipeImage(ctx, recipe.Id, fileContents)
	if err != nil {
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) createIngredientGroups(ctx context.Context, tx repository.TxClient, recipe model.Recipe) ([]model.IngredientGroup, error) {
	// create ingredients
	ingredients, err := tx.BulkCreateIngredients(ctx, recipe.GetIngredients())
	if err != nil {
		return []model.IngredientGroup{}, err
	}
	recipe.SetIngredients(ingredients)

	// link ingredients to recipe
	err = tx.SetRecipeIngredients(ctx, recipe.Id, recipe.IngredientGroups)
	if err != nil {
		return []model.IngredientGroup{}, err
	}

	return recipe.IngredientGroups, nil
}
