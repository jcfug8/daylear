package convert

import (
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	coreModel "github.com/jcfug8/daylear/server/core/model"
)

// CoreCircleAccessToCircleAccess converts a core CircleAccess model to a gorm CircleAccess model.
func CoreCircleAccessToCircleAccess(access coreModel.CircleAccess) dbModel.CircleAccess {
	return dbModel.CircleAccess{
		CircleAccessId:    access.CircleAccessId.CircleAccessId,
		CircleId:          access.CircleAccessParent.CircleId.CircleId,
		RequesterUserId:   access.Requester.UserId,
		RequesterCircleId: access.Requester.CircleId,
		RecipientUserId:   access.Recipient.UserId,
		PermissionLevel:   access.PermissionLevel,
		State:             access.State,
		AcceptTarget:      access.AcceptTarget,
	}
}

// CircleAccessToCoreCircleAccess converts a gorm CircleAccess model to a core CircleAccess model.
func CircleAccessToCoreCircleAccess(dbAccess dbModel.CircleAccess) coreModel.CircleAccess {
	return coreModel.CircleAccess{
		CircleAccessId: coreModel.CircleAccessId{
			CircleAccessId: dbAccess.CircleAccessId,
		},
		CircleAccessParent: coreModel.CircleAccessParent{
			CircleId: coreModel.CircleId{
				CircleId: dbAccess.CircleId,
			},
		},
		Requester: coreModel.CircleRequester{
			UserId:   dbAccess.RequesterUserId,
			CircleId: dbAccess.RequesterCircleId,
		},
		Recipient: coreModel.UserId{
			UserId: dbAccess.RecipientUserId,
		},
		PermissionLevel:     dbAccess.PermissionLevel,
		State:               dbAccess.State,
		RecipientUsername:   dbAccess.RecipientUsername,
		RecipientGivenName:  dbAccess.RecipientGivenName,
		RecipientFamilyName: dbAccess.RecipientFamilyName,
		AcceptTarget:        dbAccess.AcceptTarget,
	}
}
