package domain

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/core/logutil"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
	"github.com/jcfug8/daylear/server/ports/repository"
)

type determinAccessConfig struct {
	resourceVisibilityLevel   types.VisibilityLevel
	minimumPermissionLevel    types.PermissionLevel
	findStandardUserAccess    func() (model.Access, error)
	findDelegatedCircleAccess func() (model.Access, model.CircleAccess, error)
	findDelegatedUserAccess   func() (model.Access, model.UserAccess, error)
	allowPendingAccess        bool
}

type determineAccessOption func(config *determinAccessConfig)

func withResourceVisibilityLevel(visibilityLevel types.VisibilityLevel) determineAccessOption {
	return func(config *determinAccessConfig) {
		config.resourceVisibilityLevel = visibilityLevel
	}
}

func withMinimumPermissionLevel(permissionLevel types.PermissionLevel) determineAccessOption {
	return func(config *determinAccessConfig) {
		config.minimumPermissionLevel = permissionLevel
	}
}

func withAllowPendingAccess() determineAccessOption {
	return func(config *determinAccessConfig) {
		config.allowPendingAccess = true
	}
}

// determineAccess - determines the access a user has to a resource. It takes into account the user's standard access,
// their delegated circle accesses, and their delegated user accesses. The access that is returned is the one with the
// highest permission level.
//
// If the access returned is a delegated circle access, the permission level is the minimum of the circle's access to the calendar
// and the user's access to the circle. If the access returned is a delegated user access, the permission level alway set to READ.
//
// Different options can be provided to the function to customize the access level.
//
// withResourceVisibilityLevel - sets the visibility level of the resource:
//   - VISIBILITY_LEVEL_HIDDEN: Will only consider standard user access.
//   - VISIBILITY_LEVEL_PRIVATE: Will consider standard user access and delegated circle accesses.
//   - VISIBILITY_LEVEL_RESTRICTED: Will consider standard user access, delegated circle accesses, and delegated user accesses.
//   - VISIBILITY_LEVEL_PUBLIC: Will consider all accesses, but will return PUBLIC if no explicit access is found.
//   - VISIBILITY_LEVEL_UNSPECIFIED (default): Will consider all accesses, but will return UNSPECIFIED if no explicit access is found.
//
// withMinimumPermissionLevel - sets the minimum permission level required to access the resource. If the determined access
// level is below the minimum permission level, an error is returned.
func (d *Domain) determineAccess(ctx context.Context, authAccount model.AuthAccount, id model.ResourceId, options ...determineAccessOption) (access model.Access, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	config := determinAccessConfig{}

	for _, option := range options {
		option(&config)
	}

	var standardUserAccess model.Access
	var delegatedCircleAccess model.Access
	var delegatedUserAccess model.Access
	var determinedAccess model.Access

	// based on the id type, set the appropriate find functions
	switch i := id.(type) {
	case model.CalendarId:
		config.findStandardUserAccess = func() (model.Access, error) {
			return d.repo.FindStandardUserCalendarAccess(ctx, authAccount, i)
		}
		config.findDelegatedCircleAccess = func() (model.Access, model.CircleAccess, error) {
			return d.repo.FindDelegatedCircleCalendarAccess(ctx, authAccount, i)
		}
		config.findDelegatedUserAccess = func() (model.Access, model.UserAccess, error) {
			return d.repo.FindDelegatedUserCalendarAccess(ctx, authAccount, i)
		}
		determinedAccess = model.CalendarAccess{}
	case model.CircleId:
		config.findStandardUserAccess = func() (model.Access, error) {
			return d.repo.FindStandardUserCircleAccess(ctx, authAccount, i)
		}
		config.findDelegatedCircleAccess = func() (model.Access, model.CircleAccess, error) {
			return model.CircleAccess{}, model.CircleAccess{}, nil
		}
		config.findDelegatedUserAccess = func() (model.Access, model.UserAccess, error) {
			return d.repo.FindDelegatedUserCircleAccess(ctx, authAccount, i)
		}
		determinedAccess = model.CircleAccess{}
	case model.RecipeId:
		config.findStandardUserAccess = func() (model.Access, error) {
			return d.repo.FindStandardUserRecipeAccess(ctx, authAccount, i)
		}
		config.findDelegatedCircleAccess = func() (model.Access, model.CircleAccess, error) {
			return d.repo.FindDelegatedCircleRecipeAccess(ctx, authAccount, i)
		}
		config.findDelegatedUserAccess = func() (model.Access, model.UserAccess, error) {
			return d.repo.FindDelegatedUserRecipeAccess(ctx, authAccount, i)
		}
		determinedAccess = model.RecipeAccess{}
	case model.ListId:
		config.findStandardUserAccess = func() (model.Access, error) {
			return d.repo.FindStandardUserListAccess(ctx, authAccount, i)
		}
		config.findDelegatedCircleAccess = func() (model.Access, model.CircleAccess, error) {
			return d.repo.FindDelegatedCircleListAccess(ctx, authAccount, i)
		}
		config.findDelegatedUserAccess = func() (model.Access, model.UserAccess, error) {
			return d.repo.FindDelegatedUserListAccess(ctx, authAccount, i)
		}
		determinedAccess = model.ListAccess{}
	case model.UserId:
		config.findStandardUserAccess = func() (model.Access, error) {
			// if the user is requesting access to their own user, return the admin access
			if authAccount.AuthUserId == i.UserId {
				return model.UserAccess{
					UserAccessParent: model.UserAccessParent{
						UserId: model.UserId{UserId: authAccount.UserId},
					},
					PermissionLevel: types.PermissionLevel_PERMISSION_LEVEL_ADMIN,
					State:           types.AccessState_ACCESS_STATE_ACCEPTED,
					Requester:       model.UserId{UserId: authAccount.UserId},
					Recipient:       model.UserId{UserId: authAccount.UserId},
				}, nil
			}
			return d.repo.FindStandardUserUserAccess(ctx, authAccount, i)
		}
		config.findDelegatedCircleAccess = func() (model.Access, model.CircleAccess, error) {
			return model.UserAccess{}, model.CircleAccess{}, nil
		}
		config.findDelegatedUserAccess = func() (model.Access, model.UserAccess, error) {
			return model.UserAccess{}, model.UserAccess{}, nil
		}
		determinedAccess = model.UserAccess{}
	default:
		return nil, domain.ErrInvalidArgument{Msg: "invalid resource id"}
	}

	// if the resource visibility level is hidden or lower, only consider standard user access
	if config.resourceVisibilityLevel <= types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN {
		standardUserAccess, err = config.findStandardUserAccess()
		if err != nil && !errors.Is(err, repository.ErrNotFound{}) {
			log.Error().Err(err).Msg("error finding standard user resource access")
			return nil, errors.New("unable to determine resource access")
		}
		if standardUserAccess.GetPermissionLevel() > types.PermissionLevel_PERMISSION_LEVEL_PUBLIC && (config.allowPendingAccess || standardUserAccess.GetAccessState() == types.AccessState_ACCESS_STATE_ACCEPTED) {
			determinedAccess = standardUserAccess
		}
	}
	// if the resource visibility level is private or lower, consider standard user access and delegated circle accesses
	if config.resourceVisibilityLevel <= types.VisibilityLevel_VISIBILITY_LEVEL_PRIVATE {
		var circleAccess model.CircleAccess
		delegatedCircleAccess, circleAccess, err = config.findDelegatedCircleAccess()
		if err != nil && !errors.Is(err, repository.ErrNotFound{}) {
			log.Error().Err(err).Msg("error finding delegated circle resource access")
			return nil, errors.New("unable to determine resource access")
		}
		effectivePermissionLevel := min(delegatedCircleAccess.GetPermissionLevel(), circleAccess.PermissionLevel)
		if effectivePermissionLevel > determinedAccess.GetPermissionLevel() && (config.allowPendingAccess || delegatedCircleAccess.GetAccessState() == types.AccessState_ACCESS_STATE_ACCEPTED) {
			determinedAccess = delegatedCircleAccess.SetPermissionLevel(effectivePermissionLevel)
		}
	}
	// if the resource visibility level is restricted or lower, consider standard user access, delegated circle accesses, and delegated user accesses
	if config.resourceVisibilityLevel <= types.VisibilityLevel_VISIBILITY_LEVEL_RESTRICTED {
		delegatedUserAccess, _, err = config.findDelegatedUserAccess()
		if err != nil && !errors.Is(err, repository.ErrNotFound{}) {
			log.Error().Err(err).Msg("error finding delegated user resource access")
			return nil, errors.New("unable to determine resource access")
		}
		effectivePermissionLevel := min(delegatedUserAccess.GetPermissionLevel(), types.PermissionLevel_PERMISSION_LEVEL_READ)
		if effectivePermissionLevel > determinedAccess.GetPermissionLevel() && (config.allowPendingAccess || delegatedUserAccess.GetAccessState() == types.AccessState_ACCESS_STATE_ACCEPTED) {
			determinedAccess = delegatedUserAccess.SetPermissionLevel(effectivePermissionLevel)
		}
	}
	// if the resource visibility level is public, set the permission level to public if it is unspecified
	if config.resourceVisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC &&
		determinedAccess.GetPermissionLevel() == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		determinedAccess = determinedAccess.SetPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC)
	}

	// based on the minimum permission level, return an error if the determined access level is below the minimum
	if determinedAccess.GetPermissionLevel() < config.minimumPermissionLevel {
		return nil, domain.ErrPermissionDenied{Msg: "permission level too low"}
	}
	// if the determined access level is unspecified, which means no access was found, return an error
	if determinedAccess.GetPermissionLevel() == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		return nil, domain.ErrPermissionDenied{Msg: "no access"}
	}

	return determinedAccess, nil
}

