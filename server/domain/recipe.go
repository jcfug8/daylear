package domain

import (
	"context"
	"io"

	// "fmt"
	// "path"
	// "strconv"
	// "strings"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	// "github.com/jcfug8/daylear/server/ports/repository"
	// uuid "github.com/satori/go.uuid"
)

// CreateRecipe creates a new recipe.
func (d *Domain) CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (model.Recipe, error) {
	// TODO: Implement CreateRecipe
	// Implementation commented out for refactoring
	/*
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

		err = tx.BulkCreateRecipeRecipients(ctx, []model.RecipeParent{dbRecipe.Parent}, dbRecipe.Id, permPb.PermissionLevel_PERMISSION_LEVEL_WRITE)
		if err != nil {
			return model.Recipe{}, err
		}

		err = tx.Commit()
		if err != nil {
			return model.Recipe{}, err
		}

		return dbRecipe, nil
	*/
	return model.Recipe{}, domain.ErrInternal{Msg: "CreateRecipe method not implemented"}
}

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.Recipe, error) {
	// TODO: Implement DeleteRecipe
	// Implementation commented out for refactoring
	/*
		if parent.UserId == 0 {
			return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
		}

		if id.RecipeId == 0 {
			return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
		}

		recipient, err := d.repo.GetRecipeRecipient(ctx, parent, id)
		if err != nil {
			return model.Recipe{}, err
		}
		if recipient.PermissionLevel != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
		}
		if parent.CircleId != 0 {
			permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
			if err != nil {
				return model.Recipe{}, err
			}
			if permission != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
				return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
			}
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

		err = tx.BulkDeleteRecipeRecipients(ctx, []model.RecipeParent{}, id)
		if err != nil {
			return model.Recipe{}, err
		}

		err = tx.Commit()
		if err != nil {
			return model.Recipe{}, err
		}

		return recipe, nil
	*/
	return model.Recipe{}, domain.ErrInternal{Msg: "DeleteRecipe method not implemented"}
}

// GetRecipe gets a recipe.
func (d *Domain) GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, fieldMask []string) (model.Recipe, error) {
	// TODO: Implement GetRecipe
	// Implementation commented out for refactoring
	/*
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
	*/
	return model.Recipe{}, domain.ErrInternal{Msg: "GetRecipe method not implemented"}
}

// ListRecipes lists recipes.
func (d *Domain) ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64) (recipes []model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	authAccount.VisibilityLevel = types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.verifyCircleAccess(ctx, authAccount)
		if err != nil {
			return nil, err
		}
	}

	if authAccount.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return nil, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	recipes, err = d.repo.ListRecipes(ctx, authAccount, pageSize, pageOffset)
	if err != nil {
		return nil, err
	}

	return recipes, nil

}

// UpdateRecipe updates a recipe.
func (d *Domain) UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, updateMask []string) (model.Recipe, error) {
	// TODO: Implement UpdateRecipe
	// Implementation commented out for refactoring
	/*
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
		if recipient.PermissionLevel != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}
		if recipe.Parent.CircleId != 0 {
			permission, err := d.repo.GetCircleUserPermission(ctx, recipe.Parent.UserId, recipe.Parent.CircleId)
			if err != nil {
				return model.Recipe{}, err
			}
			if permission != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
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
	*/
	return model.Recipe{}, domain.ErrInternal{Msg: "UpdateRecipe method not implemented"}
}

// Recipe Image Methods

// UploadRecipeImage uploads a recipe image.
func (d *Domain) UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, imageReader io.Reader) (imageURI string, err error) {
	// TODO: Implement UploadRecipeImage
	// Implementation commented out for refactoring
	/*
		if parent.UserId == 0 {
			return "", domain.ErrInvalidArgument{Msg: "parent required"}
		}

		if id.RecipeId == 0 {
			return "", domain.ErrInvalidArgument{Msg: "id required"}
		}

		recipient, err := d.repo.GetRecipeRecipient(ctx, parent, id)
		if err != nil {
			return "", err
		}
		if recipient.PermissionLevel != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return "", domain.ErrPermissionDenied{Msg: "user does not have write permission"}
		}
		if parent.CircleId != 0 {
			permission, err := d.repo.GetCircleUserPermission(ctx, parent.UserId, parent.CircleId)
			if err != nil {
				return "", err
			}
			if permission != permPb.PermissionLevel_PERMISSION_LEVEL_WRITE {
				return "", domain.ErrPermissionDenied{Msg: "circle does not have write permission"}
			}

		}

		recipe, err := d.GetRecipe(ctx, parent, id, []string{model.RecipeFields.ImageURI})
		if err != nil {
			return "", err
		}
		oldImageURI := recipe.ImageURI

		imageURI, err = d.uploadRecipeImage(ctx, id, imageReader)
		if err != nil {
			return "", err
		}

		_, err = d.repo.UpdateRecipe(ctx, model.Recipe{
			Parent:   parent,
			Id:       id,
			ImageURI: imageURI,
		}, []string{model.RecipeFields.ImageURI})
		if err != nil {
			return "", err
		}

		go d.fileStore.DeleteFile(context.Background(), oldImageURI)

		return imageURI, nil
	*/
	return "", domain.ErrInternal{Msg: "UploadRecipeImage method not implemented"}
}

// Helper methods (commented out)

/*

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


func (d *Domain) uploadRecipeImage(ctx context.Context, id model.RecipeId, imageReader io.Reader) (imageURL string, err error) {
	file, err := d.fileInspector.Inspect(ctx, imageReader)
	if err != nil {
		return "", err
	}
	defer file.Close()

	imagePath := path.Join(RecipeImageRoot, strconv.FormatInt(id.RecipeId, 10), uuid.NewV4().String())
	imagePath = fmt.Sprintf("%s%s", imagePath, file.Extension)

	iamgeURI, err := d.fileStore.UploadPublicFile(ctx, imagePath, file)
	if err != nil {
		return "", err
	}

	return iamgeURI, nil
}

*/
