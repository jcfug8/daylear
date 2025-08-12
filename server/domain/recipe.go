package domain

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/schemaorgrecipe"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

const RecipeImageRoot = "recipes"

// CreateRecipe creates a new recipe.
func (d *Domain) CreateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe) (dbRecipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if recipe.Parent.CircleId != 0 && authAccount.CircleId != 0 && recipe.Parent.CircleId != authAccount.CircleId {
		log.Warn().Msg("both circle ids set but do not match")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "both circle ids set but do not match"}
	}

	if recipe.Parent.CircleId != 0 {
		authAccount.CircleId = recipe.Parent.CircleId
	} else if recipe.Parent.UserId != 0 {
		authAccount.UserId = recipe.Parent.UserId
	}

	recipe.Id.RecipeId = 0

	if recipe.Parent.CircleId != 0 {
		_, err = d.determineCircleAccess(ctx, authAccount, model.CircleId{CircleId: authAccount.CircleId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when creating a recipe")
			return model.Recipe{}, err
		}
	} else if recipe.Parent.UserId != 0 {
		_, err = d.determineUserAccess(ctx, authAccount, model.UserId{UserId: authAccount.UserId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when creating a recipe")
			return model.Recipe{}, err
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

	dbRecipe, err = tx.CreateRecipe(ctx, recipe, nil)
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
		Requester: model.RecipeRecipientOrRequester{
			UserId: authAccount.AuthUserId,
		},
	}

	if authAccount.CircleId != 0 {
		recipeAccess.Recipient = model.RecipeRecipientOrRequester{
			CircleId: authAccount.CircleId,
		}
	} else {
		recipeAccess.Recipient = model.RecipeRecipientOrRequester{
			UserId: authAccount.AuthUserId,
		}
	}

	dbRecipeAccess, err := tx.CreateRecipeAccess(ctx, recipeAccess, nil)
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

	dbRecipe.Parent = recipe.Parent

	return dbRecipe, nil
}

// DeleteRecipe deletes a recipe.
func (d *Domain) DeleteRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	_, err = d.determineRecipeAccess(ctx, authAccount, id, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_ADMIN))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access")
		return model.Recipe{}, err
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

	recipe.Parent = parent

	return recipe, nil
}

// GetRecipe gets a recipe.
func (d *Domain) GetRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId, fields []string) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	recipe, err = d.repo.GetRecipe(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return model.Recipe{}, err
	}

	recipe.RecipeAccess, err = d.determineRecipeAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(recipe.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access")
		return model.Recipe{}, err
	}

	return recipe, nil
}

// ListRecipes lists recipes.
func (d *Domain) ListRecipes(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, pageSize int32, pageOffset int64, filter string, fields []string) (recipes []model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user_id required")
		return nil, domain.ErrInvalidArgument{Msg: "user_id required"}
	}

	authAccount.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_ADMIN

	if parent.CircleId != 0 {
		authAccount.CircleId = parent.CircleId
		dbCircle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: parent.CircleId}, []string{model.CircleField_Visibility})
		if err != nil {
			log.Error().Err(err).Msg("unable to get circle when listing recipes")
			return nil, domain.ErrInternal{Msg: "unable to get circle when listing recipes"}
		}
		determinedCircleAccess, err := d.determineCircleAccess(
			ctx, authAccount, model.CircleId{CircleId: parent.CircleId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(dbCircle.VisibilityLevel),
		)
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing recipes")
			return nil, err
		}
		authAccount.PermissionLevel = determinedCircleAccess.GetPermissionLevel()
	} else if parent.UserId != 0 {
		authAccount.UserId = parent.UserId
		determinedUserAccess, err := d.determineUserAccess(
			ctx, authAccount, model.UserId{UserId: authAccount.UserId},
			withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
			withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC))
		if err != nil {
			log.Error().Err(err).Msg("unable to determine access when listing recipes")
			return nil, err
		}
		authAccount.PermissionLevel = determinedUserAccess.GetPermissionLevel()
	}

	recipes, err = d.repo.ListRecipes(ctx, authAccount, pageSize, pageOffset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.ListRecipes failed")
		return nil, err
	}

	return recipes, nil
}

