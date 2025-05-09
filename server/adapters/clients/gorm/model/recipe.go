package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/core/model"
)

// RecipeMap maps the Recipe fields to their corresponding
// fields in the model.
var RecipeMap = masks.NewFieldMap().
	MapFieldToFields(model.RecipeFields.Id,
		RecipeFields.RecipeId).
	MapFieldToFields(model.RecipeFields.Title,
		RecipeFields.Title).
	MapFieldToFields(model.RecipeFields.Description,
		RecipeFields.Description).
	MapFieldToFields(model.RecipeFields.Directions,
		RecipeFields.Directions).
	MapFieldToFields(model.RecipeFields.ImageURI,
		RecipeFields.ImageURI).
	MapFieldToFields(model.RecipeFields.IngredientGroups,
		RecipeFields.IngredientGroups)

// RecipeFields defines the recipe fields.
var RecipeFields = recipeFields{
	RecipeId:         "recipe_id",
	Title:            "title",
	Description:      "description",
	Directions:       "directions",
	ImageURI:         "image_uri",
	IngredientGroups: "ingredient_groups",
}

type recipeFields struct {
	RecipeId         string
	Title            string
	Description      string
	Directions       string
	ImageURI         string
	IngredientGroups string
}

// Map maps the recipe fields to their corresponding model values.
func (fields recipeFields) Map(m Recipe) map[string]any {
	return map[string]any{
		fields.RecipeId:         m.RecipeId,
		fields.Title:            m.Title,
		fields.Description:      m.Description,
		fields.Directions:       m.Directions,
		fields.ImageURI:         m.ImageURI,
		fields.IngredientGroups: m.IngredientGroups,
	}
}

// Mask returns a FieldMask for the recipe fields.
func (fields recipeFields) Mask() []string {
	return []string{
		fields.RecipeId,
		fields.Title,
		fields.Description,
		fields.Directions,
		fields.ImageURI,
		fields.IngredientGroups,
	}
}

// Recipe -
type Recipe struct {
	RecipeId         int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title            string
	Description      string
	Directions       []byte `gorm:"type:jsonb"`
	ImageURI         string
	IngredientGroups []byte `gorm:"type:jsonb"`
}

// TableName -
func (Recipe) TableName() string {
	return "recipe"
}
