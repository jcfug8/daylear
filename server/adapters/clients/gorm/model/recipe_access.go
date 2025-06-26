package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessFields defines the recipeAccess fields.
var RecipeAccessFields = recipeAccessFields{
	RecipeAccessId:    "recipe_access.recipe_access_id",
	RecipeId:          "recipe_access.recipe_id",
	requesterUserId:   "recipe_access.requester_user_id",
	requesterCircleId: "recipe_access.requester_circle_id",
	RecipientUserId:   "recipe_access.recipient_user_id",
	RecipientCircleId: "recipe_access.recipient_circle_id",
	PermissionLevel:   "recipe_access.permission_level",
	State:             "recipe_access.state",
	Title:             "recipe_access.title",
}

type recipeAccessFields struct {
	RecipeAccessId    string
	RecipeId          string
	requesterUserId   string
	requesterCircleId string
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
		fields.requesterUserId:   m.requesterUserId,
		fields.requesterCircleId: m.requesterCircleId,
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
		fields.requesterUserId,
		fields.requesterCircleId,
		fields.RecipientUserId,
		fields.RecipientCircleId,
		fields.PermissionLevel,
		fields.State,
		fields.Title,
	}
}

// RecipeAccess -
type RecipeAccess struct {
	RecipeAccessId    int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId          int64                  `gorm:"not null;index"`
	requesterUserId   int64                  `gorm:"index"`
	requesterCircleId int64                  `gorm:"index"`
	RecipientUserId   int64                  `gorm:"index"`
	RecipientCircleId int64                  `gorm:"index"`
	PermissionLevel   permPb.PermissionLevel `gorm:"not null"`
	State             pb.Access_State        `gorm:"not null"`
	Title             string                 `gorm:"->"` // read only from join
}

// TableName -
func (RecipeAccess) TableName() string {
	return "recipe_access"
}
