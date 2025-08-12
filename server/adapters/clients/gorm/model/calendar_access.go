package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CalendarAccessTable = "calendar_access"
)

const (
	CalendarAccessColumn_CalendarAccessId  = "calendar_access_id"
	CalendarAccessColumn_CalendarId        = "calendar_id"
	CalendarAccessColumn_RequesterUserId   = "requester_user_id"
	CalendarAccessColumn_RequesterCircleId = "requester_circle_id"
	CalendarAccessColumn_RecipientUserId   = "recipient_user_id"
	CalendarAccessColumn_RecipientCircleId = "recipient_circle_id"
	CalendarAccessColumn_PermissionLevel   = "permission_level"
	CalendarAccessColumn_State             = "state"
	CalendarAccessColumn_AcceptTarget      = "accept_target"
)

var CalendarAccessFieldMasker = fieldmask.NewSQLFieldMasker(CalendarAccess{}, map[string][]fieldmask.Field{
	cmodel.CalendarAccessField_Parent:          {{Name: CalendarAccessColumn_CalendarId, Table: CalendarAccessTable}},
	cmodel.CalendarAccessField_Id:              {{Name: CalendarAccessColumn_CalendarAccessId, Table: CalendarAccessTable}},
	cmodel.CalendarAccessField_PermissionLevel: {{Name: CalendarAccessColumn_PermissionLevel, Table: CalendarAccessTable, Updatable: true}},
	cmodel.CalendarAccessField_State:           {{Name: CalendarAccessColumn_State, Table: CalendarAccessTable, Updatable: true}},
	cmodel.CalendarAccessField_AcceptTarget:    {{Name: CalendarAccessColumn_AcceptTarget, Table: CalendarAccessTable, Updatable: true}},
	cmodel.CalendarAccessField_Requester: {
		{Name: CalendarAccessColumn_RequesterUserId, Table: CalendarAccessTable},
		{Name: CalendarAccessColumn_RequesterCircleId, Table: CalendarAccessTable},
	},
	cmodel.CalendarAccessField_Recipient: {
		{Name: CalendarAccessColumn_RecipientUserId, Table: CalendarAccessTable},
		{Name: CalendarAccessColumn_RecipientCircleId, Table: CalendarAccessTable},
		{Name: UserColumn_Username, Table: UserTable, Alias: "recipient_username"},
		{Name: UserColumn_GivenName, Table: UserTable, Alias: "recipient_given_name"},
		{Name: UserColumn_FamilyName, Table: UserTable, Alias: "recipient_family_name"},
		{Name: CircleColumn_Title, Table: CircleTable, Alias: "recipient_circle_title"},
		{Name: CircleColumn_Handle, Table: CircleTable, Alias: "recipient_circle_handle"},
	},
})
var CalendarAccessSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"level":               {Name: CalendarAccessColumn_PermissionLevel, Table: "calendar_access"},
	"state":               {Name: CalendarAccessColumn_State, Table: "calendar_access"},
	"recipient.user_id":   {Name: CalendarAccessColumn_RecipientUserId, Table: "calendar_access"},
	"recipient.circle_id": {Name: CalendarAccessColumn_RecipientCircleId, Table: "calendar_access"},
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
	return CalendarAccessTable
}
