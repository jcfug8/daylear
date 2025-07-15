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
	"github.com/jcfug8/daylear/server/ports/imagegenerator"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"google.golang.org/genai"
)

var _ recipescraper.Client = &RecipeGeminiClient{}
var _ imagegenerator.Client = &RecipeGeminiClient{}

var recipePromptTemplate = `
		  - Do not make up any content that is not there, but correct any typos or errors that seem useful.
		  - schema.org @type fields must included in the output.
		  - Use the full name of the unit, not the abbreviation. i.e. "1 cup" not "1 c" or "1 tablespoon" not "1 tbsp".
		  - the basic format of an ingredient should be '{amount} {unit} {ingredient}'.
		  - If there are clearly multiple sets of ingredients for different parts of the recipe, then the ingredients must be output as a itemList. The @type for each group shold be 'IngredientSection', the @type for each ingredient should be 'Ingredient', and the text of the ingredient should be keyed by 'text'. This is not in the schema.org/Recipe but that is how i want to format it.
		  - If an igredient item list an alternative ingredient append '[{alternative name}]' to the first ingredient name. The input may be formatted like, but not limited to, "1 cup honey or brown sugar" or "1 cup honey (brown sugar)"
		  - If the ingredient has two measurements that are to be treated as a range, then the output ingredient should be formatted like. "1 -- 2 cups sugar". The input may be formatted like, but not limited to, "1 - 2 cups sugar" or "1 to 2 cups sugar".
		  - If the ingredient has two measurements that should be combined or added together and are using two different units, then the output ingredient should be formatted like "1 cup && 1 tablespoon sugar". The input may be formatted like, but not limited to, "1 cup + 1 tablespoon sugar" or "1 cup and 1 tablespoon sugar".
		  - If the ingredient has two measurements but only one should be used, e.i one is volume and the other is weight, then the output ingredient should be formatted like. "1 cup || 100 grams sugar". The input may be formatted like, but not limited to, "1 cup or 100 grams sugar" or "1 cup (100 grams) sugar" or "1 cup sugar - 100 grams".
		  - Steps can, but don't have to, be grouped using a itemList if the steps are related to a specific section of the recipe. Do this if it would help clarify the recipe directions.
		  - If the recipe is in a foreign language, translate the recipe to English.
		  - If the ingredient is optional, end the ingredient with '*'.
`

type RecipeGeminiClient struct {
	logger                 zerolog.Logger
	client                 *genai.Client
	modelNamesLock         sync.Mutex
	modelNames             []string
	imageModelNames        []string
	imageModelNamesLock    sync.Mutex
	currentModelIndex      int
	currentImageModelIndex int
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

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
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

	imageModelNames := []string{
		"gemini-2.0-flash-preview-image-generation",
	}

	// randomize the model names
	rand.Shuffle(len(modelNames), func(i, j int) {
		modelNames[i], modelNames[j] = modelNames[j], modelNames[i]
	})

	return &RecipeGeminiClient{
		logger:            params.Logger,
		client:            client,
		modelNames:        modelNames,
		imageModelNames:   imageModelNames,
		currentModelIndex: 0,
	}, nil
}

func (c *RecipeGeminiClient) RecipeFromImage(ctx context.Context, files []file.File) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(c.logger, ctx)

	parts := []*genai.Part{
		genai.NewPartFromText(
			`Please convert the recipe on the image(s) to the schema.org/Recipe format as JSON. Only output the JSON. Do not include any other text or comments. If no recipe is found, return an empty JSON object and do not make up a recipe. 
		When formatting a recipe:
		  - The images may be out of order so please first figure out the order of the images, then parse the recipe from the images in the correct order` + recipePromptTemplate,
		),
	}
	for _, file := range files {
		// Read image into []byte
		imgBytes, err := io.ReadAll(file)
		if err != nil {
			return model.Recipe{}, fmt.Errorf("failed to read image: %w", err)
		}
		parts = append(parts, genai.NewPartFromBytes(imgBytes, file.ContentType))
	}

	content := genai.NewContentFromParts(parts, genai.RoleUser)

	var resp *genai.GenerateContentResponse
	tryCount := 0
	modelName := ""
	for {
		if tryCount >= len(c.modelNames) {
			log.Error().Msg("failed to get recipe from image after trying all models")
			return model.Recipe{}, fmt.Errorf("failed to get recipe from image after %d tries", tryCount)
		}
		tryCount++
		modelName = c.getModelName()
		log.Info().Str("model", modelName).Msg("attempting OCR")
		resp, err = c.client.Models.GenerateContent(ctx, modelName, []*genai.Content{content}, nil)
		if err != nil {
			log.Info().Err(err).Str("model", modelName).Msg("failed OCR")
			continue
		}
		break
	}

	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		text = string(part.Text)
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
	log.Info().Str("model", modelName).Interface("recipe", recipe).Interface("schemaRecipe", schemaRecipe).Msg("ran recipe from image")

	return recipe, nil
}

