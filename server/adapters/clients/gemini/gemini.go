package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/schemaorgrecipe"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/ingredientcleaner"
	"github.com/jcfug8/daylear/server/ports/recipeocr"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/api/option"

	genai "github.com/google/generative-ai-go/genai"
)

var _ recipeocr.Client = &RecipeGeminiClient{}
var _ ingredientcleaner.Client = &RecipeGeminiClient{}

type RecipeGeminiClient struct {
	logger            zerolog.Logger
	client            *genai.Client
	modelNamesLock    sync.Mutex
	modelNames        []string
	currentModelIndex int
}

type RecipeGeminiClientParams struct {
	fx.In

	Config config.Client
	Logger zerolog.Logger
}

func NewRecipeGeminiClient(params RecipeGeminiClientParams) (*RecipeGeminiClient, error) {
	config := params.Config.GetConfig()["gemini"].(map[string]interface{})
	apiKey, ok := config["apikey"].(string)
	if !ok {
		return nil, fmt.Errorf("missing gemini api key")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	modelNames := []string{
		"gemini-2.0-flash-lite",
		"gemini-2.0-flash",
		"gemini-2.5-flash-lite-preview-06-17",
		"gemini-2.5-flash-preview-04-17",
		"gemini-2.5-flash",
		"gemini-2.5-pro",
	}

	// randomize the model names
	rand.Shuffle(len(modelNames), func(i, j int) {
		modelNames[i], modelNames[j] = modelNames[j], modelNames[i]
	})

	return &RecipeGeminiClient{
		logger:            params.Logger,
		client:            client,
		modelNames:        modelNames,
		currentModelIndex: 0,
	}, nil
}

func (c *RecipeGeminiClient) OCRRecipe(ctx context.Context, files []file.File) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(c.logger, ctx)

	parts := []genai.Part{
		genai.Text(
			`Please convert the recipe on the image(s) to the schema.org/Recipe format as JSON. Only output the JSON. Do not include any other text or comments. If no recipe is found, return an empty JSON object and do not make up a recipe. 
		When formatting a recipe:
		  - The images may be out of order so please first figure out the order of the images, then parse the recipe from the images in the correct order.
		  - Do not make up any content that is not in the image, but correct any typos or errors that seem useful.
		  - Use the full name of the unit, not the abbreviation. i.e. "1 cup" not "1 c" or "1 tablespoon" not "1 tbsp".
		  - the basic format of an ingredient should be '{amount} {unit} {ingredient}'.
		  - If there are multiple ingredients with the same name, specify the part of the recipe that the ingredient is used in by appending '({part of recipe})' to the ingredient name. Do not do this unless there is actual ambiguity in the recipe.
		  - If an igredient item list an alternative ingredient append '[alternative name]' to the first ingredient name. The input may be formatted like, but not limited to, "1 cup honey or brown sugar" or "1 cup honey (brown sugar)"
		  - If the ingredient has two measurements that are to be treated as a range, then the output ingredient should be formatted like. "1 to 2 cups sugar". The input may be formatted like, but not limited to, "1 - 2 cups sugar" or "1 to 2 cups sugar".
		  - If the ingredient has two measurements that should be combined or added together and are using two different units, then the output ingredient should be formatted like "1 cup and 1 tablespoon sugar". The input may be formatted like, but not limited to, "1 cup + 1 tablespoon sugar" or "1 cup and 1 tablespoon sugar".
		  - If the ingredient has two measurements but only one should be used, e.i one is volume and the other is weight, then the output ingredient should be formatted like. "1 cup or 100 grams sugar". The input may be formatted like, but not limited to, "1 cup or 100 grams sugar" or "1 cup (100 grams) sugar" or "1 cup sugar - 100 grams".
		  - Steps can, but don't have to, be grouped using a itemListElement if the steps are related to a specific section of the recipe. Do this if it would help clarify the recipe directions.
		  - If the recipe is in a foreign language, translate the recipe to English.`,
		),
	}
	for _, file := range files {
		// Read image into []byte
		imgBytes, err := io.ReadAll(file)
		if err != nil {
			return model.Recipe{}, fmt.Errorf("failed to read image: %w", err)
		}
		parts = append(parts, genai.Blob{
			MIMEType: file.ContentType,
			Data:     imgBytes,
		})
	}

	var resp *genai.GenerateContentResponse
	tryCount := 0
	modelName := ""
	for {
		if tryCount > len(c.modelNames) {
			log.Error().Msg("failed to get recipe from gemini after trying all models")
			return model.Recipe{}, fmt.Errorf("failed to get recipe from gemini after %d tries", tryCount)
		}
		tryCount++
		modelName = c.getModelName()
		log.Info().Str("model", modelName).Msg("attempting OCR")
		modelHandle := c.client.GenerativeModel(modelName)

		resp, err = modelHandle.GenerateContent(ctx,
			parts...,
		)
		if err != nil {
			log.Info().Err(err).Str("model", modelName).Msg("failed OCR")
			continue
		}
		break
	}

	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		switch p := part.(type) {
		case genai.Text:
			text = string(p)
		}
	}
	if text == "" {
		log.Error().Msg("no text response from Gemini")
		return model.Recipe{}, fmt.Errorf("no text response from Gemini")
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var schemaRecipe schemaorgrecipe.SchemaOrgRecipe
	err = json.Unmarshal([]byte(text), &schemaRecipe)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal schema.org recipe")
		return model.Recipe{}, fmt.Errorf("failed to unmarshal schema.org recipe: %w", err)
	}

	recipe = schemaorgrecipe.ToModelRecipe(schemaRecipe)
	log.Info().Str("model", modelName).Interface("recipe", recipe).Interface("schemaRecipe", schemaRecipe).Msg("ran OCR")

	return recipe, nil
}

