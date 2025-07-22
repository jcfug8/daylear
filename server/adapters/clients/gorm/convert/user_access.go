package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// CoreUserAccessToUserAccess converts a core UserAccess model to a gorm UserAccess model.
func CoreUserAccessToUserAccess(access cmodel.UserAccess) gmodel.UserAccess {
	return gmodel.UserAccess{
		UserAccessId:    access.UserAccessId.UserAccessId,
		UserId:          access.UserAccessParent.UserId.UserId,
		RequesterUserId: access.Requester,
		RecipientUserId: access.Recipient,
		PermissionLevel: access.Level,
		State:           access.State,
	}
}

// UserAccessToCoreUserAccess converts a gorm UserAccess model to a core UserAccess model.
func UserAccessToCoreUserAccess(dbAccess gmodel.UserAccess) cmodel.UserAccess {
	return cmodel.UserAccess{
		UserAccessId: cmodel.UserAccessId{
			UserAccessId: dbAccess.UserAccessId,
		},
		UserAccessParent: cmodel.UserAccessParent{
			UserId: cmodel.UserId{
				UserId: dbAccess.UserId,
			},
		},
		Level:             dbAccess.PermissionLevel,
		State:             dbAccess.State,
		Requester:         dbAccess.RequesterUserId,
		Recipient:         dbAccess.RecipientUserId,
		RecipientUsername: dbAccess.RecipientUsername,
	}
}
