package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = RecipeAccess{}

const (
	RecipeAccessField_Parent          = "parent"
	RecipeAccessField_Id              = "id"
	RecipeAccessField_PermissionLevel = "permission_level"
	RecipeAccessField_State           = "state"
	RecipeAccessField_Requester       = "requester"
	RecipeAccessField_Recipient       = "recipient"
)

// RecipeAccess represents a user's or circle's access to a recipe
type RecipeAccess struct {
	RecipeAccessParent
	RecipeAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState

	Requester RecipeRecipientOrRequester

	Recipient             RecipeRecipientOrRequester
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

// GetPermissionLevel - returns the permission level of the recipe access.
func (r RecipeAccess) GetPermissionLevel() types.PermissionLevel {
	return r.PermissionLevel
}

// SetPermissionLevel - sets the permission level of the recipe access. Because this method is
// uses a value receiver, it returns a new instance of the RecipeAccess struct and the caller
// must assign the result to a variable.
func (r RecipeAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	r.PermissionLevel = permissionLevel
	return r
}

type RecipeRecipientOrRequester struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type RecipeAccessParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	RecipeId
}

type RecipeAccessId struct {
	RecipeAccessId int64 `aip_pattern:"key=access"`
}
