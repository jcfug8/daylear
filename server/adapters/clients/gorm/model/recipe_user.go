package model

import permPb "github.com/jcfug8/daylear/server/genapi/api/types"

// RecipeUserFields defines the recipeUser fields.
var RecipeUserFields = recipeUserFields{
	RecipeUserId:    "recipe_user_id",
	RecipeId:        "recipe_id",
	UserId:          "user_id",
	PermissionLevel: "permission_level",
}

type recipeUserFields struct {
	RecipeUserId    string
	RecipeId        string
	UserId          string
	PermissionLevel string
}

// Map maps the recipeUser fields to their corresponding model values.
func (fields recipeUserFields) Map(m RecipeUser) map[string]any {
	return map[string]any{
		fields.RecipeUserId:    m.RecipeUserId,
		fields.RecipeId:        m.RecipeId,
		fields.UserId:          m.UserId,
		fields.PermissionLevel: m.PermissionLevel,
	}
}

// Mask returns a FieldMask for the recipeUser fields.
func (fields recipeUserFields) Mask() []string {
	return []string{
		fields.RecipeUserId,
		fields.RecipeId,
		fields.UserId,
		fields.PermissionLevel,
	}
}

// RecipeUser -
type RecipeUser struct {
	RecipeUserId    int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId        int64
	UserId          int64
	PermissionLevel permPb.PermissionLevel `gorm:"default:100"`
}

// TableName -
func (RecipeUser) TableName() string {
	return "recipe_user"
}
