package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	RecipeAccessTable = "recipe_access"
)

const (
	RecipeAccessFields_RecipeAccessId    = "recipe_access_id"
	RecipeAccessFields_RecipeId          = "recipe_id"
	RecipeAccessFields_RequesterUserId   = "requester_user_id"
	RecipeAccessFields_RequesterCircleId = "requester_circle_id"
	RecipeAccessFields_RecipientUserId   = "recipient_user_id"
	RecipeAccessFields_RecipientCircleId = "recipient_circle_id"
	RecipeAccessFields_PermissionLevel   = "permission_level"
	RecipeAccessFields_State             = "state"
)

var RecipeAccessFieldMasker = fieldmask.NewSQLFieldMasker(RecipeAccess{}, map[string][]fieldmask.Field{
	model.RecipeAccessField_Parent:          {{Name: RecipeAccessFields_RecipeId, Table: RecipeAccessTable}},
	model.RecipeAccessField_Id:              {{Name: RecipeAccessFields_RecipeAccessId, Table: RecipeAccessTable}},
	model.RecipeAccessField_PermissionLevel: {{Name: RecipeAccessFields_PermissionLevel, Table: RecipeAccessTable, Updatable: true}},
	model.RecipeAccessField_State:           {{Name: RecipeAccessFields_State, Table: RecipeAccessTable, Updatable: true}},
	model.RecipeAccessField_Requester: {
		{Name: RecipeAccessFields_RequesterUserId, Table: RecipeAccessTable},
		{Name: RecipeAccessFields_RequesterCircleId, Table: RecipeAccessTable},
	},
	model.RecipeAccessField_Recipient: {
		{Name: RecipeAccessFields_RecipientUserId, Table: RecipeAccessTable},
		{Name: RecipeAccessFields_RecipientCircleId, Table: RecipeAccessTable},
		{Name: UserColumn_Username, Table: UserTable, Alias: "recipient_username"},
		{Name: UserColumn_GivenName, Table: UserTable, Alias: "recipient_given_name"},
		{Name: UserColumn_FamilyName, Table: UserTable, Alias: "recipient_family_name"},
		{Name: CircleColumn_Title, Table: CircleTable, Alias: "recipient_circle_title"},
		{Name: CircleColumn_Handle, Table: CircleTable, Alias: "recipient_circle_handle"},
	},
})

var RecipeAccessSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	model.RecipeAccessField_PermissionLevel: {Name: RecipeAccessFields_PermissionLevel, Table: RecipeAccessTable},
	model.RecipeAccessField_State:           {Name: RecipeAccessFields_State, Table: RecipeAccessTable},
	model.RecipeAccessField_Requester:       {Name: RecipeAccessFields_RecipientUserId, Table: RecipeAccessTable},
	model.RecipeAccessField_Recipient:       {Name: RecipeAccessFields_RecipientCircleId, Table: RecipeAccessTable},
}, true)

// RecipeAccess -
type RecipeAccess struct {
	RecipeAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId          int64                 `gorm:"not null;index;uniqueIndex:idx_recipe_id_recipient_user_id;uniqueIndex:idx_recipe_id_recipient_circle_id;"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"index;uniqueIndex:idx_recipe_id_recipient_user_id,where:recipient_user_id <> null"`
	RecipientCircleId int64                 `gorm:"index;uniqueIndex:idx_recipe_id_recipient_circle_id,where:recipient_circle_id <> null"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`

	RecipientUsername     string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName    string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName   string `gorm:"->;-:migration"` // read only from join
	RecipientCircleTitle  string `gorm:"->;-:migration"` // read only from join
	RecipientCircleHandle string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (RecipeAccess) TableName() string {
	return RecipeAccessTable
}
