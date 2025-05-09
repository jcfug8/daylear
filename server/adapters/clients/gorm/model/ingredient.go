package model

// IngredientFields defines the Ingredient fields.
var IngredientFields = ingredientFields{
	IngredientId: "ingredient_id",
	Title:        "title",
}

type ingredientFields struct {
	IngredientId string
	Title        string
}

// Map maps the Ingredient fields to their corresponding model values.
func (fields ingredientFields) Map(m Ingredient) map[string]any {
	return map[string]any{
		fields.IngredientId: m.IngredientId,
		fields.Title:        m.Title,
	}
}

// Mask returns a FieldMask for the Ingredient fields.
func (fields ingredientFields) Mask() []string {
	return []string{
		fields.IngredientId,
		fields.Title,
	}
}

// Ingredient -
type Ingredient struct {
	IngredientId int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title        string
}

// TableName -
func (Ingredient) TableName() string {
	return "ingredient"
}
