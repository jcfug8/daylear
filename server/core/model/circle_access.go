package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
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
	Level types.PermissionLevel
	State pb.Access_State
}

type CircleAccessParent struct {
	CircleId
	Requester AuthAccount
	Recipient int64
}

type CircleAccessId struct {
	CircleAccessId int64 `aip_pattern:"key=access"`
}