// determineRecipeAccess - convenience function to determine the recipe access for a given recipe id that wraps the determineAccess function
func (d *Domain) determineRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId, options ...determineAccessOption) (recipeAccess model.RecipeAccess, err error) {
	access, err := d.determineAccess(ctx, authAccount, id, options...)
	if err != nil {
		return model.RecipeAccess{}, err
	}
	return access.(model.RecipeAccess), nil
}

// determineCalendarAccess - convenience function to determine the calendar access for a given calendar id that wraps the determineAccess function
func (d *Domain) determineCalendarAccess(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId, options ...determineAccessOption) (calendarAccess model.CalendarAccess, err error) {
	access, err := d.determineAccess(ctx, authAccount, id, options...)
	if err != nil {
		return model.CalendarAccess{}, err
	}
	return access.(model.CalendarAccess), nil
}

// determineCircleAccess - convenience function to determine the circle access for a given circle id that wraps the determineAccess function
func (d *Domain) determineCircleAccess(ctx context.Context, authAccount model.AuthAccount, id model.CircleId, options ...determineAccessOption) (circleAccess model.CircleAccess, err error) {
	access, err := d.determineAccess(ctx, authAccount, id, options...)
	if err != nil {
		return model.CircleAccess{}, err
	}
	return access.(model.CircleAccess), nil
}

// determineUserAccess - convenience function to determine the user access for a given user id that wraps the determineAccess function
func (d *Domain) determineUserAccess(ctx context.Context, authAccount model.AuthAccount, id model.UserId, options ...determineAccessOption) (userAccess model.UserAccess, err error) {
	access, err := d.determineAccess(ctx, authAccount, id, options...)
	if err != nil {
		return model.UserAccess{}, err
	}
	return access.(model.UserAccess), nil
}

// determineListAccess - convenience function to determine the list access for a given list id that wraps the determineAccess function
func (d *Domain) determineListAccess(ctx context.Context, authAccount model.AuthAccount, listId model.ListId, options ...determineAccessOption) (model.ListAccess, error) {
	access, err := d.determineAccess(ctx, authAccount, listId, options...)
	if err != nil {
		return model.ListAccess{}, err
	}
	return access.(model.ListAccess), nil
}
