package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) checkRecipeAccess(ctx context.Context, authAccount model.AuthAccount, recipeId model.RecipeId, minPermLevel types.PermissionLevel) (permissionLevel types.PermissionLevel, visibilityLevel types.VisibilityLevel, err error) {
	permissionLevel, visibilityLevel, err = d.getRecipeAccessLevels(ctx, authAccount, recipeId)
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
	}

	if authAccount.CircleId != 0 {
		circlePermissionLevel, circleVisibilityLevel, err := d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
		}

		permissionLevel, visibilityLevel, err = determineRecipeAccessLevels(circleVisibilityLevel, circlePermissionLevel, visibilityLevel, permissionLevel)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
		}
	}

	if authAccount.UserId != 0 {
		userPermissionLevel, err := d.getUserAccessLevels(ctx, authAccount)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
		}

		permissionLevel, visibilityLevel, err = determineRecipeAccessLevels(visibilityLevel, userPermissionLevel, visibilityLevel, permissionLevel)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
		}
	}

	// verify requester has the required permission level
	if permissionLevel < minPermLevel {
		return 0, 0, domain.ErrPermissionDenied{Msg: "user does not have the correct permission level"}
	}

	return permissionLevel, visibilityLevel, nil
}

func (d *Domain) checkUserAccess(ctx context.Context, authAccount model.AuthAccount, userId model.UserId, minPermLevel types.PermissionLevel) (permissionLevel types.PermissionLevel, err error) {
	permissionLevel, err = d.getUserAccessLevels(ctx, authAccount)
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	if authAccount.CircleId != 0 {
		circlePermissionLevel, _, err := d.getCircleAccessLevels(ctx, authAccount)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
		}

		permissionLevel, err = determineUserAccessLevels(circlePermissionLevel, permissionLevel)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
		}
	}

	if authAccount.UserId != 0 {
		userPermissionLevel, err := d.getUserAccessLevels(ctx, authAccount)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
		}

		permissionLevel, err = determineUserAccessLevels(userPermissionLevel, permissionLevel)
		if err != nil {
			return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
		}
	}

	return permissionLevel, nil
}

func (d *Domain) getUserAccessLevels(ctx context.Context, authAccount model.AuthAccount) (types.PermissionLevel, error) {
	// verify auth account is set
	if authAccount.AuthUserId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify user id is set
	if authAccount.UserId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	// verify user exists
	user, err := d.repo.GetUser(ctx, authAccount, model.UserId{UserId: authAccount.UserId})
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	if user.Id.UserId == authAccount.AuthUserId {
		return types.PermissionLevel_PERMISSION_LEVEL_ADMIN, nil
	}

	return user.UserAccess.Level, nil
}

func (d *Domain) getCircleAccessLevels(ctx context.Context, authAccount model.AuthAccount) (types.PermissionLevel, types.VisibilityLevel, error) {
	// verify auth account is set
	if authAccount.AuthUserId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle id is set
	if authAccount.CircleId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify circle exists
	circle, err := d.repo.GetCircle(ctx, authAccount, model.CircleId{CircleId: authAccount.CircleId})
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
	}

	if circle.VisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED && circle.CircleAccess.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to circle"}
	}

	return circle.CircleAccess.PermissionLevel, circle.VisibilityLevel, nil
}

func (d *Domain) getRecipeAccessLevels(ctx context.Context, authAccount model.AuthAccount, recipeId model.RecipeId) (types.PermissionLevel, types.VisibilityLevel, error) {
	// verify auth account is set
	if authAccount.AuthUserId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle id is set
	if recipeId.RecipeId == 0 {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify circle exists
	recipe, err := d.repo.GetRecipe(ctx, authAccount, recipeId)
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
	}

	if recipe.Visibility == types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED && recipe.RecipeAccess.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to recipe"}
	}

	if recipe.Visibility == types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN && recipe.RecipeAccess.PermissionLevel != types.PermissionLevel_PERMISSION_LEVEL_ADMIN {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to recipe"}
	}

	return recipe.RecipeAccess.PermissionLevel, recipe.Visibility, nil
}

func determineRecipeAccessLevels(circleVisibilityLevel types.VisibilityLevel, circlePermissionLevel types.PermissionLevel, recipeVisibilityLevel types.VisibilityLevel, recipePermissionLevel types.PermissionLevel) (types.PermissionLevel, types.VisibilityLevel, error) {
	// If either access level is unspecified, no access
	if circlePermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED &&
		recipePermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to recipe: access not set"}
	}

	// Effective permission is minimum of circle and recipe permissions
	// User cannot have higher access to recipe than either their circle access allows
	// or the circle's access to the recipe allows
	effectivePermission := circlePermissionLevel
	if recipePermissionLevel < circlePermissionLevel {
		effectivePermission = recipePermissionLevel
	}

	// Determine effective visibility
	effectiveVisibility := recipeVisibilityLevel

	// If user only has PUBLIC permission, they can only see PUBLIC recipes
	if effectivePermission == types.PermissionLevel_PERMISSION_LEVEL_PUBLIC &&
		recipeVisibilityLevel != types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to recipe: public access not allowed"}
	}

	return effectivePermission, effectiveVisibility, nil
}

func determineUserAccessLevels(circlePermissionLevel types.PermissionLevel, userPermissionLevel types.PermissionLevel) (types.PermissionLevel, error) {
	// // If either access level is unspecified, no access
	// if circlePermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED &&
	// 	userPermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
	// 	return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to user: access not set"}
	// }

	// Effective permission is minimum of circle and user permissions
	// User cannot have higher access to user than either their circle access allows
	// or the circle's access to the user allows
	effectivePermission := circlePermissionLevel
	if userPermissionLevel < circlePermissionLevel {
		effectivePermission = userPermissionLevel
	}

	return effectivePermission, nil
}
