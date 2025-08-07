package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CircleAccessColumn_CircleAccessId    = "circle_access.circle_access_id"
	CircleAccessColumn_CircleId          = "circle_access.circle_id"
	CircleAccessColumn_RequesterUserId   = "circle_access.requester_user_id"
	CircleAccessColumn_RequesterCircleId = "circle_access.requester_circle_id"
	CircleAccessColumn_RecipientUserId   = "circle_access.recipient_user_id"
	CircleAccessColumn_PermissionLevel   = "circle_access.permission_level"
	CircleAccessColumn_State             = "circle_access.state"
)

var CircleAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CircleAccessField_Parent:          {CircleAccessColumn_CircleId},
	cmodel.CircleAccessField_Id:              {CircleAccessColumn_CircleAccessId},
	cmodel.CircleAccessField_PermissionLevel: {CircleAccessColumn_PermissionLevel},
	cmodel.CircleAccessField_State:           {CircleAccessColumn_State},

	cmodel.CircleAccessField_Requester: {
		CircleAccessColumn_RequesterUserId,
		CircleAccessColumn_RequesterCircleId,
	},
	cmodel.CircleAccessField_Recipient: {
		CircleAccessColumn_RecipientUserId,
		UserFields.Username + " as recipient_username",
		UserFields.GivenName + " as recipient_given_name",
		UserFields.FamilyName + " as recipient_family_name",
	},
})

var UpdateCircleAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CircleAccessField_PermissionLevel: {CircleAccessColumn_PermissionLevel},
	cmodel.CircleAccessField_State:           {CircleAccessColumn_State},
})

var CircleAccessSQLConverter = filter.NewSQLConverter(map[string]string{
	"permission_level":  CircleAccessColumn_PermissionLevel,
	"state":             CircleAccessColumn_State,
	"recipient.user_id": CircleAccessColumn_RecipientUserId,
}, true)

// CircleAccess -
type CircleAccess struct {
	CircleAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	CircleId          int64                 `gorm:"not null;index;uniqueIndex:idx_circle_id_recipient_user_id"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"not null;uniqueIndex:idx_circle_id_recipient_user_id,where:recipient_user_id <> 0"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`

	RecipientUsername   string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName  string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (CircleAccess) TableName() string {
	return "circle_access"
}

// CA
type CA struct {
	CircleAccess
}

// TableName -
func (CA) TableName() string {
	return "ca"
}
