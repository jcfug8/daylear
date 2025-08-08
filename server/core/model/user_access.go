package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = UserAccess{}

const (
	UserAccessField_Parent          = "parent"
	UserAccessField_Id              = "id"
	UserAccessField_PermissionLevel = "permission_level"
	UserAccessField_State           = "state"
	UserAccessField_Requester       = "requester"
	UserAccessField_Recipient       = "recipient"
)

// UserAccess represents a user's or circle's access to a user
type UserAccess struct {
	UserAccessParent
	UserAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState
	AcceptTarget    types.AcceptTarget

	Requester UserId

	Recipient           UserId
	RecipientUsername   string // username of the recipient
	RecipientGivenName  string // given name of the recipient
	RecipientFamilyName string // family name of the recipient
}

// GetAccessId - returns the access id of the user access.
func (u UserAccess) GetAccessId() int64 {
	return u.UserAccessId.UserAccessId
}

// GetPermissionLevel - returns the permission level of the user access.
func (u UserAccess) GetPermissionLevel() types.PermissionLevel {
	return u.PermissionLevel
}

// GetAcceptTarget - returns the accept target of the user access.
func (u UserAccess) GetAcceptTarget() types.AcceptTarget {
	return u.AcceptTarget
}

// GetAccessState - returns the access state of the user access.
func (u UserAccess) GetAccessState() types.AccessState {
	return u.State
}

// GetRecipientCircleId - returns the circle id of the recipient.
func (u UserAccess) GetRecipientCircleId() CircleId {
	return CircleId{}
}

// GetRecipientUserId - returns the user id of the recipient.
func (u UserAccess) GetRecipientUserId() UserId {
	return UserId{UserId: u.Recipient.UserId}
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
