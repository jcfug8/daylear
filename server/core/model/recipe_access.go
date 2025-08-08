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
	AcceptTarget    types.AcceptTarget

	Requester RecipeRecipientOrRequester

	Recipient             RecipeRecipientOrRequester
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

// GetAccessId - returns the access id of the recipe access.
func (r RecipeAccess) GetAccessId() int64 {
	return r.RecipeAccessId.RecipeAccessId
}

// GetPermissionLevel - returns the permission level of the recipe access.
func (r RecipeAccess) GetPermissionLevel() types.PermissionLevel {
	return r.PermissionLevel
}

// GetAcceptTarget - returns the accept target of the recipe access.
func (r RecipeAccess) GetAcceptTarget() types.AcceptTarget {
	return r.AcceptTarget
}

// GetAccessState - returns the access state of the recipe access.
func (r RecipeAccess) GetAccessState() types.AccessState {
	return r.State
}

// GetRecipientCircleId - returns the circle id of the recipient.
func (r RecipeAccess) GetRecipientCircleId() CircleId {
	return CircleId{CircleId: r.Recipient.CircleId}
}

// GetRecipientUserId - returns the user id of the recipient.
func (r RecipeAccess) GetRecipientUserId() UserId {
	return UserId{UserId: r.Recipient.UserId}
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
