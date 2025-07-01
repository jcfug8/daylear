package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) getCircleAccessLevels(ctx context.Context, authAccount model.AuthAccount) (types.PermissionLevel, types.VisibilityLevel, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
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

	if circle.VisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED || circle.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to circle"}
	}

	return circle.PermissionLevel, circle.VisibilityLevel, nil
}

func (d *Domain) getRecipeAccessLevels(ctx context.Context, authAccount model.AuthAccount, recipeId model.RecipeId) (types.PermissionLevel, types.VisibilityLevel, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
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

	if recipe.Visibility == types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED || recipe.Permission == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user does not have access to recipe"}
	}

	return recipe.Permission, recipe.Visibility, nil
}

func (d *Domain) getRecipeAccessLevelsForCircle(ctx context.Context, authAccount model.AuthAccount, recipeId model.RecipeId) (types.PermissionLevel, types.VisibilityLevel, error) {
	circlePermissionLevel, circleVisibilityLevel, err := d.getCircleAccessLevels(ctx, authAccount)
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
	}

	recipePermissionLevel, recipeVisibilityLevel, err := d.getRecipeAccessLevels(ctx, authAccount, recipeId)
	if err != nil {
		return types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, types.VisibilityLevel_VISIBILITY_LEVEL_UNSPECIFIED, err
	}

	return determineRecipeAccessLevels(circleVisibilityLevel, circlePermissionLevel, recipeVisibilityLevel, recipePermissionLevel)
}

func determineRecipeAccessLevels(circleVisibilityLevel types.VisibilityLevel, circlePermissionLevel types.PermissionLevel, recipeVisibilityLevel types.VisibilityLevel, recipePermissionLevel types.PermissionLevel) (types.PermissionLevel, types.VisibilityLevel, error) {
	// If either access level is unspecified, no access
	if circlePermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED ||
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