// UpdateRecipe updates a recipe.
func (d *Domain) UpdateRecipe(ctx context.Context, authAccount model.AuthAccount, recipe model.Recipe, fields []string) (dbRecipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if recipe.Id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	previousDbRecipe, err := d.repo.GetRecipe(ctx, authAccount, recipe.Id, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return model.Recipe{}, err
	}

	_, err = d.determineRecipeAccess(
		ctx, authAccount, recipe.Id,
		withResourceVisibilityLevel(previousDbRecipe.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access")
		return model.Recipe{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("repo.Begin failed")
		return model.Recipe{}, err
	}

	for _, updateMaskField := range fields {
		if updateMaskField == model.RecipeField_ImageURI && recipe.ImageURI != previousDbRecipe.ImageURI {
			recipe.ImageURI, err = d.updateRecipeImageURI(ctx, authAccount, recipe)
			if err != nil {
				log.Error().Err(err).Msg("updateRecipeImageURI failed")
				return model.Recipe{}, err
			}
		}
	}

	dbRecipe, err = d.repo.UpdateRecipe(ctx, authAccount, recipe, fields)
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipe failed")
		return model.Recipe{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("tx.Commit failed")
		return model.Recipe{}, err
	}

	dbRecipe.Parent = recipe.Parent

	return dbRecipe, nil
}

// Recipe Image Methods

// UploadRecipeImage uploads a recipe image.
func (d *Domain) UploadRecipeImage(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId, imageReader io.Reader) (imageURI string, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("parent required")
		return "", domain.ErrInvalidArgument{Msg: "parent required"}
	}

	if id.RecipeId == 0 {
		log.Warn().Msg("id required")
		return "", domain.ErrInvalidArgument{Msg: "id required"}
	}

	dbRecipe, err := d.repo.GetRecipe(ctx, authAccount, id, []string{model.RecipeField_ImageURI, model.RecipeField_VisibilityLevel})
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return "", err
	}
	oldImageURI := dbRecipe.ImageURI

	_, err = d.determineRecipeAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(dbRecipe.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access")
		return "", err
	}

	imageURI, err = d.uploadRecipeImage(ctx, id, imageReader)
	if err != nil {
		log.Error().Err(err).Msg("uploadRecipeImage failed")
		return "", err
	}

	_, err = d.repo.UpdateRecipe(ctx, authAccount, model.Recipe{
		Id:       id,
		ImageURI: imageURI,
	}, []string{model.RecipeField_ImageURI})
	if err != nil {
		log.Error().Err(err).Msg("repo.UpdateRecipe failed")
		return "", err
	}

	go d.fileStore.DeleteFile(context.Background(), oldImageURI)

	return imageURI, nil
}

func (d *Domain) GenerateRecipeImage(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId) (file.File, error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return file.File{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	dbRecipe, err := d.repo.GetRecipe(ctx, authAccount, id, nil)
	if err != nil {
		log.Error().Err(err).Msg("repo.GetRecipe failed")
		return file.File{}, err
	}

	_, err = d.determineRecipeAccess(
		ctx, authAccount, id,
		withResourceVisibilityLevel(dbRecipe.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access")
		return file.File{}, err
	}

	f, err := d.imageGenerator.GenerateRecipeImage(ctx, dbRecipe)
	if err != nil {
		log.Error().Err(err).Msg("imageGenerator.GenerateRecipeImage failed")
		return file.File{}, err
	}

	return f, nil
}

func (d *Domain) ScrapeRecipe(ctx context.Context, authAccount model.AuthAccount, uri string) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required")
		return model.Recipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	recipe.Id.RecipeId = 0

	request, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to create request")
		return model.Recipe{}, fmt.Errorf("failed to create request: %w", err)
	}
	// Set a common browser User-Agent
	request.Header.Set("User-Agent", "Daylear/1.0")
	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Language", "en-US,en;q=0.5")
	request.Header.Set("Accept-Encoding", "gzip, deflate")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to fetch url")
		return model.Recipe{}, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		l := log.Warn().Str("uri", uri).Int("status_code", resp.StatusCode)
		if len(body) > 0 {
			l.Str("body", string(body))
		}
		l.Msg("non-200 response")
		return model.Recipe{}, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return model.Recipe{}, err
		}
		defer r.Close()
		reader = r
	case "deflate":
		r := flate.NewReader(resp.Body)
		defer r.Close()
		reader = r
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("unable to read body", err.Error())
		return
	}

	extraRecipeData := d.appendSchemaOrgRecipe(log, body)
	body = bluemonday.StrictPolicy().SanitizeBytes(body)

	recipe, err = d.recipeScraper.RecipeFromData(ctx, body)
	if err != nil {
		log.Error().Err(err).Msg("recipeScraper.RecipeFromData failed")
		return model.Recipe{}, err
	}

	if extraRecipeData.ImageURI != "" {
		recipe.ImageURI = extraRecipeData.ImageURI
	}
	recipe.Citation = uri

	return recipe, nil
}