func (c *RecipeGeminiClient) CleanIngredients(ctx context.Context, ingredients []string) (cleanedIngredients []string, err error) {
	log := logutil.EnrichLoggerWithContext(c.logger, ctx)

	parts := []genai.Part{
		genai.Text(
			`Please clean up this newline separated list of recipe ingredients and return the cleaned newline separated list of ingredients. The output must be a newline separated list of ingredients. Do not include any other text or comments.
		When cleaning the list of recipe ingredients:
		  - Use the full name of the unit, not the abbreviation. i.e. "1 cup" not "1 c" or "1 tablespoon" not "1 tbsp".
		  - a basic format of an ingredient should be '{amount} {unit} {ingredient}'.
		  - If there are multiple ingredients with the same name, specify the part of the recipe that the ingredient is used in by appending '({part of recipe})' to the ingredient name. Do not do this unless there is actual ambiguity in the recipe.
		  - If an igredient item list an alternative ingredient append '[alternative name]' to the first ingredient name. The input may be formatted like, but not limited to, "1 cup honey or brown sugar" or "1 cup honey (brown sugar)"
		  - If the ingredient has two measurements that are to be treated as a range, then the output ingredient should be formatted like. "1 to 2 cups sugar". The input may be formatted like, but not limited to, "1 - 2 cups sugar" or "1 to 2 cups sugar".
		  - If the ingredient has two measurements that should be combined or added together and are using two different units, then the output ingredient should be formatted like "1 cup and 1 tablespoon sugar". The input may be formatted like, but not limited to, "1 cup + 1 tablespoon sugar" or "1 cup and 1 tablespoon sugar".
		  - If the ingredient has two measurements but only one should be used, e.i one is volume and the other is weight, then the output ingredient should be formatted like. "1 cup or 100 grams sugar". The input may be formatted like, but not limited to, "1 cup or 100 grams sugar" or "1 cup (100 grams) sugar" or "1 cup sugar - 100 grams".`,
		),
		genai.Text(strings.Join(ingredients, "\n")),
	}

	var resp *genai.GenerateContentResponse
	tryCount := 0
	modelName := ""
	for {
		if tryCount > len(c.modelNames) {
			log.Error().Msg("failed to get recipe from gemini after trying all models")
			return nil, fmt.Errorf("failed to get recipe from gemini after %d tries", tryCount)
		}
		tryCount++
		modelName = c.getModelName()
		log.Info().Str("model", modelName).Msg("attempting cleaning ingredients")
		modelHandle := c.client.GenerativeModel(modelName)

		resp, err = modelHandle.GenerateContent(ctx,
			parts...,
		)
		if err != nil {
			log.Info().Err(err).Str("model", modelName).Msg("failed cleaning ingredients")
			continue
		}
		break
	}

	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		switch p := part.(type) {
		case genai.Text:
			text = string(p)
		}
	}
	if text == "" {
		log.Warn().Msg("no text response from Gemini")
		return nil, fmt.Errorf("no text response from Gemini")
	}

	cleanedIngredients = strings.Split(text, "\n")

	log.Info().Str("model", modelName).Interface("ingredients", ingredients).Interface("cleanedIngredients", cleanedIngredients).Msg("cleaned ingredients")

	return cleanedIngredients, nil
}

func (c *RecipeGeminiClient) getModelName() string {
	c.modelNamesLock.Lock()
	defer c.modelNamesLock.Unlock()

	modelName := c.modelNames[c.currentModelIndex]
	c.currentModelIndex = (c.currentModelIndex + 1) % len(c.modelNames)
	return modelName
}
