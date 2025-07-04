package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

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
	client *genai.Client
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

	return &RecipeOCRClient{client: client}, nil
}

func (c *RecipeOCRClient) OCRRecipe(ctx context.Context, file file.File) (recipe model.Recipe, err error) {
	// Read image into []byte
	imgBytes, err := io.ReadAll(file)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to read image: %w", err)
	}

	modelHandle := c.client.GenerativeModel("gemini-2.0-flash-lite")

	prompt := "Please convert the recipe on this image to the schema.org/Recipe format as JSON. Only output the JSON. Do not include any other text or comments."

	resp, err := modelHandle.GenerateContent(ctx,
		genai.Blob{
			MIMEType: file.ContentType,
			Data:     imgBytes,
		},
		genai.Text(prompt),
	)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("gemini api error: %w", err)
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
