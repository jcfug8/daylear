package recipescraper

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
	Identifier         interface{} `json:"identifier"`         // https://schema.org/identifier
}

// Helper to extract a string from interface{}
func asString(v interface{}) string {
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
func asStringSlice(v interface{}) []string {
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
