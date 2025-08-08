package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	UserAccessTable = "user_access"
)

const (
	UserAccessColumn_UserAccessId    = "user_access_id"
	UserAccessColumn_UserId          = "user_id"
	UserAccessColumn_RequesterUserId = "requester_user_id"
	UserAccessColumn_RecipientUserId = "recipient_user_id"
	UserAccessColumn_PermissionLevel = "permission_level"
	UserAccessColumn_State           = "state"
)

var UserAccessFieldMasker = fieldmask.NewSQLFieldMasker(UserAccess{}, map[string][]fieldmask.Field{
	cmodel.UserAccessField_Parent:          {{Name: UserAccessColumn_UserId, Table: UserAccessTable}},
	cmodel.UserAccessField_Id:              {{Name: UserAccessColumn_UserAccessId, Table: UserAccessTable}},
	cmodel.UserAccessField_PermissionLevel: {{Name: UserAccessColumn_PermissionLevel, Table: UserAccessTable, Updatable: true}},
	cmodel.UserAccessField_State:           {{Name: UserAccessColumn_State, Table: UserAccessTable, Updatable: true}},

	cmodel.UserAccessField_Requester: {
		{Name: UserAccessColumn_RequesterUserId, Table: UserAccessTable},
	},
	cmodel.UserAccessField_Recipient: {
		{Name: UserAccessColumn_RecipientUserId, Table: UserAccessTable},
		{Name: UserColumn_Username, Table: UserTable, Alias: "recipient_username"},
		{Name: UserColumn_GivenName, Table: UserTable, Alias: "recipient_given_name"},
		{Name: UserColumn_FamilyName, Table: UserTable, Alias: "recipient_family_name"},
	},
})

var UserAccessSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"permission_level":  {Name: UserAccessColumn_PermissionLevel, Table: UserAccessTable},
	"state":             {Name: UserAccessColumn_State, Table: UserAccessTable},
	"recipient_user_id": {Name: UserAccessColumn_RecipientUserId, Table: UserAccessTable},
}, true)

// UserAccess -
type UserAccess struct {
	UserAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	UserId          int64                 `gorm:"not null;index;uniqueIndex:idx_user_access_user_recipient_unique,priority:1"`
	RequesterUserId int64                 `gorm:"index"`
	RecipientUserId int64                 `gorm:"not null;index;uniqueIndex:idx_user_access_user_recipient_unique,priority:2"`
	PermissionLevel types.PermissionLevel `gorm:"not null"`
	State           types.AccessState     `gorm:"not null"`

	RecipientUsername   string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName  string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (UserAccess) TableName() string {
	return UserAccessTable
}

type UA struct {
	UserAccess
}

// TableName -
func (UA) TableName() string {
	return "ua"
}
