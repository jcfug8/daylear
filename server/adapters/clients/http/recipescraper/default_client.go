package recipescraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/fx"

	"encoding/json"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/schemaorgrecipe"
	"github.com/jcfug8/daylear/server/ports/ingredientcleaner"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
)

var _ recipescraper.DefaultClient = &DefaultClient{}

type DefaultClient struct {
	ingredientCleaner ingredientcleaner.Client
}

type DefaultClientParams struct {
	fx.In

	IngredientCleaner ingredientcleaner.Client
}

func NewDefaultClient(params DefaultClientParams) *DefaultClient {
	return &DefaultClient{
		ingredientCleaner: params.IngredientCleaner,
	}
}

func (r *DefaultClient) ScrapeRecipe(ctx context.Context, uri string) (model.Recipe, error) {
	// Fetch the HTML
	resp, err := http.Get(uri)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return model.Recipe{}, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to parse HTML: %w", err)
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
		return model.Recipe{}, errors.New("no schema.org recipe found in ld+json")
	}

	schemaRecipe.RecipeIngredient, err = r.ingredientCleaner.CleanIngredients(ctx, schemaorgrecipe.AsStringSlice(schemaRecipe.RecipeIngredient))
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to clean ingredients: %w", err)
	}

	return schemaorgrecipe.ToModelRecipe(schemaRecipe), nil
}
