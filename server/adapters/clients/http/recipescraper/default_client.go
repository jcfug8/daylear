package recipescraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/fx"

	"encoding/json"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/schemaorgrecipe"
	"github.com/jcfug8/daylear/server/ports/ingredientcleaner"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
	"github.com/rs/zerolog"
)

var _ recipescraper.DefaultClient = &DefaultClient{}

type DefaultClient struct {
	ingredientCleaner ingredientcleaner.Client
	log               zerolog.Logger
}

type DefaultClientParams struct {
	fx.In

	IngredientCleaner ingredientcleaner.Client
}

func NewDefaultClient(params DefaultClientParams, log zerolog.Logger) *DefaultClient {
	return &DefaultClient{
		ingredientCleaner: params.IngredientCleaner,
		log:               log,
	}
}

func (r *DefaultClient) ScrapeRecipe(ctx context.Context, uri string) (model.Recipe, error) {
	log := logutil.EnrichLoggerWithContext(r.log, ctx)
	// Fetch the HTML with a browser-like User-Agent
	request, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to create request")
		return model.Recipe{}, fmt.Errorf("failed to create request: %w", err)
	}
	// Set a common browser User-Agent
	request.Header.Set("User-Agent", "Daylear/1.0")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to fetch url")
		return model.Recipe{}, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Warn().Str("uri", uri).Int("status_code", resp.StatusCode).Msg("non-200 response")
		return model.Recipe{}, fmt.Errorf("non-200 response: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to parse HTML")
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
		log.Warn().Str("uri", uri).Msg("no schema.org recipe found in ld+json")
		return model.Recipe{}, errors.New("no schema.org recipe found in ld+json")
	}

	schemaRecipe.RecipeIngredient, err = r.ingredientCleaner.CleanIngredients(ctx, schemaorgrecipe.AsStringSlice(schemaRecipe.RecipeIngredient))
	if err != nil {
		log.Error().Err(err).Str("uri", uri).Msg("failed to clean ingredients")
		return model.Recipe{}, fmt.Errorf("failed to clean ingredients: %w", err)
	}

	return schemaorgrecipe.ToModelRecipe(schemaRecipe), nil
}
