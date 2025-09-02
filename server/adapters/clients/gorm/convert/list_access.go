package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// ListAccessFromCoreModel converts a core ListAccess to a GORM ListAccess
func ListAccessFromCoreModel(listAccess cmodel.ListAccess) gmodel.ListAccess {
	return gmodel.ListAccess{
		ListAccessId:      listAccess.ListAccessId.ListAccessId,
		ListId:            listAccess.ListAccessParent.ListId.ListId,
		RequesterUserId:   listAccess.Requester.UserId,
		RequesterCircleId: listAccess.Requester.CircleId,
		RecipientUserId:   listAccess.Recipient.UserId,
		RecipientCircleId: listAccess.Recipient.CircleId,
		PermissionLevel:   listAccess.PermissionLevel,
		State:             listAccess.State,
		AcceptTarget:      listAccess.AcceptTarget,
	}
}

// ListAccessToCoreModel converts a GORM ListAccess to a core ListAccess
func ListAccessToCoreModel(gormListAccess gmodel.ListAccess) cmodel.ListAccess {
	listAccess := cmodel.ListAccess{
		ListAccessParent: cmodel.ListAccessParent{
			ListId: cmodel.ListId{ListId: gormListAccess.ListId},
		},
		ListAccessId:    cmodel.ListAccessId{ListAccessId: gormListAccess.ListAccessId},
		PermissionLevel: gormListAccess.PermissionLevel,
		State:           gormListAccess.State,
		AcceptTarget:    gormListAccess.AcceptTarget,
		Requester: cmodel.ListRecipientOrRequester{
			UserId:   gormListAccess.RequesterUserId,
			CircleId: gormListAccess.RequesterCircleId,
		},
		Recipient: cmodel.ListRecipientOrRequester{
			UserId:   gormListAccess.RecipientUserId,
			CircleId: gormListAccess.RecipientCircleId,
		},
		RecipientUsername:     gormListAccess.RecipientUsername,
		RecipientGivenName:    gormListAccess.RecipientGivenName,
		RecipientFamilyName:   gormListAccess.RecipientFamilyName,
		RecipientCircleTitle:  gormListAccess.RecipientCircleTitle,
		RecipientCircleHandle: gormListAccess.RecipientCircleHandle,
	}

	return listAccess
}

// ListAccessListFromCoreModel converts a list of core ListAccess to a list of GORM ListAccess
func ListAccessListFromCoreModel(listAccesses []cmodel.ListAccess) []gmodel.ListAccess {
	res := make([]gmodel.ListAccess, len(listAccesses))
	for i, v := range listAccesses {
		res[i] = ListAccessFromCoreModel(v)
	}
	return res
}

// ListAccessListToCoreModel converts a list of GORM ListAccess to a list of core ListAccess
func ListAccessListToCoreModel(gormListAccesses []gmodel.ListAccess) []cmodel.ListAccess {
	res := make([]cmodel.ListAccess, len(gormListAccesses))
	for i, v := range gormListAccesses {
		res[i] = ListAccessToCoreModel(v)
	}
	return res
}