func (d *Domain) OCRRecipe(ctx context.Context, authAccount model.AuthAccount, imageReaders []io.Reader) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
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

	recipe, err = d.recipeScraper.RecipeFromImage(ctx, files)
	if err != nil {
		log.Error().Err(err).Msg("recipeScraper.RecipeFromImage failed")
		return model.Recipe{}, err
	}

	return recipe, nil
}

func (d *Domain) appendSchemaOrgRecipe(log zerolog.Logger, body []byte) model.Recipe {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse HTML")
		return model.Recipe{}
	}

	// Try to parse JSON-LD first
	schemaTag := doc.Find("script[type='application/ld+json']")
	found := false
	schemaRecipes := []schemaorgrecipe.SchemaOrgRecipe{}
	schemaRecipe := schemaorgrecipe.SchemaOrgRecipe{}
	schemaTag.EachWithBreak(func(i int, s *goquery.Selection) bool {
		schemaRecipes = []schemaorgrecipe.SchemaOrgRecipe{}
		schemaRecipe = schemaorgrecipe.SchemaOrgRecipe{}

		jsonText := s.Text()
		// Try @graph object
		var graphObj map[string]interface{}
		err := json.Unmarshal([]byte(jsonText), &graphObj)
		if err == nil {
			if graph, ok := graphObj["@graph"]; ok {
				if arr, ok := graph.([]interface{}); ok {
					for _, item := range arr {
						if m, ok := item.(map[string]interface{}); ok {
							typeVal, _ := m["@type"].(string)
							if typeVal == "Recipe" {
								b, _ := json.Marshal(m)
								_ = json.Unmarshal(b, &schemaRecipe)
								found = true
								return false // break
							}
						}
					}
				}
			}
		}
		// Try array of recipes
		err = json.Unmarshal([]byte(jsonText), &schemaRecipes)
		if err == nil && len(schemaRecipes) > 0 {
			for _, rec := range schemaRecipes {
				if schemaorgrecipe.AsString(rec.Type) == "Recipe" {
					schemaRecipe = rec
					found = true
					return false // break
				}
			}
		}
		// Try single recipe
		err = json.Unmarshal([]byte(jsonText), &schemaRecipe)
		if err == nil && schemaorgrecipe.AsString(schemaRecipe.Type) == "Recipe" {
			found = true
			return false // break
		}
		return true // continue
	})
	if !found {
		log.Warn().Msg("no schema.org recipe found in ld+json")
		return model.Recipe{}
	}

	return schemaorgrecipe.ToModelRecipe(schemaRecipe)
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
		err := d.removeRecipeImage(ctx, authAccount, recipe.Parent, recipe.Id)
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

func (d *Domain) removeRecipeImage(ctx context.Context, authAccount model.AuthAccount, parent model.RecipeParent, id model.RecipeId) (err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	recipe, err := d.GetRecipe(ctx, authAccount, parent, id, []string{model.RecipeField_ImageURI})
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
