package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ Access = CalendarAccess{}

// CalendarAccessFields defines the calendar access fields for filtering.
const (
	CalendarAccessField_Parent          = "parent"
	CalendarAccessField_Id              = "id"
	CalendarAccessField_PermissionLevel = "permission_level"
	CalendarAccessField_State           = "state"
	CalendarAccessField_AcceptTarget    = "accept_target"
	CalendarAccessField_Requester       = "requester"
	CalendarAccessField_Recipient       = "recipient"
)

// CalendarAccess represents a user's or circle's access to a calendar
type CalendarAccess struct {
	CalendarAccessParent
	CalendarAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState
	AcceptTarget    types.AcceptTarget

	// the requester of the access
	Requester CalendarRecipientOrRequester

	// the recipient of the access
	Recipient             CalendarRecipientOrRequester
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

// GetAccessId - returns the access id of the calendar access.
func (c CalendarAccess) GetAccessId() int64 {
	return c.CalendarAccessId.CalendarAccessId
}

// GetPermissionLevel - returns the permission level of the calendar access.
func (c CalendarAccess) GetPermissionLevel() types.PermissionLevel {
	return c.PermissionLevel
}

// GetAcceptTarget - returns the accept target of the calendar access.
func (c CalendarAccess) GetAcceptTarget() types.AcceptTarget {
	return c.AcceptTarget
}

// GetAccessState - returns the access state of the calendar access.
func (c CalendarAccess) GetAccessState() types.AccessState {
	return c.State
}

// GetRecipientCircleId - returns the circle id of the recipient.
func (c CalendarAccess) GetRecipientCircleId() CircleId {
	return CircleId{CircleId: c.Recipient.CircleId}
}

// GetRecipientUserId - returns the user id of the recipient.
func (c CalendarAccess) GetRecipientUserId() UserId {
	return UserId{UserId: c.Recipient.UserId}
}

// SetPermissionLevel - sets the permission level of the calendar access. Because this method is
// uses a value receiver, it returns a new instance of the CalendarAccess struct and the caller
// must assign the result to a variable.
func (c CalendarAccess) SetPermissionLevel(permissionLevel types.PermissionLevel) Access {
	c.PermissionLevel = permissionLevel
	return c
}

type CalendarRecipientOrRequester struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type CalendarAccessParent struct {
	UserId     int64 `aip_pattern:"key=user"`
	CircleId   int64 `aip_pattern:"key=circle"`
	CalendarId int64 `aip_pattern:"key=calendar"`
}

type CalendarAccessId struct {
	CalendarAccessId int64 `aip_pattern:"key=access"`
}
