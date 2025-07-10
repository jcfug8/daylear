package schemaorgrecipe

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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

// ParseIngredient parses an ingredient string into two measurements (amount/unit), a conjunction, and a name
func ParseIngredient(text string) (amount1 float64, unit1 string, conj string, amount2 float64, unit2 string, name string) {
	text = ReplaceUnicodeFractions(text)
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return 0, "", "", 0, "", text
	}

	// Find conjunction index ("and", "or", "to")
	conjIdx := -1
	conjWord := ""
	for i, p := range parts {
		lp := strings.ToLower(p)
		if lp == "and" || lp == "+" || lp == "or" || lp == "to" || lp == "-" {
			conjIdx = i
			conjWord = lp
			break
		} else if i > 5 {
			break
		}
	}

	if conjIdx > 0 {
		amount1Idx := 0
		unit1Idx := 1
		amount2Idx := 3
		unit2Idx := 4

		// make sure to include stuff like "1 1/2"
		extraIdx := amount1Idx
		for i := extraIdx; i < len(parts); i++ {
			if i == conjIdx {
				unit1Idx = 3
				amount2Idx = 2
				unit2Idx = 3
				break
			}
			tempAmount, err := ParseAmount(strings.Join(parts[amount1Idx:i+1], " "))
			if err != nil {
				break
			}
			extraIdx = i
			amount1 = tempAmount
		}
		// push the unit1Idx, amount2Idx, and unit2Idx to the right
		unit1Idx += extraIdx - amount1Idx
		amount2Idx += extraIdx - amount1Idx
		unit2Idx += extraIdx - amount1Idx

		// make sure to include stuff like "1 1/2"
		extraIdx = amount2Idx
		for i := extraIdx; i < len(parts); i++ {
			tempAmount, err := ParseAmount(strings.Join(parts[amount2Idx:i+1], " "))
			if err != nil {
				break
			}
			extraIdx = i
			amount2 = tempAmount
		}

		if unit2Idx == unit1Idx {
			unit1Idx += extraIdx - amount2Idx
		}
		unit2Idx += extraIdx - amount2Idx

		unit1 = parts[unit1Idx]
		unit2 = parts[unit2Idx]
		name = strings.Join(parts[unit2Idx+1:], " ")

		if unit1 == unit2 && (conjWord == "and" || conjWord == "+") {
			return amount1 + amount2, unit1, "", 0, "", name
		}
		return amount1, unit1, conjWord, amount2, unit2, name
	}

	// Fallback: single measurement
	amountIdx := 0
	unitIdx := 1
	extraIdx := 0
	for i := extraIdx; i < len(parts); i++ {
		tempAmount, err := ParseAmount(strings.Join(parts[amountIdx:i+1], " "))
		if err != nil {
			break
		}
		extraIdx = i
		amount1 = tempAmount
	}

	unitIdx += extraIdx - amountIdx

	if len(parts) > unitIdx {
		unit1 = parts[unitIdx]
		name = strings.Join(parts[unitIdx+1:], " ")
	} else {
		unit1 = ""
		name = strings.Join(parts[unitIdx:], " ")
	}
	return amount1, unit1, "", 0, "", name
}

