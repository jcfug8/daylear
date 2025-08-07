package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	RecipeAccessFields_RecipeAccessId    = "recipe_access.recipe_access_id"
	RecipeAccessFields_RecipeId          = "recipe_access.recipe_id"
	RecipeAccessFields_RequesterUserId   = "recipe_access.requester_user_id"
	RecipeAccessFields_RequesterCircleId = "recipe_access.requester_circle_id"
	RecipeAccessFields_RecipientUserId   = "recipe_access.recipient_user_id"
	RecipeAccessFields_RecipientCircleId = "recipe_access.recipient_circle_id"
	RecipeAccessFields_PermissionLevel   = "recipe_access.permission_level"
	RecipeAccessFields_State             = "recipe_access.state"
)

var RecipeAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	model.RecipeAccessField_Parent:          {RecipeAccessFields_RecipeId},
	model.RecipeAccessField_Id:              {RecipeAccessFields_RecipeAccessId},
	model.RecipeAccessField_PermissionLevel: {RecipeAccessFields_PermissionLevel},
	model.RecipeAccessField_State:           {RecipeAccessFields_State},
	model.RecipeAccessField_Requester: {
		RecipeAccessFields_RequesterUserId,
		RecipeAccessFields_RequesterCircleId,
	},
	model.RecipeAccessField_Recipient: {
		RecipeAccessFields_RecipientUserId,
		RecipeAccessFields_RecipientCircleId,
		UserColumn_Username + "as recipient_username",
		UserColumn_GivenName + "as recipient_given_name",
		UserColumn_FamilyName + "as recipient_family_name",
		CircleColumn_Title + "as recipient_circle_title",
		CircleColumn_Handle + "as recipient_circle_handle",
	},
})
var UpdateRecipeAccessFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	model.RecipeAccessField_PermissionLevel: {RecipeAccessFields_PermissionLevel},
	model.RecipeAccessField_State:           {RecipeAccessFields_State},
})
var RecipeAccessSQLConverter = filter.NewSQLConverter(map[string]string{
	model.RecipeAccessField_PermissionLevel: RecipeAccessFields_PermissionLevel,
	model.RecipeAccessField_State:           RecipeAccessFields_State,
	model.RecipeAccessField_Requester:       RecipeAccessFields_RecipientUserId,
	model.RecipeAccessField_Recipient:       RecipeAccessFields_RecipientCircleId,
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
	return "recipe_access"
}
