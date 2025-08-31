package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateEventRecipe creates a new event recipe connection
func (d *Domain) CreateEventRecipe(ctx context.Context, authAccount model.AuthAccount, eventRecipe model.EventRecipe) (dbEventRecipe model.EventRecipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)
	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when creating an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	// Validate that the user has access to the calendar
	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: eventRecipe.Parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access when creating an event recipe")
		return model.EventRecipe{}, err
	}

	// Validate that the user has access to the recipe
	_, err = d.determineRecipeAccess(ctx, authAccount, eventRecipe.RecipeId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine recipe access when creating an event recipe")
		return model.EventRecipe{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin creating event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to begin creating event recipe"}
	}
	defer tx.Rollback()

	dbEventRecipe, err = tx.CreateEventRecipe(ctx, eventRecipe, []string{})
	if err != nil {
		log.Error().Err(err).Msg("unable to create event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to create event recipe"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish creating event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to finish creating event recipe"}
	}

	dbEventRecipe.Parent = eventRecipe.Parent

	return dbEventRecipe, nil
}

// DeleteEventRecipe deletes an event recipe connection
func (d *Domain) DeleteEventRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, id model.EventRecipeId) (dbEventRecipe model.EventRecipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("auth user id required when deleting an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "auth user id required"}
	}

	if id.EventRecipeId == 0 {
		log.Warn().Msg("event recipe id required when deleting an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "event recipe id required"}
	}

	// Validate that the user has access to the calendar
	_, err = d.determineCalendarAccess(ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId}, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine calendar access when deleting an event recipe")
		return model.EventRecipe{}, err
	}

	tx, err := d.repo.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to begin deleting event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to begin deleting event recipe"}
	}
	defer tx.Rollback()

	dbEventRecipe, err = tx.DeleteEventRecipe(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to delete event recipe"}
	}

	err = tx.Commit()
	if err != nil {
		log.Error().Err(err).Msg("unable to finish deleting event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to finish deleting event recipe"}
	}

	dbEventRecipe.Parent = parent

	return dbEventRecipe, nil
}

// GetEventRecipe retrieves an event recipe
func (d *Domain) GetEventRecipe(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, id model.EventRecipeId, fields []string) (dbEventRecipe model.EventRecipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when getting an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "user id required"}
	}

	if id.EventRecipeId == 0 {
		log.Warn().Msg("id required when getting an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "id required"}
	}

	if parent.CalendarId == 0 {
		log.Warn().Msg("calendar id required when getting an event recipe")
		return model.EventRecipe{}, domain.ErrInvalidArgument{Msg: "calendar id required"}
	}

	dbCalendar, err := d.repo.GetCalendar(ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId}, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar when getting an event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to get calendar"}
	}

	_, err = d.determineCalendarAccess(
		ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
		withResourceVisibilityLevel(dbCalendar.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when getting an event recipe")
		return model.EventRecipe{}, err
	}

	dbEventRecipe, err = d.repo.GetEventRecipe(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get event recipe")
		return model.EventRecipe{}, domain.ErrInternal{Msg: "unable to get event recipe"}
	}

	dbEventRecipe.Parent = parent

	return dbEventRecipe, nil
}

// ListEventRecipes lists event recipes
func (d *Domain) ListEventRecipes(ctx context.Context, authAccount model.AuthAccount, parent model.EventRecipeParent, pageSize int32, offset int64, filter string, fields []string) (dbEventRecipes []model.EventRecipe, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Warn().Msg("user id required when listing event recipes")
		return nil, domain.ErrInvalidArgument{Msg: "user id required"}
	}
	if parent.CalendarId == 0 {
		log.Warn().Msg("calendar id required when listing event recipes")
		return nil, domain.ErrInvalidArgument{Msg: "calendar id required"}
	}

	dbCalendar, err := d.repo.GetCalendar(ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId}, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get calendar when listing event recipes")
		return nil, domain.ErrInternal{Msg: "unable to get calendar"}
	}

	_, err = d.determineCalendarAccess(
		ctx, authAccount, model.CalendarId{CalendarId: parent.CalendarId},
		withResourceVisibilityLevel(dbCalendar.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_PUBLIC),
		withAllowPendingAccess(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when listing event recipes")
		return nil, err
	}

	dbEventRecipes, err = d.repo.ListEventRecipes(ctx, authAccount, parent, pageSize, offset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list event recipes")
		return nil, domain.ErrInternal{Msg: "unable to list event recipes"}
	}

	// Set the parent context for all returned event recipes
	for i := range dbEventRecipes {
		dbEventRecipes[i].Parent = parent
	}

	return dbEventRecipes, nil
}
