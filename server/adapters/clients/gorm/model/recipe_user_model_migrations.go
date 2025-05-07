package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// RecipeUser -
type RecipeUser struct {
	RecipeUserId    int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId        int64
	UserId          int64
	PermissionLevel pb.ShareRecipeRequest_ResourcePermission `gorm:"default:100"`
}

// TableName -
func (RecipeUser) TableName() string {
	return "recipe_user"
}
