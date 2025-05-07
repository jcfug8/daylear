package model

// RecipeIngredientFields defines the recipeIngredient fields.
var RecipeIngredientFields = recipeIngredientFields{
	RecipeIngredientId: "recipe_ingredient_id",
	RecipeId:           "recipe_id",
	IngredientId:       "ingredient_id",
	MeasurementAmount:  "measurement_amount",
	MeasurementType:    "measurement_type",
}

type recipeIngredientFields struct {
	RecipeIngredientId string
	RecipeId           string
	IngredientId       string
	MeasurementAmount  string
	MeasurementType    string
}

// Map maps the recipeIngredient fields to their corresponding model values.
func (fields recipeIngredientFields) Map(m RecipeIngredient) map[string]any {
	return map[string]any{
		fields.RecipeIngredientId: m.RecipeIngredientId,
		fields.RecipeId:           m.RecipeId,
		fields.IngredientId:       m.IngredientId,
		fields.MeasurementAmount:  m.MeasurementAmount,
		fields.MeasurementType:    m.MeasurementType,
	}
}

// Mask returns a FieldMask for the recipeIngredient fields.
func (fields recipeIngredientFields) Mask() []string {
	return []string{
		fields.RecipeIngredientId,
		fields.RecipeId,
		fields.IngredientId,
		fields.MeasurementAmount,
		fields.MeasurementType,
	}
}
