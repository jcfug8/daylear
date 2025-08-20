package model

import (
	"time"

	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ ResourceId = RecipeId{}

// ----------------------------------------------------------------------------
// Fields

// RecipeFields defines the recipe fields.
const (
	RecipeField_Parent               = "parent"
	RecipeField_Id                   = "id"
	RecipeField_Title                = "title"
	RecipeField_Description          = "description"
	RecipeField_Directions           = "directions"
	RecipeField_IngredientGroups     = "ingredient_groups"
	RecipeField_ImageURI             = "image_uri"
	RecipeField_VisibilityLevel      = "visibility_level"
	RecipeField_Citation             = "citation"
	RecipeField_PrepDurationSeconds  = "prep_duration_seconds"
	RecipeField_CookDurationSeconds  = "cook_duration_seconds"
	RecipeField_TotalDurationSeconds = "total_duration_seconds"
	RecipeField_CookingMethod        = "cooking_method"
	RecipeField_Categories           = "categories"
	RecipeField_YieldAmount          = "yield_amount"
	RecipeField_Cuisines             = "cuisines"
	RecipeField_CreateTime           = "create_time"
	RecipeField_UpdateTime           = "update_time"
	RecipeField_Favorited            = "favorited"

	RecipeField_RecipeAccess = "recipe_access"
)

// Recipe defines the model for a recipe.
type Recipe struct {
	Id               RecipeId
	Parent           RecipeParent
	Title            string
	Description      string
	Directions       []RecipeDirection
	IngredientGroups []IngredientGroup
	ImageURI         string
	VisibilityLevel  types.VisibilityLevel
	Citation         string
	PrepDuration     time.Duration
	CookDuration     time.Duration
	TotalDuration    time.Duration
	CookingMethod    string
	Categories       []string
	YieldAmount      string
	Cuisines         []string
	CreateTime       time.Time
	UpdateTime       time.Time
	Favorited        bool

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

// isResourceId - implements the ResourceId interface.
func (r RecipeId) isResourceId() {
}

// RecipeParent defines the name for a recipe parent.
// Supports both circle and user parents for recipes.
type RecipeParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	UserId   int64 `aip_pattern:"key=user"`
}
