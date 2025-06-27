package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

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
	Level types.PermissionLevel
	State types.AccessState
}

type UserAccessParent struct {
	UserId
	Requester int64
	Recipient int64
}

type UserAccessId struct {
	UserAccessId int64 `aip_pattern:"key=access"`
}
