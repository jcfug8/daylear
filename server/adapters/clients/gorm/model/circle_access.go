package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CircleAccessTable = "circle_access"
)

const (
	CircleAccessColumn_CircleAccessId    = "circle_access_id"
	CircleAccessColumn_CircleId          = "circle_id"
	CircleAccessColumn_RequesterUserId   = "requester_user_id"
	CircleAccessColumn_RequesterCircleId = "requester_circle_id"
	CircleAccessColumn_RecipientUserId   = "recipient_user_id"
	CircleAccessColumn_PermissionLevel   = "permission_level"
	CircleAccessColumn_State             = "state"
	CircleAccessColumn_AcceptTarget      = "accept_target"
)

var CircleAccessFieldMasker = fieldmask.NewSQLFieldMasker(CircleAccess{}, map[string][]fieldmask.Field{
	cmodel.CircleAccessField_Parent:          {{Name: CircleAccessColumn_CircleId, Table: CircleAccessTable}},
	cmodel.CircleAccessField_Id:              {{Name: CircleAccessColumn_CircleAccessId, Table: CircleAccessTable}},
	cmodel.CircleAccessField_PermissionLevel: {{Name: CircleAccessColumn_PermissionLevel, Table: CircleAccessTable, Updatable: true}},
	cmodel.CircleAccessField_State:           {{Name: CircleAccessColumn_State, Table: CircleAccessTable, Updatable: true}},
	cmodel.CircleAccessField_AcceptTarget:    {{Name: CircleAccessColumn_AcceptTarget, Table: CircleAccessTable, Updatable: true}},
	cmodel.CircleAccessField_Requester: {
		{Name: CircleAccessColumn_RequesterUserId, Table: CircleAccessTable},
		{Name: CircleAccessColumn_RequesterCircleId, Table: CircleAccessTable},
	},
	cmodel.CircleAccessField_Recipient: {
		{Name: CircleAccessColumn_RecipientUserId, Table: CircleAccessTable},
		{Name: UserColumn_Username, Table: UserTable, Alias: "recipient_username"},
		{Name: UserColumn_GivenName, Table: UserTable, Alias: "recipient_given_name"},
		{Name: UserColumn_FamilyName, Table: UserTable, Alias: "recipient_family_name"},
	},
})

var CircleAccessSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"permission_level":  {Name: CircleAccessColumn_PermissionLevel, Table: CircleAccessTable},
	"state":             {Name: CircleAccessColumn_State, Table: CircleAccessTable},
	"recipient.user_id": {Name: CircleAccessColumn_RecipientUserId, Table: CircleAccessTable},
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
	AcceptTarget      types.AcceptTarget    `gorm:"not null;default:0"`

	RecipientUsername   string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName  string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (CircleAccess) TableName() string {
	return CircleAccessTable
}

// CA
type CA struct {
	CircleAccess
}

// TableName -
func (CA) TableName() string {
	return "ca"
}
