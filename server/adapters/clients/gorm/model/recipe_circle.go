package model

import permPb "github.com/jcfug8/daylear/server/genapi/api/types"

// RecipeCircleFields defines the recipeCircle fields.
var RecipeCircleFields = recipeCircleFields{
	RecipeCircleId:  "recipe_circle_id",
	RecipeId:        "recipe_id",
	CircleId:        "circle_id",
	PermissionLevel: "permission_level",
	Title:           "title",
}

type recipeCircleFields struct {
	RecipeCircleId  string
	RecipeId        string
	CircleId        string
	PermissionLevel string
	Title           string
}

// Map maps the recipeCircle fields to their corresponding model values.
func (fields recipeCircleFields) Map(m RecipeCircle) map[string]any {
	return map[string]any{
		fields.RecipeCircleId:  m.RecipeCircleId,
		fields.RecipeId:        m.RecipeId,
		fields.CircleId:        m.CircleId,
		fields.PermissionLevel: m.PermissionLevel,
		fields.Title:           m.Title,
	}
}

// Mask returns a FieldMask for the recipeCircle fields.
func (fields recipeCircleFields) Mask() []string {
	return []string{
		fields.RecipeCircleId,
		fields.RecipeId,
		fields.CircleId,
		fields.PermissionLevel,
		fields.Title,
	}
}

// RecipeCircle -
type RecipeCircle struct {
	RecipeCircleId  int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId        int64                  `gorm:"not null;index"`
	CircleId        int64                  `gorm:"not null;index"`
	PermissionLevel permPb.PermissionLevel `gorm:"default:100"`
	Title           string                 `gorm:"->"` // read only from join
}

// TableName -
func (RecipeCircle) TableName() string {
	return "recipe_circle"
}