// ParseAmount parses a string like "1", "1/2", "1-1/2", "1 1/2", or "1 and 1/2" into a float64
func ParseAmount(s string) (float64, error) {
	s = strings.ReplaceAll(s, "-", " ")
	if strings.Contains(s, " ") {
		parts := strings.Split(s, " ")
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

// Helper to map conjunction string to proto enum
func MapConjunctionToProto(conj string) pb.Recipe_Ingredient_MeasurementConjunction {
	switch conj {
	case "and", "+":
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_AND
	case "or":
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_OR
	case "to", "-":
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_TO
	default:
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_UNSPECIFIED
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

// Helper to parse ISO 8601 duration (e.g., PT30M) to time.Duration
func parseISODuration(s string) (time.Duration, error) {
	// This is a simple implementation for PT#H#M#S
	// For more complex cases, use a library or extend as needed
	var d time.Duration
	var numStr string
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case 'P', 'T':
			// skip
		case 'H', 'M', 'S':
			if numStr != "" {
				num, err := strconv.Atoi(numStr)
				if err != nil {
					return 0, err
				}
				switch c {
				case 'H':
					d += time.Duration(num) * time.Hour
				case 'M':
					d += time.Duration(num) * time.Minute
				case 'S':
					d += time.Duration(num) * time.Second
				}
				numStr = ""
			}
		default:
			if c >= '0' && c <= '9' {
				numStr += string(c)
			}
		}
	}
	return d, nil
}

// Helper to parse ISO 8601 date/time (e.g., 2023-01-01T12:00:00Z)
func parseISOTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// ToModelRecipe converts a SchemaOrgRecipe to a model.Recipe
func ToModelRecipe(schemaRecipe SchemaOrgRecipe) model.Recipe {
	var coreRecipe model.Recipe
	coreRecipe.Title = AsString(schemaRecipe.Name)
	coreRecipe.Description = AsString(schemaRecipe.Description)
	coreRecipe.ImageURI = AsString(schemaRecipe.Image)

	// New fields
	// Citation: try mainEntityOfPage, url, isBasedOn
	coreRecipe.Citation = AsString(schemaRecipe.MainEntityOfPage)
	if coreRecipe.Citation == "" {
		coreRecipe.Citation = AsString(schemaRecipe.Context) // fallback, or add more sources as needed
	}
	// CookDuration: parse ISO 8601 duration from cookTime
	if cookTimeStr := AsString(schemaRecipe.CookTime); cookTimeStr != "" {
		if d, err := parseISODuration(cookTimeStr); err == nil {
			coreRecipe.CookDuration = d
		}
	}
	// PrepDuration: parse ISO 8601 duration from prepTime
	if prepTimeStr := AsString(schemaRecipe.PrepTime); prepTimeStr != "" {
		if d, err := parseISODuration(prepTimeStr); err == nil {
			coreRecipe.PrepDuration = d
		}
	}
	// TotalDuration: parse ISO 8601 duration from totalTime
	if totalTimeStr := AsString(schemaRecipe.TotalTime); totalTimeStr != "" {
		if d, err := parseISODuration(totalTimeStr); err == nil {
			coreRecipe.TotalDuration = d
		}
	}
	// CookingMethod
	coreRecipe.CookingMethod = AsString(schemaRecipe.CookingMethod)
	// Categories
	coreRecipe.Categories = AsStringSlice(schemaRecipe.RecipeCategory)
	// YieldAmount
	coreRecipe.YieldAmount = AsString(schemaRecipe.RecipeYield)
	// Cuisines
	coreRecipe.Cuisines = AsStringSlice(schemaRecipe.RecipeCuisine)
	// CreateTime: parse datePublished
	if dateStr := AsString(schemaRecipe.DatePublished); dateStr != "" {
		if t, err := parseISOTime(dateStr); err == nil {
			coreRecipe.CreateTime = t
		}
	}
	// UpdateTime: parse dateModified if present
	if dateStr := AsString(schemaRecipe.DatePublished); dateStr != "" {
		if t, err := parseISOTime(dateStr); err == nil {
			coreRecipe.UpdateTime = t
		}
	}

	// Ingredients
	var ingredientGroup model.IngredientGroup
	var ingredients = AsStringSlice(schemaRecipe.RecipeIngredient)
	for _, ing := range ingredients {
		amount1, unit1, conj, amount2, unit2, name := ParseIngredient(ing)
		measurementType1 := MapUnitToMeasurementType(unit1)
		measurementType2 := MapUnitToMeasurementType(unit2)
		if measurementType1 == pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED {
			name = unit1 + " " + name
			name = strings.TrimSpace(name)
		}
		ingredientGroup.RecipeIngredients = append(ingredientGroup.RecipeIngredients, model.RecipeIngredient{
			Optional:                false,
			MeasurementAmount:       amount1,
			MeasurementType:         measurementType1,
			MeasurementConjunction:  MapConjunctionToProto(conj),
			SecondMeasurementAmount: amount2,
			SecondMeasurementType:   measurementType2,
			Title:                   name,
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