func (c *RecipeGeminiClient) RecipeFromData(ctx context.Context, data []byte) (recipe model.Recipe, err error) {
	log := logutil.EnrichLoggerWithContext(c.logger, ctx)

	parts := []*genai.Part{
		genai.NewPartFromText(
			`Please use the scraped web page data from https://recipeyumm.com/the-best-cowboy-caviar-ever/ and convert itinto the schema.org/Recipe format as JSON. Only output the JSON. Do not include any other text or comments. If no recipe is found, return an empty JSON object and do not make up a recipe. 
		When formatting a recipe:` + recipePromptTemplate,
		),
		genai.NewPartFromText(fmt.Sprintf("`Here is the scraped web page data: \n%s`", string(data))),
	}

	content := genai.NewContentFromParts(parts, genai.RoleUser)

	var resp *genai.GenerateContentResponse
	tryCount := 0
	modelName := ""
	for {
		if tryCount >= len(c.modelNames) {
			log.Error().Msg("failed to get recipe from data after trying all models")
			return model.Recipe{}, fmt.Errorf("failed to get recipe from data after %d tries", tryCount)
		}
		tryCount++
		modelName = c.getModelName()
		log.Info().Str("model", modelName).Msg("attempting to get recipe from data")
		resp, err = c.client.Models.GenerateContent(ctx, modelName, []*genai.Content{content}, nil)
		if err != nil {
			log.Info().Err(err).Str("model", modelName).Msg("failed to get recipe from data")
			continue
		}
		break
	}

	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		text = string(part.Text)
	}
	if text == "" {
		log.Warn().Msg("no text response from Gemini")
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
	log.Info().Str("model", modelName).Interface("recipe", recipe).Interface("schemaRecipe", schemaRecipe).Msg("ran recipe from data")

	return recipe, nil
}

func (c *RecipeGeminiClient) GenerateRecipeImage(ctx context.Context, recipe model.Recipe) (file.File, error) {
	log := logutil.EnrichLoggerWithContext(c.logger, ctx)

	schemaRecipe := schemaorgrecipe.ToSchemaOrgRecipe(recipe)
	schemaRecipeJson, err := json.Marshal(schemaRecipe)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal schema.org recipe")
		return file.File{}, fmt.Errorf("failed to marshal schema.org recipe: %w", err)
	}

	parts := []*genai.Part{
		genai.NewPartFromText(
			fmt.Sprintf(`Please generate a high-quality, realistic image of the finished dish described by the following schema.org/Recipe JSON. The image should look like a professional food photograph suitable for use as the main image on a recipe website. 

- The dish should be fully prepared and attractively presented, with attention to the key ingredients and typical serving style.
- Use natural lighting and a clean, appetizing background.
- Do not include any text, watermarks, or extraneous objects in the image.
- Focus on making the food look delicious and visually appealing.
- Pay attention to the ingredients and the recipe steps to ensure the image is accurate and realistic.

Here is the recipe:
				%s`,
				string(schemaRecipeJson),
			),
		),
	}

	content := genai.NewContentFromParts(parts, genai.RoleUser)

	var resp *genai.GenerateContentResponse
	modelName := ""
	tryCount := 0
	for {
		if tryCount >= len(c.imageModelNames) {
			return file.File{}, fmt.Errorf("failed to generate recipe image after %d tries", tryCount)
		}
		tryCount++
		modelName = c.getImageModelName()
		log.Info().Str("model", modelName).Msg("attempting to generate recipe image")
		resp, err = c.client.Models.GenerateContent(ctx, modelName, []*genai.Content{content}, &genai.GenerateContentConfig{
			ResponseModalities: []string{"TEXT", "IMAGE"},
		})
		if err != nil {
			log.Info().Err(err).Str("model", modelName).Msg("failed to generate recipe image")
			continue
		}
		break
	}

	var reader io.ReadSeekCloser
	var length int64
	var contentType string
	var extension string
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			if part.InlineData != nil {
				length = int64(len(part.InlineData.Data))
				reader = file.NewReadSeekCloser(part.InlineData.Data)
				contentType = part.InlineData.MIMEType
				extension = "." + strings.TrimPrefix(contentType, "image/")
				break
			}
		}
	}

	return file.File{
		Extension:      extension,
		ContentType:    contentType,
		ReadSeekCloser: reader,
		ContentLength:  int64(length),
	}, nil
}

func (c *RecipeGeminiClient) getImageModelName() string {
	c.imageModelNamesLock.Lock()
	defer c.imageModelNamesLock.Unlock()

	modelName := c.imageModelNames[c.currentImageModelIndex]
	c.currentImageModelIndex = (c.currentImageModelIndex + 1) % len(c.imageModelNames)
	return modelName
}

func (c *RecipeGeminiClient) getModelName() string {
	c.modelNamesLock.Lock()
	defer c.modelNamesLock.Unlock()

	modelName := c.modelNames[c.currentModelIndex]
	c.currentModelIndex = (c.currentModelIndex + 1) % len(c.modelNames)
	return modelName
}
