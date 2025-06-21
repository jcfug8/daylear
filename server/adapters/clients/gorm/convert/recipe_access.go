package convert

import (
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	coreModel "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeAccessLevelToPermissionLevel converts a recipe access level to a permission level.
func RecipeAccessLevelToPermissionLevel(level pb.Access_Level) permPb.PermissionLevel {
	switch level {
	case pb.Access_LEVEL_READ:
		return permPb.PermissionLevel_RESOURCE_PERMISSION_READ
	case pb.Access_LEVEL_WRITE:
		return permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE
	case pb.Access_LEVEL_ADMIN:
		return permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE // TODO: change to admin
	default:
		return permPb.PermissionLevel_RESOURCE_PERMISSION_UNSPECIFIED
	}
}

// PermissionLevelToRecipeAccessLevel converts a permission level to a recipe access level.
func PermissionLevelToRecipeAccessLevel(level permPb.PermissionLevel) pb.Access_Level {
	switch level {
	case permPb.PermissionLevel_RESOURCE_PERMISSION_READ:
		return pb.Access_LEVEL_READ
	case permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE:
		return pb.Access_LEVEL_WRITE
	default:
		return pb.Access_LEVEL_UNSPECIFIED
	}
}

// CoreRecipeAccessToRecipeUser converts a core RecipeAccess model to a gorm RecipeUser model.
func CoreRecipeAccessToRecipeUser(access coreModel.RecipeAccess) dbModel.RecipeUser {
	return dbModel.RecipeUser{
		RecipeUserId:    access.RecipeAccessId.RecipeAccessId,
		RecipeId:        access.RecipeId.RecipeId,
		UserId:          access.Recipient.UserId,
		PermissionLevel: RecipeAccessLevelToPermissionLevel(access.Level),
		State:           access.State,
	}
}

// CoreRecipeAccessToRecipeCircle converts a core RecipeAccess model to a gorm RecipeCircle model.
func CoreRecipeAccessToRecipeCircle(access coreModel.RecipeAccess) dbModel.RecipeCircle {
	return dbModel.RecipeCircle{
		RecipeCircleId:  access.RecipeAccessId.RecipeAccessId,
		RecipeId:        access.RecipeId.RecipeId,
		CircleId:        access.Recipient.CircleId,
		PermissionLevel: RecipeAccessLevelToPermissionLevel(access.Level),
		State:           access.State,
	}
}

// RecipeUserToCoreRecipeAccess converts a gorm RecipeUser model to a core RecipeAccess model.
func RecipeUserToCoreRecipeAccess(user dbModel.RecipeUser) coreModel.RecipeAccess {
	return coreModel.RecipeAccess{
		RecipeAccessId: coreModel.RecipeAccessId{
			RecipeAccessId: user.RecipeUserId,
		},
		RecipeAccessParent: coreModel.RecipeAccessParent{
			RecipeId: coreModel.RecipeId{
				RecipeId: user.RecipeId,
			},
			Recipient: coreModel.RecipeParent{
				UserId: user.UserId,
			},
		},
		Level: PermissionLevelToRecipeAccessLevel(user.PermissionLevel),
		State: user.State,
	}
}

// RecipeCircleToCoreRecipeAccess converts a gorm RecipeCircle model to a core RecipeAccess model.
func RecipeCircleToCoreRecipeAccess(circle dbModel.RecipeCircle) coreModel.RecipeAccess {
	return coreModel.RecipeAccess{
		RecipeAccessId: coreModel.RecipeAccessId{
			RecipeAccessId: circle.RecipeCircleId,
		},
		RecipeAccessParent: coreModel.RecipeAccessParent{
			RecipeId: coreModel.RecipeId{
				RecipeId: circle.RecipeId,
			},
			Recipient: coreModel.RecipeParent{
				CircleId: circle.CircleId,
			},
		},
		Level: PermissionLevelToRecipeAccessLevel(circle.PermissionLevel),
		State: circle.State,
	}
}
