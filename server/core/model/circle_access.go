package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

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
	Recipient           int64 `aip_pattern:"key=user"`
	Requester           AuthAccount
	PermissionLevel     types.PermissionLevel
	State               types.AccessState
	RecipientUsername   string // username of the recipient (if user)
	RecipientGivenName  string // given name of the recipient (if user)
	RecipientFamilyName string // family name of the recipient (if user)
}

type CircleAccessParent struct {
	CircleId
}

type CircleAccessId struct {
	CircleAccessId int64 `aip_pattern:"key=access"`
}
