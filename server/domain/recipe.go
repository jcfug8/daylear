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
	"github.com/jcfug8/daylear/server/core/logutil"
)

const RecipeImageRoot = "recipes"

// CreateRecipe creates a new recipe.
func (d *Domain) CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (dbRecipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("user id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	recipe.Id.RecipeId = 0

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
			return model.Recipe{}, err
		}
		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			log.Warn().Msg("user does not have access")
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Recipe{}, err
	}
	defer tx.Rollback()

	recipe.ImageURI, err = d.createRecipeImageURI(ctx, recipe)
	if err != nil {
		log.Error().Err(err).Msg("createRecipeImageURI failed")
		return model.Recipe{}, err
	}

	dbRecipe, err = tx.CreateRecipe(ctx, recipe)
	if err != nil {
		log.Error().Err(err).Msg("tx.CreateRecipe failed")
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
		log.Error().Err(err).Msg("tx.CreateRecipeAccess failed")
		return model.Recipe{}, err
	}

	dbRecipe.RecipeAccess = dbRecipeAccess

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevelsForCircle failed")
			return model.Recipe{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevels failed")
			return model.Recipe{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
		log.Warn().Msg("user does not have access")
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Recipe{}, err
	}

	defer tx.Rollback()

	recipe, err = tx.DeleteRecipe(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("tx.DeleteRecipe failed")
		return model.Recipe{}, err
	}

	if recipe.ImageURI != "" {
		go d.fileStore.DeleteFile(context.Background(), recipe.ImageURI)
	}

	err = tx.BulkDeleteRecipeAccess(ctx, model.RecipeAccessParent{RecipeId: id})
	if err != nil {
		log.Error().Err(err).Msg("tx.BulkDeleteRecipeAccess failed")
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Recipe{}, err
	}

	return recipe, nil
}

// GetRecipe gets a recipe.
func (d *Domain) GetRecipe(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevelsForCircle failed")
			return model.Recipe{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevels failed")
			return model.Recipe{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_PUBLIC {
		log.Warn().Msg("user does not have access")
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	recipe, err = d.repo.GetRecipe(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return model.Recipe{}, err
	}

	if authAccount.PermissionLevel < recipe.RecipeAccess.PermissionLevel {
		recipe.RecipeAccess.PermissionLevel = authAccount.PermissionLevel
	}

	return recipe, nil
}

// ListRecipes lists recipes.
func (d *Domain) ListRecipes(ctx context.Context, authAccount model.AuthAccount, pageSize int32, pageOffset int64, filter string) (recipes []model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("user_id required")
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	recipes, err = d.repo.ListRecipes(ctx, authAccount, pageSize, pageOffset, filter)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListRecipes failed")
		return nil, err
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
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
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if recipe.Id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, recipe.Id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevelsForCircle failed")
			return model.Recipe{}, err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, recipe.Id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevels failed")
			return model.Recipe{}, err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("user does not have access")
		return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	previousDbRecipe, err := d.repo.GetRecipe(ctx, authAccount, recipe.Id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return model.Recipe{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Recipe{}, err
	}

	for _, updateMaskField := range updateMask {
		if updateMaskField == model.RecipeFields.ImageURI && recipe.ImageURI != previousDbRecipe.ImageURI {
			recipe.ImageURI, err = d.updateRecipeImageURI(ctx, authAccount, recipe)
			if err != nil {
				log.Error().Err(err).Msg("updateRecipeImageURI failed")
				return model.Recipe{}, err
			}
		}
	}

	dbRecipe, err = d.repo.UpdateRecipe(ctx, authAccount, recipe, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipe failed")
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Recipe{}, err
	}

	return dbRecipe, nil
}

// Recipe Image Methods

// UploadRecipeImage uploads a recipe image.
func (d *Domain) UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, imageReader io.Reader) (imageURI string, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("parent required")
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevelsForCircle(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevelsForCircle failed")
			return "", err
		}
	} else {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, id)
		if err != nil {
			log.Error().Err(err).Msg("getRecipeAccessLevels failed")
			return "", err
		}
	}

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
		log.Warn().Msg("user does not have access")
		return "", domain.ErrPermissionDenied{Msg: "user does not have access"}
	}

	recipe, err := d.repo.GetRecipe(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return "", err
	}
	oldImageURI := recipe.ImageURI

	imageURI, err = d.uploadRecipeImage(ctx, id, imageReader)
	if err != nil {
		log.Error().Err(err).Msg("uploadRecipeImage failed")
		return "", err
	}

	_, err = d.repo.UpdateRecipe(ctx, authAccount, model.Recipe{
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.RecipeFields.ImageURI})
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipe failed")
		return "", err
	}

	go d.fileStore.DeleteFile(context.Background(), oldImageURI)

	return imageURI, nil
}

