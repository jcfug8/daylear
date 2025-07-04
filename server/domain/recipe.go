package domain

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strconv"

	// "fmt"
	// "path"
	// "strconv"
	// "strings"

	"github.com/jcfug8/daylear/server/core/file"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	uuid "github.com/satori/go.uuid"
	// "github.com/jcfug8/daylear/server/ports/repository"
	// uuid "github.com/satori/go.uuid"
)

const RecipeImageRoot = "recipes"

// CreateRecipe creates a new recipe.
func (d *Domain) CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (dbRecipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	recipe.Id.RecipeId = 0

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			return model.Recipe{}, err
		}
		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}
	defer tx.Rollback()

	recipe.ImageURI, err = d.createRecipeImageURI(ctx, recipe)
	if err != nil {
		return model.Recipe{}, err
	}

	dbRecipe, err = tx.CreateRecipe(ctx, recipe)
	if err != nil {
		return model.Recipe{}, err
	}

	recipeAccess := model.RecipeAccess{
		RecipeAccessParent: model.RecipeAccessParent{
			RecipeId: dbRecipe.Id,
		},
		PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
		State:           types.AccessState_ACCESS_STATE_ACCEPTED,
	}

	if authAccount.CircleId != 0 {
		recipeAccess.Recipient = model.AuthAccount{
			CircleId: authAccount.CircleId,
		}
	} else {
		recipeAccess.Recipient = model.AuthAccount{
			UserId: authAccount.UserId,
		}
	}

	dbRecipeAccess, err := tx.CreateRecipeAccess(ctx, recipeAccess)
	if err != nil {
		return model.Recipe{}, err
	}

	dbRecipe.RecipeAccess = dbRecipeAccess

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (recipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
	if err != nil {
		return model.Recipe{}, err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}

	defer tx.Rollback()

	recipe, err = tx.DeleteRecipe(ctx, authAccount, id)
	if err != nil {
		return model.Recipe{}, err
	}

	err = tx.BulkDeleteRecipeAccess(ctx, model.RecipeAccessParent{RecipeId: id})
	if err != nil {
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return recipe, nil
}

// GetRecipe gets a recipe.
func (d *Domain) GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (recipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
		if err != nil {
			return model.Recipe{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, id)
		if err != nil {
			return model.Recipe{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_PUBLIC {
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	recipe, err = d.repo.GetRecipe(ctx, authAccount, id)
	if err != nil {
		return model.Recipe{}, err
	}

	if authAccount.PermissionLevel < recipe.RecipeAccess.PermissionLevel {
		recipe.RecipeAccess.PermissionLevel = authAccount.PermissionLevel
	}

	return recipe, nil
}

// ListRecipes lists recipes.
func (d *Domain) ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string) (recipes []model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	recipes, err = d.repo.ListRecipes(ctx, authAccount, pageSize, pageOffset, filter)
	if err != nil {
		return nil, err
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			return nil, err
		}
		for _, recipe := range recipes {
			if recipe.RecipeAccess.PermissionLevel < authAccount.PermissionLevel {
				recipe.RecipeAccess.PermissionLevel = authAccount.PermissionLevel
			}
		}
	}

	return recipes, nil

}

// UpdateRecipe updates a recipe.
func (d *Domain) UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, updateMask []string) (dbRecipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if recipe.Id.RecipeId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, recipe.Id)
		if err != nil {
			return model.Recipe{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, recipe.Id)
		if err != nil {
			return model.Recipe{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	previousDbRecipe, err := d.repo.GetRecipe(ctx, authAccount, recipe.Id)
	if err != nil {
		return model.Recipe{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		return model.Recipe{}, err
	}

	for _, updateMaskField := range updateMask {
		if updateMaskField == model.RecipeFields.ImageURI && recipe.ImageURI != previousDbRecipe.ImageURI {
			recipe.ImageURI, err = d.updateRecipeImageURI(ctx, authAccount, recipe)
			if err != nil {
				return model.Recipe{}, err
			}
		}
	}

	dbRecipe, err = d.repo.UpdateRecipe(ctx, authAccount, recipe, updateMask)
	if err != nil {
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

// Recipe Image Methods

// UploadRecipeImage uploads a recipe image.
func (d *Domain) UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, imageReader io.Reader) (imageURI string, err error) {
	if authAccount.UserId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
	if err != nil {
		return "", err
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		return "", domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	recipe, err := d.repo.GetRecipe(ctx, authAccount, id)
	if err != nil {
		return "", err
	}
	oldImageURI := recipe.ImageURI

	imageURI, err = d.uploadRecipeImage(ctx, id, imageReader)
	if err != nil {
		return "", err
	}

	_, err = d.repo.UpdateRecipe(ctx, authAccount, model.Recipe{
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.RecipeFields.ImageURI})
	if err != nil {
		return "", err
	}

	go d.fileStore.DeleteFile(context.Background(), oldImageURI)

	return imageURI, nil
}

func (d *Domain) ScrapeRecipe(ctx context.Context, authAccount model.AuthAccount, uri string) (recipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	recipe.Id.RecipeId = 0

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			return model.Recipe{}, err
		}
		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "invalid uri"}
	}

	host := parsedURI.Host

	var scraper recipescraper.DefaultClient

	scraper, ok := d.recipeScrapers[host]
	if !ok {
		scraper = d.defaultRecipeScraper
	}

	recipe, err = scraper.ScrapeRecipe(ctx, uri)
	if err != nil {
		return model.Recipe{}, err
	}

	return recipe, nil
}

func (d *Domain) OCRRecipe(ctx context.Context, authAccount model.AuthAccount, imageReaders []io.Reader) (recipe model.Recipe, err error) {
	if authAccount.UserId == 0 {
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	var files []file.File

	for _, imageReader := range imageReaders {
		file, err := d.fileInspector.Inspect(ctx, imageReader)
		if err != nil {
			return model.Recipe{}, err
		}
		defer file.Close()
		files = append(files, file)
	}

	recipe, err = d.recipeOCR.OCRRecipe(ctx, files)
	if err != nil {
		return model.Recipe{}, err
	}

	return recipe, nil
}

// Helper methods

func (d *Domain) createRecipeImageURI(ctx context.Context, recipe model.Recipe) (string, error) {
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

func (d *Domain) updateRecipeImageURI(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (string, error) {
	if recipe.ImageURI == "" {
		err := d.removeRecipeImage(ctx, authAccount, recipe.Id)
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

func (d *Domain) removeRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (err error) {
	recipe, err := d.GetRecipe(ctx, authAccount, id)
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
