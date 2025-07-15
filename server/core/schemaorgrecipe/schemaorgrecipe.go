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
		if lp == "&&" || lp == "||" || lp == "--" {
			conjIdx = i
			conjWord = lp
			break
		} else if i > 5 {
			break
		}
	}

	amountFound := false

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
			amountFound = true
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
			amountFound = true
			extraIdx = i
			amount2 = tempAmount
		}

		if !amountFound {
			return 0, "", "", 0, "", text
		}

		if unit2Idx == unit1Idx {
			unit1Idx += extraIdx - amount2Idx
		}
		unit2Idx += extraIdx - amount2Idx

		unit1 = parts[unit1Idx]
		unit2 = parts[unit2Idx]
		name = strings.Join(parts[unit2Idx+1:], " ")

		if unit1 == unit2 && (conjWord == "&&") {
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
		amountFound = true
		extraIdx = i
		amount1 = tempAmount
	}

	if !amountFound {
		return 0, "", "", 0, "", text
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
	case "&&":
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_AND
	case "||":
		return pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_OR
	case "--":
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
	coreRecipe.ImageURI = parseImageURI(schemaRecipe.Image)

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
	var ingredientGroups []model.IngredientGroup
	var parseIngredientGroups func(interface{}, string)
	parseIngredient := func(ing string) model.RecipeIngredient {
		amount1, unit1, conj, amount2, unit2, name := ParseIngredient(ing)
		measurementType1 := MapUnitToMeasurementType(unit1)
		measurementType2 := MapUnitToMeasurementType(unit2)
		if measurementType1 == pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED {
			name = unit1 + " " + name
			name = strings.TrimSpace(name)
		}
		optional := false
		if strings.HasSuffix(name, "*") {
			optional = true
			name = strings.TrimSuffix(name, "*")
		}
		return model.RecipeIngredient{
			Optional:                optional,
			MeasurementAmount:       amount1,
			MeasurementType:         measurementType1,
			MeasurementConjunction:  MapConjunctionToProto(conj),
			SecondMeasurementAmount: amount2,
			SecondMeasurementType:   measurementType2,
			Title:                   name,
		}
	}
	parseIngredientGroups = func(instr interface{}, sectionTitle string) {
		var ingredients []model.RecipeIngredient
		switch v := instr.(type) {
		case string:
			if v != "" {
				ingredients = append(ingredients, parseIngredient(v))
			}
		case []interface{}:
			for _, step := range v {
				switch st := step.(type) {
				case string:
					if st != "" {
						ingredients = append(ingredients, parseIngredient(st))
					}
				case map[string]interface{}:
					typeVal, _ := st["@type"].(string)
					if typeVal == "IngredientSection" || typeVal == "ItemList" {
						name, _ := st["name"].(string)
						parseIngredientGroups(st["itemListElement"], name)
					} else if typeVal == "Ingredient" {
						if txt, ok := st["text"].(string); ok && txt != "" {
							ingredients = append(ingredients, parseIngredient(txt))
						} else if txt, ok := st["name"].(string); ok && txt != "" {
							ingredients = append(ingredients, parseIngredient(txt))
						}
					}
				}
			}
		}
		if len(ingredients) > 0 {
			ingredientGroups = append(ingredientGroups, model.IngredientGroup{Title: sectionTitle, RecipeIngredients: ingredients})
		}
	}
	parseIngredientGroups(schemaRecipe.RecipeIngredient, "")
	if len(ingredientGroups) > 0 {
		coreRecipe.IngredientGroups = ingredientGroups
	}

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
					if typeVal == "HowToSection" || typeVal == "ItemList" {
						name, _ := st["name"].(string)
						parseInstructions(st["itemListElement"], name)
					} else if typeVal == "HowToStep" {
						if txt, ok := st["text"].(string); ok && txt != "" {
							steps = append(steps, txt)
						} else if txt, ok := st["name"].(string); ok && txt != "" {
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
	coreRecipe.Visibility = 1
	return coreRecipe
}

// parse image uri from schema.org/Recipe
func parseImageURI(imageData interface{}) string {
	switch v := imageData.(type) {
	case string:
		return v
	case []interface{}:
		for _, item := range v {
			if img, ok := item.(string); ok {
				return img
			} else if img, ok := item.(map[string]interface{}); ok {
				if imgURI, ok := img["url"].(string); ok {
					return imgURI
				}
			}
		}
	}
	return ""
}

// durationToISO8601 converts a time.Duration to an ISO 8601 duration string (e.g., PT1H30M)
func durationToISO8601(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	result := "PT"
	if hours > 0 {
		result += strconv.Itoa(hours) + "H"
	}
	if minutes > 0 {
		result += strconv.Itoa(minutes) + "M"
	}
	if seconds > 0 || result == "PT" {
		result += strconv.Itoa(seconds) + "S"
	}
	return result
}

// timeToRFC3339 converts a time.Time to an RFC3339 string, or empty if zero
func timeToRFC3339(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

// ingredientToString converts a RecipeIngredient to a string for schema.org
func ingredientToString(ing model.RecipeIngredient) string {
	// Try to reconstruct the original string as best as possible
	var parts []string
	if ing.MeasurementAmount > 0 {
		parts = append(parts, strings.TrimRight(strings.TrimRight(strconv.FormatFloat(ing.MeasurementAmount, 'f', -1, 64), "0"), "."))
	}
	if ing.MeasurementType != pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED {
		parts = append(parts, pb.Recipe_MeasurementType_name[int32(ing.MeasurementType)])
	}
	if ing.MeasurementConjunction != pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_UNSPECIFIED && ing.SecondMeasurementAmount > 0 {
		conj := ""
		switch ing.MeasurementConjunction {
		case pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_AND:
			conj = "and"
		case pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_OR:
			conj = "or"
		case pb.Recipe_Ingredient_MEASUREMENT_CONJUNCTION_TO:
			conj = "to"
		}
		if conj != "" {
			parts = append(parts, conj)
		}
		parts = append(parts, strings.TrimRight(strings.TrimRight(strconv.FormatFloat(ing.SecondMeasurementAmount, 'f', -1, 64), "0"), "."))
		if ing.SecondMeasurementType != pb.Recipe_MEASUREMENT_TYPE_UNSPECIFIED {
			parts = append(parts, pb.Recipe_MeasurementType_name[int32(ing.SecondMeasurementType)])
		}
	}
	if ing.Title != "" {
		parts = append(parts, ing.Title)
	}
	return strings.Join(parts, " ")
}

// directionsToSchemaOrg converts []model.RecipeDirection to schema.org format
func directionsToSchemaOrg(directions []model.RecipeDirection) interface{} {
	if len(directions) == 1 && directions[0].Title == "" {
		// Single section, no title: return steps as []string
		return directions[0].Steps
	}
	var out []map[string]interface{}
	for _, dir := range directions {
		section := map[string]interface{}{
			"@type":           "HowToSection",
			"name":            dir.Title,
			"itemListElement": dir.Steps,
		}
		out = append(out, section)
	}
	return out
}

func ingredientsToSchemaOrg(ingredients []model.IngredientGroup) interface{} {
	if len(ingredients) == 1 && ingredients[0].Title == "" {
		// Single section, no title: return steps as []string
		return ingredients[0].RecipeIngredients
	}
	var out []map[string]interface{}
	for _, group := range ingredients {
		var ings []map[string]interface{}
		for _, ingredient := range group.RecipeIngredients {
			ings = append(ings, map[string]interface{}{
				"@type": "HowToStep",
				"text":  ingredientToString(ingredient),
			})
		}
		section := map[string]interface{}{
			"@type":           "IngredientSection",
			"name":            group.Title,
			"itemListElement": ings,
		}
		out = append(out, section)
	}
	return out
}

// ToSchemaOrgRecipe converts a model.Recipe to a SchemaOrgRecipe
func ToSchemaOrgRecipe(r model.Recipe) SchemaOrgRecipe {
	return SchemaOrgRecipe{
		Context:            "https://schema.org/",
		Type:               "Recipe",
		Name:               r.Title,
		Description:        r.Description,
		Image:              r.ImageURI,
		RecipeIngredient:   ingredientsToSchemaOrg(r.IngredientGroups),
		RecipeInstructions: directionsToSchemaOrg(r.Directions),
		RecipeYield:        r.YieldAmount,
		DatePublished:      timeToRFC3339(r.CreateTime),
		PrepTime:           durationToISO8601(r.PrepDuration),
		CookTime:           durationToISO8601(r.CookDuration),
		TotalTime:          durationToISO8601(r.TotalDuration),
		RecipeCategory:     r.Categories,
		RecipeCuisine:      r.Cuisines,
		CookingMethod:      r.CookingMethod,
		MainEntityOfPage:   r.Citation,
	}
}
