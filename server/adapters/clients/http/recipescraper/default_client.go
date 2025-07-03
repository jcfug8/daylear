package recipescraper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"encoding/json"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/recipescraper"
)

var _ recipescraper.DefaultClient = &DefaultClient{}

type DefaultClient struct{}

func NewDefaultClient() *DefaultClient {
	return &DefaultClient{}
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
	var coreRecipe model.Recipe
	schemaTag := doc.Find("script[type='application/ld+json']")
	found := false
	schemaRecipes := []SchemaOrgRecipe{}
	schemaRecipe := SchemaOrgRecipe{}
	schemaTag.EachWithBreak(func(i int, s *goquery.Selection) bool {
		schemaRecipes = []SchemaOrgRecipe{}
		schemaRecipe = SchemaOrgRecipe{}

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
				if asString(rec.Type) == "Recipe" {
					schemaRecipe = rec
					found = true
					return false // break
				}
			}
		}
		// Try single recipe
		err = json.Unmarshal([]byte(jsonText), &schemaRecipe)
		if err == nil && asString(schemaRecipe.Type) == "Recipe" {
			found = true
			return false // break
		}
		return true // continue
	})
	if !found {
		return model.Recipe{}, errors.New("no schema.org recipe found in ld+json")
	}

	// Map fields to core model
	coreRecipe.Title = asString(schemaRecipe.Name)
	coreRecipe.Description = asString(schemaRecipe.Description)
	// Image (handle string or array)
	coreRecipe.ImageURI = asString(schemaRecipe.Image)
	// Ingredients (handle string, array, or []string)
	var ingredientGroup model.IngredientGroup
	var ingredients = asStringSlice(schemaRecipe.RecipeIngredient)
	for _, ing := range ingredients {
		amount, unit, name := parseIngredient(ing)
		measurementType := mapUnitToMeasurementType(unit)
		if measurementType == pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED {
			name = unit + " " + name
		}
		ingredientGroup.RecipeIngredients = append(ingredientGroup.RecipeIngredients, model.RecipeIngredient{
			Optional:          false,
			MeasurementAmount: amount,
			MeasurementType:   measurementType,
			Title:             name,
		})
	}
	coreRecipe.IngredientGroups = []model.IngredientGroup{ingredientGroup}
	// Directions (handle string, array, HowToStep, HowToSection)
	var directions []model.RecipeDirection
	var parseInstructions func(interface{}, string)
	parseInstructions = func(instr interface{}, sectionTitle string) {
		var steps []string
		switch v := instr.(type) {
		case string:
			if v != "" {
				steps = append(steps, v)
			}
		case []interface{}:
			for _, step := range v {
				switch st := step.(type) {
				case string:
					if st != "" {
						steps = append(steps, st)
					}
				case map[string]interface{}:
					typeVal, _ := st["@type"].(string)
					if typeVal == "HowToSection" {
						name, _ := st["name"].(string)
						parseInstructions(st["itemListElement"], name)
					} else if typeVal == "HowToStep" {
						if txt, ok := st["text"].(string); ok && txt != "" {
							steps = append(steps, txt)
						}
					}
				}
			}
		}
		if len(steps) > 0 {
			directions = append(directions, model.RecipeDirection{Title: sectionTitle, Steps: steps})
		}
	}
	parseInstructions(schemaRecipe.RecipeInstructions, "")
	if len(directions) > 0 {
		coreRecipe.Directions = directions
	}
	coreRecipe.Visibility = 300
	return coreRecipe, nil
}

// parseIngredient attempts to split an ingredient string into amount, unit, and name.
func parseIngredient(text string) (float64, string, string) {
	// Example: "1 cup sugar" or "2 (14.5 ounce) cans stewed tomatoes"
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return 0, "", text
	}
	// Try to parse the first part as a number (amount)
	amount, err := parseAmount(parts[0])
	if err != nil {
		return 0, "", text
	}
	if len(parts) == 1 {
		return amount, "", ""
	}
	// Next part is likely the unit
	unit := parts[1]
	name := strings.Join(parts[2:], " ")
	return amount, unit, name
}

// parseAmount parses a string like "1", "1/2", or "1-1/2" into a float64
func parseAmount(s string) (float64, error) {
	// Handle mixed numbers like "1-1/2"
	if strings.Contains(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			whole, err1 := strconv.Atoi(parts[0])
			frac, err2 := parseFraction(parts[1])
			if err1 == nil && err2 == nil {
				return float64(whole) + frac, nil
			}
		}
	}
	// Handle fractions like "1/2"
	if strings.Contains(s, "/") {
		return parseFraction(s)
	}
	return strconv.ParseFloat(s, 64)
}

func parseFraction(s string) (float64, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return 0, errors.New("invalid fraction")
	}
	num, err1 := strconv.Atoi(parts[0])
	den, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || den == 0 {
		return 0, errors.New("invalid fraction")
	}
	return float64(num) / float64(den), nil
}

// mapUnitToMeasurementType maps a unit string to the MeasurementType enum
func mapUnitToMeasurementType(unit string) pb.Recipe_MeasurementType {
	switch strings.ToLower(unit) {
	case "cup", "cups":
		return pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED // TODO: Add CUP if needed
	case "tablespoon", "tablespoons", "tbsp":
		return pb.Recipe_MEASUREMENT_TYPE_TABLESPOON
	case "teaspoon", "teaspoons", "tsp":
		return pb.Recipe_MEASUREMENT_TYPE_TEASPOON
	case "ounce", "ounces", "oz":
		return pb.Recipe_MEASUREMENT_TYPE_OUNCE
	case "pound", "pounds", "lb", "lbs":
		return pb.Recipe_MEASUREMENT_TYPE_POUND
	case "gram", "grams", "g":
		return pb.Recipe_MEASUREMENT_TYPE_GRAM
	case "milliliter", "milliliters", "ml":
		return pb.Recipe_MEASUREMENT_TYPE_MILLILITER
	case "liter", "liters", "l":
		return pb.Recipe_MEASUREMENT_TYPE_LITER
	default:
		fmt.Println("unknown unit: ", unit)
		return pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED // TODO: Add more units if needed
	}
}
