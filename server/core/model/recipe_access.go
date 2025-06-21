package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// RecipeAccessFields defines the recipe access fields for filtering.
var RecipeAccessFields = recipeAccessFields{
	Level: "permission_level",
	State: "state",
}

type recipeAccessFields struct {
	Level string
	State string
}

// RecipeAccessMap maps the core model fields to the database model fields.
var RecipeAccessMap = masks.NewFieldMap().
	MapFieldToFields(RecipeAccessFields.Level, "permission_level").
	MapFieldToFields(RecipeAccessFields.State, "state")

// RecipeAccess represents a user's or circle's access to a recipe
type RecipeAccess struct {
	RecipeAccessParent
	RecipeAccessId
	Level pb.Access_Level
	State pb.Access_State
}

type RecipeAccessParent struct {
	RecipeId
	Issuer    RecipeParent
	Recipient RecipeParent
}

type RecipeAccessId struct {
	RecipeAccessId int64 `aip_pattern:"key=access"`
}
