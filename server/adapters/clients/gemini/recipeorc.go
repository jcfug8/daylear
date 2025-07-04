package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/jcfug8/daylear/server/core/file"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/schemaorgrecipe"
	"github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/recipeocr"
	"go.uber.org/fx"
	"google.golang.org/api/option"

	genai "github.com/google/generative-ai-go/genai"
)

var _ recipeocr.Client = &RecipeOCRClient{}

type RecipeOCRClient struct {
	client            *genai.Client
	modelNamesLock    sync.Mutex
	modelNames        []string
	currentModelIndex int
}

type RecipeOCRClientParams struct {
	fx.In

	Config config.Client
}

func NewRecipeOCRClient(params RecipeOCRClientParams) (recipeocr.Client, error) {
	config := params.Config.GetConfig()["gemini"].(map[string]interface{})
	apiKey, ok := config["apikey"].(string)
	if !ok {
		return nil, fmt.Errorf("missing gemini api key")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &RecipeOCRClient{
		client: client,
		modelNames: []string{
			"gemini-2.0-flash-lite",
			"gemini-2.0-flash",
			"gemini-2.5-flash-lite-preview-06-17",
			"gemini-2.5-flash-preview-04-17",
			"gemini-2.5-flash",
			"gemini-2.5-pro",
		},
		currentModelIndex: 0,
	}, nil
}

func (c *RecipeOCRClient) OCRRecipe(ctx context.Context, files []file.File) (recipe model.Recipe, err error) {
	parts := []genai.Part{
		genai.Text(
			`Please convert the recipe on the image(s) to the schema.org/Recipe format as JSON. Only output the JSON. Do not include any other text or comments. If no recipe is found, return an empty JSON object and do not make up a recipe. 
		When formatting a recipe:
		  - The images may be out of order so please first figure out the order of the images, then parse the recipe from the images in the correct order.
		  - Do not make up any content that is not in the image, but correct any typos or errors that seem useful.
		  - Ingredients should be '{amount} {unit} {ingredient}'.
		  - Use the full name of the unit, not the abbreviation. i.e. "1 cup" not "1 c".
		  - If there are multiple ingredients with the same name, try to specify the part of the recipe that the ingredient is used in by appending '({part of recipe})' to the ingredient name.
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
	for {
		if tryCount > len(c.modelNames) {
			return model.Recipe{}, fmt.Errorf("failed to get recipe from gemini after %d tries", tryCount)
		}
		tryCount++
		modelHandle := c.client.GenerativeModel(c.getModelName())

		resp, err = modelHandle.GenerateContent(ctx,
			parts...,
		)
		if err != nil {
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
		return model.Recipe{}, fmt.Errorf("no text response from Gemini")
	}

	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var schemaRecipe schemaorgrecipe.SchemaOrgRecipe
	err = json.Unmarshal([]byte(text), &schemaRecipe)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to unmarshal schema.org recipe: %w", err)
	}

	recipe = schemaorgrecipe.ToModelRecipe(schemaRecipe)

	return recipe, nil
}

func (c *RecipeOCRClient) getModelName() string {
	c.modelNamesLock.Lock()
	defer c.modelNamesLock.Unlock()

	modelName := c.modelNames[c.currentModelIndex]
	c.currentModelIndex = (c.currentModelIndex + 1) % len(c.modelNames)
	return modelName
}
