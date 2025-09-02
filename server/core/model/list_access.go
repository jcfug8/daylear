package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = ListAccess{}

const (
	ListAccessField_Parent          = "parent"
	ListAccessField_Id              = "id"
	ListAccessField_PermissionLevel = "permission_level"
	ListAccessField_State           = "state"
	ListAccessField_AcceptTarget    = "accept_target"
	ListAccessField_Requester       = "requester"
	ListAccessField_Recipient       = "recipient"
)

// ListAccess represents a user's or circle's access to a list
type ListAccess struct {
	ListAccessParent
	ListAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState
	AcceptTarget    types.AcceptTarget

	Requester ListRecipientOrRequester

	Recipient             ListRecipientOrRequester
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

// GetAccessId - returns the access id of the list access.
func (l ListAccess) GetAccessId() int64 {
	return l.ListAccessId.ListAccessId
}

// GetPermissionLevel - returns the permission level of the list access.
func (l ListAccess) GetPermissionLevel() types.PermissionLevel {
	return l.PermissionLevel
}

// GetAcceptTarget - returns the accept target of the list access.
func (l ListAccess) GetAcceptTarget() types.AcceptTarget {
	return l.AcceptTarget
}

// GetAccessState - returns the access state of the list access.
func (l ListAccess) GetAccessState() types.AccessState {
	return l.State
}

// GetRecipientCircleId - returns the circle id of the recipient.
func (l ListAccess) GetRecipientCircleId() CircleId {
	return CircleId{CircleId: l.Recipient.CircleId}
}

// GetRecipientUserId - returns the user id of the recipient.
func (l ListAccess) GetRecipientUserId() UserId {
	return UserId{UserId: l.Recipient.UserId}
}

// SetPermissionLevel - sets the permission level of the list access. Because this method is
// uses a value receiver, it returns a new instance of the ListAccess struct and the caller
// must assign the result to a variable.
func (l ListAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	l.PermissionLevel = permissionLevel
	return l
}

type ListRecipientOrRequester struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type ListAccessParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	ListId
}

type ListAccessId struct {
	ListAccessId int64 `aip_pattern:"key=access"`
}
