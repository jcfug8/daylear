package model

import "github.com/jcfug8/daylear/server/genapi/api/types"

type Access interface {
	GetPermissionLevel() types.PermissionLevel
	SetPermissionLevel(permissionLevel types.PermissionLevel) Access
}
