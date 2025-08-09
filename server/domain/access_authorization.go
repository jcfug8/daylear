package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

type accessOwnershipDetails struct {
	isRecipientOwner       bool
	isResourceOwner        bool
	maximumPermissionLevel types.PermissionLevel
	acceptTarget           types.AcceptTarget
	accessState            types.AccessState
}

type determineAccessOwnershipDetailsOption func(config *determineAccessOwnershipDetailsConfig)

func withAllowAutoOmitAccessChecks() determineAccessOwnershipDetailsOption {
	return func(config *determineAccessOwnershipDetailsConfig) {
		config.allowAutoOmitAccessChecks = true
	}
}

func withForceOmitResourceCheck() determineAccessOwnershipDetailsOption {
	return func(config *determineAccessOwnershipDetailsConfig) {
		config.forceOmitResourceCheck = true
	}
}

func withForceOmitRecipientCheck() determineAccessOwnershipDetailsOption {
	return func(config *determineAccessOwnershipDetailsConfig) {
		config.forceOmitRecipientCheck = true
	}
}

func withMinimimRecipientPermissionLevel(permissionLevel types.PermissionLevel) determineAccessOwnershipDetailsOption {
	return func(config *determineAccessOwnershipDetailsConfig) {
		config.minimimRecipientPermissionLevel = permissionLevel
	}
}

type determineAccessOwnershipDetailsConfig struct {
	determineResourceAccess         func() (model.Access, error)
	allowAutoOmitAccessChecks       bool
	forceOmitResourceCheck          bool
	forceOmitRecipientCheck         bool
	minimimRecipientPermissionLevel types.PermissionLevel
}

