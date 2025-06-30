package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// Recipe defines the model for a recipe.
type Recipe struct {
	Id RecipeId

	Title       string
	Description string

	Directions       []RecipeDirection
	IngredientGroups []IngredientGroup

	ImageURI string

	Visibility types.VisibilityLevel
	Permission types.PermissionLevel
}

type RecipeDirection struct {
	Title string
	Steps []string
}

type IngredientGroup struct {
	Title             string
	RecipeIngredients []RecipeIngredient
}

type RecipeIngredient struct {
	Optional          bool
	MeasurementAmount float64
	MeasurementType   pb.Recipe_MeasurementType
	Title             string
}

// RecipeId defines the name for a recipe.
type RecipeId struct {
	RecipeId int64 `aip_pattern:"key=recipe"`
}

// ----------------------------------------------------------------------------
// Fields

// RecipeFields defines the recipe fields.
var RecipeFields = recipeFields{
	Id: "id",

	Title:       "title",
	Description: "description",

	Directions:       "directions",
	IngredientGroups: "ingredient_groups",

	ImageURI: "image_uri",

	VisibilityLevel: "visibility_level",
	PermissionLevel: "permission_level",
}

type recipeFields struct {
	Id string

	Title       string
	Description string

	Directions       string
	IngredientGroups string

	ImageURI string

	VisibilityLevel string
	PermissionLevel string
}

// Mask returns a FieldMask for the recipe fields.
func (fields recipeFields) Mask() []string {
	return []string{
		fields.Id,

		fields.Title,
		fields.Description,

		fields.Directions,
		fields.IngredientGroups,

		fields.ImageURI,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields recipeFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Title,
		fields.Description,

		fields.Directions,
		fields.IngredientGroups,

		fields.ImageURI,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
