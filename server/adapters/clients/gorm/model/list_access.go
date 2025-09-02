package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	ListAccessTable = "list_access"
)

const (
	ListAccessFields_ListAccessId      = "list_access_id"
	ListAccessFields_ListId            = "list_id"
	ListAccessFields_RequesterUserId   = "requester_user_id"
	ListAccessFields_RequesterCircleId = "requester_circle_id"
	ListAccessFields_RecipientUserId   = "recipient_user_id"
	ListAccessFields_RecipientCircleId = "recipient_circle_id"
	ListAccessFields_PermissionLevel   = "permission_level"
	ListAccessFields_State             = "state"
	ListAccessFields_AcceptTarget      = "accept_target"
)

var ListAccessFieldMasker = fieldmask.NewSQLFieldMasker(ListAccess{}, map[string][]fieldmask.Field{
	model.ListAccessField_Parent:          {{Name: ListAccessFields_ListId}},
	model.ListAccessField_Id:              {{Name: ListAccessFields_ListAccessId}},
	model.ListAccessField_PermissionLevel: {{Name: ListAccessFields_PermissionLevel}},
	model.ListAccessField_State:           {{Name: ListAccessFields_State}},
	model.ListAccessField_AcceptTarget:    {{Name: ListAccessFields_AcceptTarget}},
	model.ListAccessField_Requester: {
		{Name: ListAccessFields_RequesterUserId, Table: ListAccessTable},
		{Name: ListAccessFields_RequesterCircleId, Table: ListAccessTable},
	},
	model.ListAccessField_Recipient: {
		{Name: ListAccessFields_RecipientUserId, Table: ListAccessTable},
		{Name: ListAccessFields_RecipientCircleId, Table: ListAccessTable},
		{Name: UserColumn_Username, Table: UserTable, Alias: "recipient_username"},
		{Name: UserColumn_GivenName, Table: UserTable, Alias: "recipient_given_name"},
		{Name: UserColumn_FamilyName, Table: UserTable, Alias: "recipient_family_name"},
		{Name: CircleColumn_Title, Table: CircleTable, Alias: "recipient_circle_title"},
		{Name: CircleColumn_Handle, Table: CircleTable, Alias: "recipient_circle_handle"},
	},
})

var ListAccessSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	model.ListAccessField_PermissionLevel: {Name: ListAccessFields_PermissionLevel, Table: ListAccessTable},
	model.ListAccessField_State:           {Name: ListAccessFields_State, Table: ListAccessTable},
	model.ListAccessField_Requester:       {Name: ListAccessFields_RequesterUserId, Table: ListAccessTable},
	model.ListAccessField_Recipient:       {Name: ListAccessFields_RecipientCircleId, Table: ListAccessTable},
}, true)

// ListAccess -
type ListAccess struct {
	ListAccessId      int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	ListId            int64                 `gorm:"not null;index;uniqueIndex:idx_list_id_recipient_user_id;uniqueIndex:idx_list_id_recipient_circle_id;"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"index;uniqueIndex:idx_list_id_recipient_user_id,where:recipient_user_id <> null"`
	RecipientCircleId int64                 `gorm:"index;uniqueIndex:idx_list_id_recipient_circle_id,where:recipient_circle_id <> null"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`
	AcceptTarget      types.AcceptTarget    `gorm:"not null;default:0"`

	RecipientUsername     string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName    string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName   string `gorm:"->;-:migration"` // read only from join
	RecipientCircleTitle  string `gorm:"->;-:migration"` // read only from join
	RecipientCircleHandle string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (ListAccess) TableName() string {
	return ListAccessTable
}