func (d *Domain) ScrapeRecipe(ctx context.Context, authAccount model.AuthAccount, uri string) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("user id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	recipe.Id.RecipeId = 0

	if authAccount.CircleId != 0 {
		authAccount.PermissionLevel, authAccount.VisibilityLevel, err = d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			log.Error().Err(err).Msg("getCircleAccessLevels failed")
			return model.Recipe{}, err
		}
		if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_WRITE {
			log.Warn().Msg("user does not have access")
			return model.Recipe{}, domain.ErrPermissionDenied{Msg: "user does not have access"}
		}
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		log.Warn().Err(err).Msg("invalid uri")
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
		log.Error().Err(err).Msg("scraper.ScrapeRecipe failed")
		return model.Recipe{}, err
	}

	return recipe, nil
}

func (d *Domain) OCRRecipe(ctx context.Context, authAccount model.AuthAccount, imageReaders []io.Reader) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.UserId == 0 {
		log.Warn().Msg("user id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	var files []file.File

	for _, imageReader := range imageReaders {
		image, err := d.imageClient.CreateImage(ctx, imageReader)
		if err != nil {
			log.Error().Err(err).Msg("imageClient.CreateImage failed")
			return model.Recipe{}, err
		}

		file, err := image.GetFile()
		if err != nil {
			log.Error().Err(err).Msg("image.GetFile failed")
			return model.Recipe{}, err
		}

		defer image.Remove(ctx)

		files = append(files, file)
	}

	recipe, err = d.recipeOCR.OCRRecipe(ctx, files)
	if err != nil {
		log.Error().Err(err).Msg("recipeOCR.OCRRecipe failed")
		return model.Recipe{}, err
	}

	return recipe, nil
}

// Helper methods

func (d *Domain) createRecipeImageURI(ctx context.Context, recipe model.Recipe) (string, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if recipe.ImageURI == "" {
		return "", nil
	}

	fileContents, err := d.fileRetriever.GetFileContents(ctx, recipe.ImageURI)
	if err != nil {
		log.Error().Err(err).Msg("fileRetriever.GetFileContents failed")
		return "", err
	}
	defer fileContents.Close()

	imageURI, err := d.uploadRecipeImage(ctx, recipe.Id, fileContents)
	if err != nil {
		log.Error().Err(err).Msg("uploadRecipeImage failed")
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) updateRecipeImageURI(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (string, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if recipe.ImageURI == "" {
		err := d.removeRecipeImage(ctx, authAccount, recipe.Id)
		if err != nil {
			log.Error().Err(err).Msg("removeRecipeImage failed")
			return "", err
		}
		return "", nil
	}

	fileContents, err := d.fileRetriever.GetFileContents(ctx, recipe.ImageURI)
	if err != nil {
		log.Error().Err(err).Msg("fileRetriever.GetFileContents failed")
		return "", err
	}
	defer fileContents.Close()

	imageURI, err := d.uploadRecipeImage(ctx, recipe.Id, fileContents)
	if err != nil {
		log.Error().Err(err).Msg("uploadRecipeImage failed")
		return "", err
	}

	return imageURI, nil
}

func (d *Domain) uploadRecipeImage(ctx context.Context, id model.RecipeId, imageReader io.Reader) (imageURL string, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	image, err := d.imageClient.CreateImage(ctx, imageReader)
	if err != nil {
		log.Error().Err(err).Msg("imageClient.CreateImage failed")
		return "", err
	}

	err = image.Convert(ctx, "jpeg")
	if err != nil {
		log.Error().Err(err).Msg("image.Convert failed")
		return "", err
	}

	width, height, err := image.GetDimensions(ctx)
	if err != nil {
		log.Error().Err(err).Msg("image.GetDimensions failed")
		return "", err
	}

	if width > maxImageWidth || height > maxImageHeight {
		newWidth, newHeight := resizeToFit(width, height, maxImageWidth)
		err = image.Resize(ctx, newWidth, newHeight)
		if err != nil {
			log.Error().Err(err).Msg("image.Resize failed")
			return "", err
		}
	}

	file, err := image.GetFile()
	if err != nil {
		log.Error().Err(err).Msg("image.GetFile failed")
		return "", err
	}
	defer image.Remove(ctx)

	imagePath := path.Join(RecipeImageRoot, strconv.FormatInt(id.RecipeId, 10), uuid.NewV4().String())
	imagePath = fmt.Sprintf("%s%s", imagePath, file.Extension)

	iamgeURI, err := d.fileStore.UploadPublicFile(ctx, imagePath, file)
	if err != nil {
		log.Error().Err(err).Msg("fileStore.UploadPublicFile failed")
		return "", err
	}

	return iamgeURI, nil
}

func (d *Domain) removeRecipeImage(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	recipe, err := d.GetRecipe(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("GetRecipe failed")
		return err
	}

	if recipe.ImageURI == "" {
		return nil
	}

	err = d.fileStore.DeleteFile(ctx, recipe.ImageURI)
	if err != nil {
		log.Error().Err(err).Msg("fileStore.DeleteFile failed")
		return err
	}

	return nil
}

func resizeToFit(width, height, maxDim int) (newWidth, newHeight int) {
	scale := float64(maxDim) / float64(width)
	if height > width {
		scale = float64(maxDim) / float64(height)
	}
	newWidth = int(float64(width) * scale)
	newHeight = int(float64(height) * scale)
	return
}
