package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessFields defines the recipe access fields for filtering.
var RecipeAccessFields = recipeAccessFields{
	Level:           "permission_level",
	State:           "state",
	RecipientUser:   "user_id",
	RecipientCircle: "circle_id",
}

type recipeAccessFields struct {
	Level           string
	State           string
	RecipientUser   string
	RecipientCircle string
}

// RecipeAccess represents a user's or circle's access to a recipe
type RecipeAccess struct {
	RecipeAccessParent
	RecipeAccessId
	Level types.PermissionLevel
	State pb.Access_State
}

type RecipeAccessParent struct {
	RecipeId
	Requester AuthAccount
	Recipient AuthAccount
}

type RecipeAccessId struct {
	RecipeAccessId int64 `aip_pattern:"key=access"`
}
