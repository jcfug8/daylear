package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = UserAccess{}

// UserAccessFields defines the user access fields for filtering.
var UserAccessFields = userAccessFields{
	Level:           "permission_level",
	State:           "state",
	RecipientUser:   "user_id",
	RecipientCircle: "circle_id",
}

type userAccessFields struct {
	Level           string
	State           string
	RecipientUser   string
	RecipientCircle string
}

// UserAccess represents a user's or circle's access to a user
type UserAccess struct {
	UserAccessParent
	UserAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState

	Requester UserId

	Recipient           UserId
	RecipientUsername   string // username of the recipient
	RecipientGivenName  string // given name of the recipient
	RecipientFamilyName string // family name of the recipient
}

// GetPermissionLevel - returns the permission level of the user access.
func (u UserAccess) GetPermissionLevel() types.PermissionLevel {
	return u.PermissionLevel
}

// SetPermissionLevel - sets the permission level of the user access. Because this method is
// uses a value receiver, it returns a new instance of the UserAccess struct and the caller
// must assign the result to a variable.
func (u UserAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	u.PermissionLevel = permissionLevel
	return u
}

type UserAccessParent struct {
	UserId
}

type UserAccessId struct {
	UserAccessId int64 `aip_pattern:"key=access"`
}
