package model

import "github.com/jcfug8/daylear/server/genapi/api/types"

type Access interface {
	GetAccessId() int64
	GetPermissionLevel() types.PermissionLevel
	SetPermissionLevel(permissionLevel types.PermissionLevel) Access
	GetRecipientCircleId() CircleId
	GetRecipientUserId() UserId
	GetAcceptTarget() types.AcceptTarget
	GetAccessState() types.AccessState
}
