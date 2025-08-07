package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	UserAccessColumn_UserAccessId    = "user_access.user_access_id"
	UserAccessColumn_UserId          = "user_access.user_id"
	UserAccessColumn_RequesterUserId = "user_access.requester_user_id"
	UserAccessColumn_RecipientUserId = "user_access.recipient_user_id"
	UserAccessColumn_PermissionLevel = "user_access.permission_level"
	UserAccessColumn_State           = "user_access.state"
)

var UserAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.UserAccessField_Parent:          {UserAccessColumn_UserId},
	cmodel.UserAccessField_Id:              {UserAccessColumn_UserAccessId},
	cmodel.UserAccessField_PermissionLevel: {UserAccessColumn_PermissionLevel},
	cmodel.UserAccessField_State:           {UserAccessColumn_State},

	cmodel.UserAccessField_Requester: {
		UserAccessColumn_RequesterUserId,
	},
	cmodel.UserAccessField_Recipient: {
		UserAccessColumn_RecipientUserId,
		UserColumn_Username + " as recipient_username",
		UserColumn_GivenName + " as recipient_given_name",
		UserColumn_FamilyName + " as recipient_family_name",
	},
})

var UpdateUserAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.UserField_AccessPermissionLevel: {UserAccessColumn_PermissionLevel},
	cmodel.UserField_AccessState:           {UserAccessColumn_State},
})

var UserAccessSQLConverter = filter.NewSQLConverter(map[string]string{
	"permission_level":  UserAccessColumn_PermissionLevel,
	"state":             UserAccessColumn_State,
	"recipient_user_id": UserAccessColumn_RecipientUserId,
}, true)

// UserAccess -
type UserAccess struct {
	UserAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	UserId          int64                 `gorm:"not null;index;uniqueIndex:idx_user_access_user_recipient_unique,priority:1"`
	RequesterUserId int64                 `gorm:"index"`
	RecipientUserId int64                 `gorm:"not null;index;uniqueIndex:idx_user_access_user_recipient_unique,priority:2"`
	PermissionLevel types.PermissionLevel `gorm:"not null"`
	State           types.AccessState     `gorm:"not null"`

	RecipientUsername string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (UserAccess) TableName() string {
	return "user_access"
}

type UA struct {
	UserAccess
}

// TableName -
func (UA) TableName() string {
	return "ua"
}
