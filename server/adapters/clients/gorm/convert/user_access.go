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
		RequesterUserId: access.Requester.UserId,
		RecipientUserId: access.Recipient.UserId,
		PermissionLevel: access.PermissionLevel,
		State:           access.State,
		AcceptTarget:    access.AcceptTarget,
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
		PermissionLevel: dbAccess.PermissionLevel,
		State:           dbAccess.State,
		Requester: cmodel.UserId{
			UserId: dbAccess.RequesterUserId,
		},
		Recipient: cmodel.UserId{
			UserId: dbAccess.RecipientUserId,
		},
		RecipientUsername: dbAccess.RecipientUsername,
		AcceptTarget:      dbAccess.AcceptTarget,
	}
}
