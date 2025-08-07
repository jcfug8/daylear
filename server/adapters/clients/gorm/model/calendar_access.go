package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CalendarAccessColumn_CalendarAccessId  = "calendar_access.calendar_access_id"
	CalendarAccessColumn_CalendarId        = "calendar_access.calendar_id"
	CalendarAccessColumn_RequesterUserId   = "calendar_access.requester_user_id"
	CalendarAccessColumn_RequesterCircleId = "calendar_access.requester_circle_id"
	CalendarAccessColumn_RecipientUserId   = "calendar_access.recipient_user_id"
	CalendarAccessColumn_RecipientCircleId = "calendar_access.recipient_circle_id"
	CalendarAccessColumn_PermissionLevel   = "calendar_access.permission_level"
	CalendarAccessColumn_State             = "calendar_access.state"
)

var CalendarAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CalendarAccessField_Parent:          {CalendarAccessColumn_CalendarId},
	cmodel.CalendarAccessField_Id:              {CalendarAccessColumn_CalendarAccessId},
	cmodel.CalendarAccessField_PermissionLevel: {CalendarAccessColumn_PermissionLevel},
	cmodel.CalendarAccessField_State:           {CalendarAccessColumn_State},
	cmodel.CalendarAccessField_Recipient: {
		CalendarAccessColumn_RecipientUserId,
		CalendarAccessColumn_RecipientCircleId,
		UserFields.Username + " as recipient_username",
		UserFields.GivenName + " as recipient_given_name",
		UserFields.FamilyName + " as recipient_family_name",
		CircleColumn_Title + " as recipient_circle_title",
		CircleColumn_Handle + " as recipient_circle_handle",
	},
	cmodel.CalendarAccessField_Requester: {
		CalendarAccessColumn_RequesterUserId,
		CalendarAccessColumn_RequesterCircleId,
	},
})
var UpdateCalendarAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CalendarAccessField_PermissionLevel: {CalendarAccessColumn_PermissionLevel},
	cmodel.CalendarAccessField_State:           {CalendarAccessColumn_State},
})
var CalendarAccessSQLConverter = filter.NewSQLConverter(map[string]string{
	"permission_level":    CalendarAccessColumn_PermissionLevel,
	"state":               CalendarAccessColumn_State,
	"recipient.user_id":   CalendarAccessColumn_RecipientUserId,
	"recipient.circle_id": CalendarAccessColumn_RecipientCircleId,
}, true)

// CalendarAccess represents access control for a calendar
type CalendarAccess struct {
	CalendarAccessId  int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	CalendarId        int64                 `gorm:"not null;index;uniqueIndex:idx_calendar_id_recipient_user_id,idx_calendar_id_recipient_circle_id"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"not null;uniqueIndex:idx_calendar_id_recipient_user_id,where:recipient_user_id <> 0"`
	RecipientCircleId int64                 `gorm:"not null;uniqueIndex:idx_calendar_id_recipient_circle_id,where:recipient_circle_id <> 0"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`
	AcceptTarget      types.AcceptTarget    `gorm:"not null"`

	// Read-only fields from joins
	RecipientUsername     string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName    string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName   string `gorm:"->;-:migration"` // read only from join
	RecipientCircleTitle  string `gorm:"->;-:migration"` // read only from join
	RecipientCircleHandle string `gorm:"->;-:migration"` // read only from join
}

// TableName sets the table name for the CalendarAccess model.
func (CalendarAccess) TableName() string {
	return "calendar_access"
}
