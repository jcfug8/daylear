package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessFields defines the recipeAccess fields.
var RecipeAccessFields = recipeAccessFields{
	RecipeAccessId:  "recipe_access.recipe_access_id",
	RecipeId:        "recipe_access.recipe_id",
	UserId:          "recipe_access.user_id",
	CircleId:        "recipe_access.circle_id",
	PermissionLevel: "recipe_access.permission_level",
	State:           "recipe_access.state",
	Title:           "recipe_access.title",
}

type recipeAccessFields struct {
	RecipeAccessId  string
	RecipeId        string
	UserId          string
	CircleId        string
	PermissionLevel string
	State           string
	Title           string
}

// Map maps the recipeAccess fields to their corresponding model values.
func (fields recipeAccessFields) Map(m RecipeAccess) map[string]any {
	return map[string]any{
		fields.RecipeAccessId:  m.RecipeAccessId,
		fields.RecipeId:        m.RecipeId,
		fields.UserId:          m.UserId,
		fields.CircleId:        m.CircleId,
		fields.PermissionLevel: m.PermissionLevel,
		fields.State:           m.State,
		fields.Title:           m.Title,
	}
}

// Mask returns a FieldMask for the recipeAccess fields.
func (fields recipeAccessFields) Mask() []string {
	return []string{
		fields.RecipeAccessId,
		fields.RecipeId,
		fields.UserId,
		fields.CircleId,
		fields.PermissionLevel,
		fields.State,
		fields.Title,
	}
}

// RecipeAccess -
type RecipeAccess struct {
	RecipeAccessId  int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId        int64                  `gorm:"not null;index"`
	UserId          int64                  `gorm:"not null;index"`
	CircleId        int64                  `gorm:"not null;index"`
	PermissionLevel permPb.PermissionLevel `gorm:"not null"`
	State           pb.Access_State        `gorm:"not null"`
	Title           string                 `gorm:"->"` // read only from join
}

// TableName -
func (RecipeAccess) TableName() string {
	return "recipe_access"
}