func (d *Domain) determineAccessOwnershipDetails(ctx context.Context, authAccount model.AuthAccount, access model.Access, options ...determineAccessOwnershipDetailsOption) (details accessOwnershipDetails, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	var determinedRecipientAccess model.Access
	config := determineAccessOwnershipDetailsConfig{
		minimimRecipientPermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_READ,
	}

	for _, option := range options {
		option(&config)
	}

	omitResourceCheck := config.forceOmitResourceCheck || (config.allowAutoOmitAccessChecks && access.GetAccessId() != 0 && access.GetAcceptTarget() == types.AcceptTarget_ACCEPT_TARGET_RECIPIENT)
	omitRecipientCheck := config.forceOmitRecipientCheck || (config.allowAutoOmitAccessChecks && access.GetAccessId() != 0 && access.GetAcceptTarget() == types.AcceptTarget_ACCEPT_TARGET_RESOURCE)

	switch a := access.(type) {
	case model.CircleAccess:
		config.determineResourceAccess = func() (model.Access, error) {
			dbCircle, err := d.repo.GetCircle(ctx, authAccount, a.CircleId, []string{model.CircleField_Visibility})
			if err != nil {
				log.Error().Err(err).Msg("unable to get circle when determining access ownership details")
				return model.CircleAccess{}, err
			}
			return d.determineCircleAccess(ctx, authAccount, a.CircleId, withResourceVisibilityLevel(dbCircle.VisibilityLevel), withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC))
		}
	case model.CalendarAccess:
		config.determineResourceAccess = func() (model.Access, error) {
			dbCalendar, err := d.repo.GetCalendar(ctx, authAccount, model.CalendarId{CalendarId: a.CalendarId}, []string{model.CalendarField_Visibility})
			if err != nil {
				log.Error().Err(err).Msg("unable to get calendar when determining access ownership details")
				return model.CalendarAccess{}, err
			}
			return d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: a.CalendarId}, withResourceVisibilityLevel(dbCalendar.VisibilityLevel), withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC))
		}
	case model.UserAccess:
		config.determineResourceAccess = func() (model.Access, error) {
			return d.determineUserAccess(ctx, authAccount, a.UserAccessParent.UserId, withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC), withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC))
		}
	case model.RecipeAccess:
		config.determineResourceAccess = func() (model.Access, error) {
			dbRecipe, err := d.repo.GetRecipe(ctx, authAccount, a.RecipeId, []string{model.RecipeField_VisibilityLevel})
			if err != nil {
				log.Error().Err(err).Msg("unable to get recipe when determining access ownership details")
				return model.RecipeAccess{}, err
			}
			return d.determineRecipeAccess(ctx, authAccount, a.RecipeId, withResourceVisibilityLevel(dbRecipe.VisibilityLevel), withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC))
		}
	default:
		log.Warn().Msg("unable to determine access ownership details for unknown access type")
		return accessOwnershipDetails{}, domain.ErrInternal{Msg: "unable to determine access ownership details for unknown access type"}
	}

	if !omitRecipientCheck {
		if access.GetRecipientCircleId().CircleId != 0 {
			dbCircle, err := d.repo.GetCircle(ctx, authAccount, access.GetRecipientCircleId(), []string{model.CircleField_Visibility})
			if err != nil {
				log.Error().Err(err).Msg("unable to get circle when creating a calendar access")
				return accessOwnershipDetails{}, err
			}
			determinedRecipientAccess, err = d.determineCircleAccess(ctx, authAccount, access.GetRecipientCircleId(), withResourceVisibilityLevel(dbCircle.VisibilityLevel), withMinimumPermissionLevel(config.minimimRecipientPermissionLevel))
			if err != nil {
				log.Error().Err(err).Msg("unable to determine circle access when creating a calendar access")
				return accessOwnershipDetails{}, err
			}
		} else if access.GetRecipientUserId().UserId != 0 {
			determinedRecipientAccess, err = d.determineUserAccess(ctx, authAccount, access.GetRecipientUserId(), withResourceVisibilityLevel(types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC), withMinimumPermissionLevel(config.minimimRecipientPermissionLevel))
			if err != nil {
				log.Error().Err(err).Msg("unable to determine user access when creating a calendar access")
				return accessOwnershipDetails{}, err
			}
		} else {
			log.Warn().Msg("no recipient provided when determining access ownership details")
			return accessOwnershipDetails{}, domain.ErrInternal{Msg: "unable to determine recipient access when determining access ownership details"}
		}

		details.isRecipientOwner = determinedRecipientAccess.GetPermissionLevel() >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN
	}

	if !omitResourceCheck {
		determinedResourceAccess, err := config.determineResourceAccess()
		if err != nil {
			log.Error().Err(err).Msg("unable to determine resource access when determining access ownership details")
			return accessOwnershipDetails{}, err
		}

		details.isResourceOwner = determinedResourceAccess.GetPermissionLevel() >= types.PermissionLevel_PERMISSION_LEVEL_ADMIN
		details.maximumPermissionLevel = max(determinedResourceAccess.GetPermissionLevel(), types.PermissionLevel_PERMISSION_LEVEL_READ)
	}

	details.acceptTarget = access.GetAcceptTarget()
	details.accessState = access.GetAccessState()
	if access.GetAccessId() == 0 {
		if !details.isResourceOwner && details.isRecipientOwner {
			details.accessState = types.AccessState_ACCESS_STATE_PENDING
			details.acceptTarget = types.AcceptTarget_ACCEPT_TARGET_RESOURCE
		} else if details.isResourceOwner && !details.isRecipientOwner {
			details.accessState = types.AccessState_ACCESS_STATE_PENDING
			details.acceptTarget = types.AcceptTarget_ACCEPT_TARGET_RECIPIENT
		} else if details.isResourceOwner && details.isRecipientOwner {
			details.accessState = types.AccessState_ACCESS_STATE_ACCEPTED
			details.acceptTarget = types.AcceptTarget_ACCEPT_TARGET_UNSPECIFIED
		} else {
			log.Warn().Msg("unable to determine access state when creating a circle access")
			return accessOwnershipDetails{}, domain.ErrInternal{Msg: "unable to determine access state"}
		}
	}
	return details, nil
}
