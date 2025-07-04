package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessFields defines the recipeAccess fields.
var RecipeAccessFields = recipeAccessFields{
	RecipeAccessId:    "recipe_access.recipe_access_id",
	RecipeId:          "recipe_access.recipe_id",
	RequesterUserId:   "recipe_access.requester_user_id",
	RequesterCircleId: "recipe_access.requester_circle_id",
	RecipientUserId:   "recipe_access.recipient_user_id",
	RecipientCircleId: "recipe_access.recipient_circle_id",
	PermissionLevel:   "recipe_access.permission_level",
	State:             "recipe_access.state",
	Title:             "recipe_access.title",
}

type recipeAccessFields struct {
	RecipeAccessId    string
	RecipeId          string
	RequesterUserId   string
	RequesterCircleId string
	RecipientUserId   string
	RecipientCircleId string
	PermissionLevel   string
	State             string
	Title             string
}

// Map maps the recipeAccess fields to their corresponding model values.
func (fields recipeAccessFields) Map(m RecipeAccess) map[string]any {
	return map[string]any{
		fields.RecipeAccessId:    m.RecipeAccessId,
		fields.RecipeId:          m.RecipeId,
		fields.RequesterUserId:   m.requesterUserId,
		fields.RequesterCircleId: m.requesterCircleId,
		fields.RecipientUserId:   m.RecipientUserId,
		fields.RecipientCircleId: m.RecipientCircleId,
		fields.PermissionLevel:   m.PermissionLevel,
		fields.State:             m.State,
		fields.Title:             m.Title,
	}
}

// Mask returns a FieldMask for the recipeAccess fields.
func (fields recipeAccessFields) Mask() []string {
	return []string{
		fields.RecipeAccessId,
		fields.RecipeId,
		fields.RequesterUserId,
		fields.RequesterCircleId,
		fields.RecipientUserId,
		fields.RecipientCircleId,
		fields.PermissionLevel,
		fields.State,
		fields.Title,
	}
}

// RecipeAccess -
type RecipeAccess struct {
	RecipeAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId          int64                 `gorm:"not null;index;uniqueIndex:idx_recipe_id_recipient_user_id;uniqueIndex:idx_recipe_id_recipient_circle_id;"`
	requesterUserId   int64                 `gorm:"index"`
	requesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"index;uniqueIndex:idx_recipe_id_recipient_user_id,where:recipient_user_id <> null"`
	RecipientCircleId int64                 `gorm:"index;uniqueIndex:idx_recipe_id_recipient_circle_id,where:recipient_circle_id <> null"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`
	Title             string                `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (RecipeAccess) TableName() string {
	return "recipe_access"
}
