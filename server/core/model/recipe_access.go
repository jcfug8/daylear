package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessFields defines the recipe access fields for filtering.
var RecipeAccessFields = recipeAccessFields{
	Level:           "permission_level",
	State:           "state",
	RecipientUser:   "recipient.user_id",
	RecipientCircle: "recipient.circle_id",
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
	PermissionLevel       types.PermissionLevel
	State                 types.AccessState
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

type RecipeAccessParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	RecipeId
	Requester AuthAccount
	Recipient AuthAccount
}

type RecipeAccessId struct {
	RecipeAccessId int64 `aip_pattern:"key=access"`
}
