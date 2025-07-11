package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/masks"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/genapi/api/types"
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

	Visibility types.VisibilityLevel

	Citation      string
	PrepDuration  time.Duration
	CookDuration  time.Duration
	TotalDuration time.Duration
	CookingMethod string
	Categories    []string
	YieldAmount   string
	Cuisines      []string
	CreateTime    time.Time
	UpdateTime    time.Time

	// The access details for the current user/circle
	RecipeAccess RecipeAccess
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
	Optional                bool
	MeasurementAmount       float64
	MeasurementType         pb.Recipe_MeasurementType
	MeasurementConjunction  pb.Recipe_Ingredient_MeasurementConjunction
	SecondMeasurementAmount float64
	SecondMeasurementType   pb.Recipe_MeasurementType
	Title                   string
}

// RecipeId defines the name for a recipe.
type RecipeId struct {
	RecipeId int64 `aip_pattern:"key=recipe"`
}

// RecipeParent defines the name for a recipe parent.
type RecipeParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
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

	AccessId:        "access_id",
	PermissionLevel: "permission_level",
	State:           "state",

	Citation:             "citation",
	PrepDurationSeconds:  "prep_duration_seconds",
	CookDurationSeconds:  "cook_duration_seconds",
	TotalDurationSeconds: "total_duration_seconds",
	CookingMethod:        "cooking_method",
	Categories:           "categories",
	YieldAmount:          "yield_amount",
	Cuisines:             "cuisines",
	CreateTime:           "create_time",
	UpdateTime:           "update_time",
}

type recipeFields struct {
	Id string

	Title       string
	Description string

	Directions       string
	IngredientGroups string

	ImageURI string

	VisibilityLevel string

	AccessId        string
	PermissionLevel string
	State           string

	Citation             string
	CookDurationSeconds  string
	PrepDurationSeconds  string
	TotalDurationSeconds string
	CookingMethod        string
	Categories           string
	YieldAmount          string
	Cuisines             string
	CreateTime           string
	UpdateTime           string
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

		fields.VisibilityLevel,

		fields.AccessId,
		fields.PermissionLevel,
		fields.State,

		fields.Citation,
		fields.CookDurationSeconds,
		fields.PrepDurationSeconds,
		fields.TotalDurationSeconds,
		fields.CookingMethod,
		fields.Categories,
		fields.YieldAmount,
		fields.Cuisines,
		fields.CreateTime,
		fields.UpdateTime,
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

		fields.VisibilityLevel,

		fields.Citation,
		fields.CookDurationSeconds,
		fields.PrepDurationSeconds,
		fields.TotalDurationSeconds,
		fields.CookingMethod,
		fields.Categories,
		fields.YieldAmount,
		fields.Cuisines,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
