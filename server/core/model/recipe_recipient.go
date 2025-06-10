package model

import (
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

type RecipeRecipient struct {
	RecipeId
	RecipeParent
	Title           string
	PermissionLevel permPb.PermissionLevel
}
