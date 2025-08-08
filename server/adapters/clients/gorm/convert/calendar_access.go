package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// CalendarAccessToGorm converts a core CalendarAccess to a GORM CalendarAccess
func CalendarAccessToGorm(access cmodel.CalendarAccess) gmodel.CalendarAccess {
	return gmodel.CalendarAccess{
		CalendarAccessId:  access.CalendarAccessId.CalendarAccessId,
		CalendarId:        access.CalendarId,
		RequesterUserId:   access.Requester.UserId,
		RequesterCircleId: access.Requester.CircleId,
		RecipientUserId:   access.Recipient.UserId,
		RecipientCircleId: access.Recipient.CircleId,
		PermissionLevel:   access.PermissionLevel,
		State:             access.State,
		AcceptTarget:      access.AcceptTarget,
	}
}

// CalendarAccessFromGorm converts a GORM CalendarAccess to a core CalendarAccess
func CalendarAccessFromGorm(gormAccess gmodel.CalendarAccess) cmodel.CalendarAccess {
	return model.CalendarAccess{
		CalendarAccessParent: model.CalendarAccessParent{
			CalendarId: gormAccess.CalendarId,
		},
		CalendarAccessId: model.CalendarAccessId{
			CalendarAccessId: gormAccess.CalendarAccessId,
		},
		Requester: model.CalendarRecipientOrRequester{
			UserId:   gormAccess.RequesterUserId,
			CircleId: gormAccess.RequesterCircleId,
		},
		Recipient: model.CalendarRecipientOrRequester{
			UserId:   gormAccess.RecipientUserId,
			CircleId: gormAccess.RecipientCircleId,
		},
		PermissionLevel: gormAccess.PermissionLevel,
		State:           gormAccess.State,
		AcceptTarget:    gormAccess.AcceptTarget,
	}
}
