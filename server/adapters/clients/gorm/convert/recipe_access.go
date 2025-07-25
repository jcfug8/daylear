package convert

import (
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	coreModel "github.com/jcfug8/daylear/server/core/model"
)

// CoreRecipeAccessToRecipeAccess converts a core RecipeAccess model to a gorm RecipeAccess model.
func CoreRecipeAccessToRecipeAccess(access coreModel.RecipeAccess) dbModel.RecipeAccess {
	return dbModel.RecipeAccess{
		RecipeAccessId:    access.RecipeAccessId.RecipeAccessId,
		RecipeId:          access.RecipeAccessParent.RecipeId.RecipeId,
		RecipientUserId:   access.Recipient.UserId,
		RecipientCircleId: access.Recipient.CircleId,
		PermissionLevel:   access.PermissionLevel,
		State:             access.State,
	}
}

// RecipeAccessToCoreRecipeAccess converts a gorm RecipeAccess model to a core RecipeAccess model.
func RecipeAccessToCoreRecipeAccess(dbAccess dbModel.RecipeAccess) coreModel.RecipeAccess {
	return coreModel.RecipeAccess{
		RecipeAccessId: coreModel.RecipeAccessId{
			RecipeAccessId: dbAccess.RecipeAccessId,
		},
		RecipeAccessParent: coreModel.RecipeAccessParent{
			RecipeId: coreModel.RecipeId{
				RecipeId: dbAccess.RecipeId,
			},
		},
		Recipient: coreModel.RecipeRecipientOrRequester{
			UserId:   dbAccess.RecipientUserId,
			CircleId: dbAccess.RecipientCircleId,
		},
		PermissionLevel:       dbAccess.PermissionLevel,
		State:                 dbAccess.State,
		RecipientUsername:     dbAccess.RecipientUsername,
		RecipientGivenName:    dbAccess.RecipientGivenName,
		RecipientFamilyName:   dbAccess.RecipientFamilyName,
		RecipientCircleTitle:  dbAccess.RecipientCircleTitle,
		RecipientCircleHandle: dbAccess.RecipientCircleHandle,
	}
}
