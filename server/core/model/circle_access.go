package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = CircleAccess{}

const (
	CircleAccessField_Parent          = "parent"
	CircleAccessField_Id              = "id"
	CircleAccessField_PermissionLevel = "permission_level"
	CircleAccessField_State           = "state"
	CircleAccessField_AcceptTarget    = "accept_target"
	CircleAccessField_Requester       = "requester"
	CircleAccessField_Recipient       = "recipient"
)

// CircleAccess represents a user's or circle's access to a circle
type CircleAccess struct {
	CircleAccessParent
	CircleAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState
	AcceptTarget    types.AcceptTarget

	Requester CircleRequester

	Recipient           UserId
	RecipientUsername   string // username of the recipient
	RecipientGivenName  string // given name of the recipient
	RecipientFamilyName string // family name of the recipient
}

// GetAccessId - returns the access id of the circle access.
func (c CircleAccess) GetAccessId() int64 {
	return c.CircleAccessId.CircleAccessId
}

// GetPermissionLevel - returns the permission level of the circle access.
func (c CircleAccess) GetPermissionLevel() types.PermissionLevel {
	return c.PermissionLevel
}

// GetAcceptTarget - returns the accept target of the circle access.
func (c CircleAccess) GetAcceptTarget() types.AcceptTarget {
	return c.AcceptTarget
}

// GetAccessState - returns the access state of the circle access.
func (c CircleAccess) GetAccessState() types.AccessState {
	return c.State
}

// SetPermissionLevel - sets the permission level of the circle access. Because this method is
// uses a value receiver, it returns a new instance of the CircleAccess struct and the caller
// must assign the result to a variable.
func (c CircleAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	c.PermissionLevel = permissionLevel
	return c
}

// GetRecipientCircleId - returns the circle id of the recipient.
func (c CircleAccess) GetRecipientCircleId() CircleId {
	return CircleId{}
}

// GetRecipientUserId - returns the user id of the recipient.
func (c CircleAccess) GetRecipientUserId() UserId {
	return UserId{UserId: c.Recipient.UserId}
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
