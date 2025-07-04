package schemaorgrecipe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// SchemaOrgRecipe matches the schema.org Recipe definition (generic)
type SchemaOrgRecipe struct {
	Context            interface{} `json:"@context"`           // https://schema.org/@context
	Type               interface{} `json:"@type"`              // https://schema.org/@type
	Name               interface{} `json:"name"`               // https://schema.org/name
	Description        interface{} `json:"description"`        // https://schema.org/description
	Image              interface{} `json:"image"`              // https://schema.org/image
	RecipeIngredient   interface{} `json:"recipeIngredient"`   // https://schema.org/recipeIngredient
	RecipeInstructions interface{} `json:"recipeInstructions"` // https://schema.org/recipeInstructions
	RecipeYield        interface{} `json:"recipeYield"`        // https://schema.org/recipeYield
	Nutrition          interface{} `json:"nutrition"`          // https://schema.org/nutrition
	Author             interface{} `json:"author"`             // https://schema.org/author
	DatePublished      interface{} `json:"datePublished"`      // https://schema.org/datePublished
	PrepTime           interface{} `json:"prepTime"`           // https://schema.org/prepTime
	CookTime           interface{} `json:"cookTime"`           // https://schema.org/cookTime
	TotalTime          interface{} `json:"totalTime"`          // https://schema.org/totalTime
	Keywords           interface{} `json:"keywords"`           // https://schema.org/keywords
	RecipeCategory     interface{} `json:"recipeCategory"`     // https://schema.org/recipeCategory
	RecipeCuisine      interface{} `json:"recipeCuisine"`      // https://schema.org/recipeCuisine
	AggregateRating    interface{} `json:"aggregateRating"`    // https://schema.org/aggregateRating
	SuitableForDiet    interface{} `json:"suitableForDiet"`    // https://schema.org/suitableForDiet
	Video              interface{} `json:"video"`              // https://schema.org/video
	CookingMethod      interface{} `json:"cookingMethod"`      // https://schema.org/cookingMethod
	MainEntityOfPage   interface{} `json:"mainEntityOfPage"`   // https://schema.org/mainEntityOfPage
	Identifier         interface{} `json:"identifier"`         // https://schema.org/identifier`
}

// Helper to extract a string from interface{}
func AsString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case []interface{}:
		if len(val) > 0 {
			if s, ok := val[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// Helper to extract []string from interface{}
func AsStringSlice(v interface{}) []string {
	var out []string
	switch val := v.(type) {
	case []interface{}:
		for _, item := range val {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
	case []string:
		return val
	case string:
		if val != "" {
			out = append(out, val)
		}
	}
	return out
}

// ToModelRecipe converts a SchemaOrgRecipe to a model.Recipe
func ToModelRecipe(schemaRecipe SchemaOrgRecipe) model.Recipe {
	var coreRecipe model.Recipe
	coreRecipe.Title = AsString(schemaRecipe.Name)
	coreRecipe.Description = AsString(schemaRecipe.Description)
	coreRecipe.ImageURI = AsString(schemaRecipe.Image)

	// Ingredients
	var ingredientGroup model.IngredientGroup
	var ingredients = AsStringSlice(schemaRecipe.RecipeIngredient)
	for _, ing := range ingredients {
		amount, unit, name := ParseIngredient(ing)
		measurementType := MapUnitToMeasurementType(unit)
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

	// Directions
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
	return coreRecipe
}

// ParseIngredient parses an ingredient string into amount, unit, and name
func ParseIngredient(text string) (float64, string, string) {
	text = ReplaceUnicodeFractions(text)
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return 0, "", text
	}
	amount := 0.0
	unitIdx := 1
	if len(parts) >= 2 {
		if frac, err := ParseFraction(parts[1]); err == nil {
			if whole, err := strconv.Atoi(parts[0]); err == nil {
				amount = float64(whole) + frac
				unitIdx = 2
			}
		}
		if len(parts) >= 3 && strings.ToLower(parts[1]) == "and" {
			if frac, err := ParseFraction(parts[2]); err == nil {
				if whole, err := strconv.Atoi(parts[0]); err == nil {
					amount = float64(whole) + frac
					unitIdx = 3
				}
			}
		}
	}
	if amount == 0.0 {
		var err error
		amount, err = ParseAmount(parts[0])
		if err != nil {
			return 0, "", text
		}
		unitIdx = 1
	}
	if len(parts) <= unitIdx {
		return amount, "", ""
	}
	unit := parts[unitIdx]
	name := strings.Join(parts[unitIdx+1:], " ")
	return amount, unit, name
}

// ParseAmount parses a string like "1", "1/2", "1-1/2", "1 1/2", or "1 and 1/2" into a float64
func ParseAmount(s string) (float64, error) {
	if strings.Contains(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) == 2 {
			whole, err1 := strconv.Atoi(parts[0])
			frac, err2 := ParseFraction(parts[1])
			if err1 == nil && err2 == nil {
				return float64(whole) + frac, nil
			}
		}
	}
	if strings.Contains(s, "/") {
		return ParseFraction(s)
	}
	return strconv.ParseFloat(s, 64)
}

// ParseFraction parses a string like "1/2" into a float64
func ParseFraction(s string) (float64, error) {
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid fraction")
	}
	num, err1 := strconv.Atoi(parts[0])
	den, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || den == 0 {
		return 0, fmt.Errorf("invalid fraction")
	}
	return float64(num) / float64(den), nil
}

// MapUnitToMeasurementType maps a unit string to the MeasurementType enum
func MapUnitToMeasurementType(unit string) pb.Recipe_MeasurementType {
	switch strings.ToLower(unit) {
	case "cup", "cups":
		return pb.Recipe_MEASUREMENT_TYPE_CUP
	case "tablespoon", "tablespoons", "tbsp", "tbs":
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
		return pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED
	}
}

// ReplaceUnicodeFractions replaces unicode fractions in a string with decimals
func ReplaceUnicodeFractions(s string) string {
	unicodeFractions := map[rune]float64{
		'¼': 0.25,
		'½': 0.5,
		'¾': 0.75,
		'⅓': 1.0 / 3.0,
		'⅔': 2.0 / 3.0,
		'⅛': 0.125,
		'⅜': 0.375,
		'⅝': 0.625,
		'⅞': 0.875,
	}
	out := ""
	for _, r := range s {
		if v, ok := unicodeFractions[r]; ok {
			out += fmt.Sprintf("%g", v)
		} else {
			out += string(r)
		}
	}
	return out
}
