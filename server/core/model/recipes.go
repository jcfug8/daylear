package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// Recipe defines the model for a recipe.
type Recipe struct {
	Id     RecipeId
	Parent RecipeParent

	Title       string
	Description string

	Directions       []RecipeDirection
	IngredientGroups []IngredientGroup

	ImageURI string
}

type RecipeDirection struct {
	Title string
	Steps []string
}

type IngredientGroup struct {
	Title             string
	RecipeIngredients []RecipeIngredient `json:"-"`
}

type RecipeIngredient struct {
	RecipeIngredientId   int64
	Optional             bool
	MeasurementAmount    float64
	MeasurementType      pb.Recipe_MeasurementType
	IngredientGroupIndex int
	Ingredient
}

type Ingredient struct {
	IngredientId int64
	Title        string
}

func (r Recipe) GetIngredients() []Ingredient {
	ingredients := make([]Ingredient, 0)
	for _, ig := range r.IngredientGroups {
		for _, ri := range ig.RecipeIngredients {
			ingredients = append(ingredients, ri.Ingredient)
		}
	}
	return ingredients
}

func (r *Recipe) SetIngredients(ingredients []Ingredient) {
	k := 0
	for i := range r.IngredientGroups {
		for j := range r.IngredientGroups[i].RecipeIngredients {
			if k >= len(ingredients) {
				continue
			}
			r.IngredientGroups[i].RecipeIngredients[j].Ingredient = ingredients[k]
			k++
		}
	}
}

func (r *Recipe) SetRecipeIngredients(recipeIngredients []RecipeIngredient) {
	for i := range recipeIngredients {
		groupIndex := recipeIngredients[i].IngredientGroupIndex
		for len(r.IngredientGroups) <= int(groupIndex) {
			r.IngredientGroups = append(r.IngredientGroups, IngredientGroup{})
		}
		r.IngredientGroups[groupIndex].RecipeIngredients = append(r.IngredientGroups[groupIndex].RecipeIngredients, recipeIngredients[i])
	}
}

// RecipeId defines the name for a recipe.
type RecipeId struct {
	RecipeId int64 `aip_pattern:"key=recipe"`
}

// RecipeParent defines the owner for a recipe.
type RecipeParent struct {
	UserId   int64 `aip_pattern:"key=user,public_user"`
	CircleId int64 `aip_pattern:"key=circle,public_circle,circle"`
}

// ----------------------------------------------------------------------------
// Fields

// RecipeFields defines the recipe fields.
var RecipeFields = recipeFields{
	Id:     "id",
	Parent: "parent",

	Title:       "title",
	Description: "description",

	Directions:       "directions",
	IngredientGroups: "ingredient_groups",

	ImageURI: "image_uri",
}

type recipeFields struct {
	Id     string
	Parent string

	Title       string
	Description string

	Directions       string
	IngredientGroups string

	ImageURI string
}

// Mask returns a FieldMask for the recipe fields.
func (fields recipeFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Parent,

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
