package domain

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/repository"
)

// UpdateRecipe updates a recipe.
func (d *Domain) UpdateRecipe(ctx context.Context, recipe model.Recipe, updateMask []string) (model.Recipe, error) {
	if recipe.Parent.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if recipe.Id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	recipient, err := d.repo.GetRecipeRecipient(ctx, recipe.Parent, recipe.Id)
	if err != nil {
		return model.Recipe{}, err
	}
	if recipient.PermissionLevel != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
	}
	if recipe.Parent.CircleId != 0 {
		permission, err := d.repo.GetCircleUserPermission(ctx, recipe.Parent.UserId, recipe.Parent.CircleId)
		if err != nil {
			return model.Recipe{}, err
		}
		if permission != permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}

	for _, updateMaskField := range updateMask {
		if updateMaskField == model.RecipeFields.ImageURI {
			recipe.ImageURI, err = d.updateImageURI(ctx, recipe)
			if err != nil {
				return model.Recipe{}, err
			}
		}
	}

	dbRecipe, err := tx.UpdateRecipe(ctx, recipe, updateMask)
	if err != nil {
		return model.Recipe{}, err
	}
	dbRecipe.Parent = recipe.Parent

	for _, updateMaskField := range updateMask {
		if updateMaskField == model.RecipeFields.IngredientGroups {
			dbRecipe.IngredientGroups = recipe.IngredientGroups
			dbRecipe.IngredientGroups, err = d.updateIngredientGroups(ctx, tx, dbRecipe)
			if err != nil {
				return model.Recipe{}, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

func (d *Domain) updateImageURI(ctx context.Context, recipe model.Recipe) (string, error) {
	if recipe.ImageURI == "" {
		err := d.removeRecipeImage(ctx, recipe.Parent, recipe.Id)
		if err != nil {
			return "", err
		}
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

func (d *Domain) removeRecipeImage(ctx context.Context, parent model.RecipeParent, id model.RecipeId) (err error) {
	recipe, err := d.GetRecipe(ctx, parent, id, []string{model.RecipeFields.ImageURI})
	if err != nil {
		return err
	}

	if recipe.ImageURI == "" {
		return nil
	}

	err = d.fileStore.DeleteFile(ctx, recipe.ImageURI)
	if err != nil {
		return err
	}

	return nil
}

func (d *Domain) updateIngredientGroups(ctx context.Context, tx repository.TxClient, dbRecipe model.Recipe) ([]model.IngredientGroup, error) {
	err := d.deleteIngredientGroups(ctx, tx, dbRecipe)
	if err != nil {
		return []model.IngredientGroup{}, err
	}

	dbRecipe.IngredientGroups, err = d.createIngredientGroups(ctx, tx, dbRecipe)
	if err != nil {
		return []model.IngredientGroup{}, err
	}

	return dbRecipe.IngredientGroups, nil
}

func (d *Domain) deleteIngredientGroups(ctx context.Context, tx repository.TxClient, recipe model.Recipe) error {
	filter := fmt.Sprintf("recipe_id = %d", recipe.Id.RecipeId)
	dbRecipeIngredients, err := tx.BulkDeleteRecipeIngredients(ctx, filter)
	if err != nil {
		return err
	}

	// delete ingredients if there are any
	if len(dbRecipeIngredients) > 0 {
		dbRecipeIds := []string{}
		for _, dbRecipeIngredient := range dbRecipeIngredients {
			dbRecipeIds = append(dbRecipeIds, strconv.FormatInt(dbRecipeIngredient.IngredientId, 10))
		}

		filter = fmt.Sprintf("ingredient_id = any(%s)", strings.Join(dbRecipeIds, ","))
		_, err = tx.BulkDeleteIngredients(ctx, filter)
		if err != nil {
			return err
		}
	}

	return nil
}
