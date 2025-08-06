package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = CircleAccess{}

// CircleAccessFields defines the circle access fields for filtering.
var CircleAccessFields = circleAccessFields{
	Level:           "permission_level",
	State:           "state",
	RecipientUser:   "user_id",
	RecipientCircle: "circle_id",
}

type circleAccessFields struct {
	Level           string
	State           string
	RecipientUser   string
	RecipientCircle string
}

// CircleAccess represents a user's or circle's access to a circle
type CircleAccess struct {
	CircleAccessParent
	CircleAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState

	Requester CircleRequester

	Recipient           UserId
	RecipientUsername   string // username of the recipient
	RecipientGivenName  string // given name of the recipient
	RecipientFamilyName string // family name of the recipient
}

// GetPermissionLevel - returns the permission level of the circle access.
func (c CircleAccess) GetPermissionLevel() types.PermissionLevel {
	return c.PermissionLevel
}

// SetPermissionLevel - sets the permission level of the circle access. Because this method is
// uses a value receiver, it returns a new instance of the CircleAccess struct and the caller
// must assign the result to a variable.
func (c CircleAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	c.PermissionLevel = permissionLevel
	return c
}

type CircleRequester struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type CircleAccessParent struct {
	CircleId
}

type CircleAccessId struct {
	CircleAccessId int64 `aip_pattern:"key=access"`
}
