package model

import (
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

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
