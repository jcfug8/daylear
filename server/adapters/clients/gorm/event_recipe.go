package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateEventRecipe creates a new event recipe
func (repo *Client) CreateEventRecipe(ctx context.Context, m cmodel.EventRecipe, fields []string) (cmodel.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.EventRecipeFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid event recipe when creating event recipe row")
		return cmodel.EventRecipe{}, repository.ErrInvalidArgument{Msg: "invalid event recipe"}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.EventRecipeFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create event recipe row")
		return cmodel.EventRecipe{}, ConvertGormError(err)
	}

	m, err = convert.EventRecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid event recipe row when creating event recipe")
		return cmodel.EventRecipe{}, repository.ErrInternal{Msg: "invalid event recipe row when creating event recipe"}
	}

	return m, nil
}

// DeleteEventRecipe deletes an event recipe
func (repo *Client) DeleteEventRecipe(ctx context.Context, id cmodel.EventRecipeId) (cmodel.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("eventRecipeId", id.EventRecipeId).
		Logger()

	gm := gmodel.EventRecipe{EventRecipeId: id.EventRecipeId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.EventRecipeFieldMasker.Get()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete event recipe row")
		return cmodel.EventRecipe{}, ConvertGormError(err)
	}

	m, err := convert.EventRecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid event recipe row when deleting event recipe")
		return cmodel.EventRecipe{}, repository.ErrInternal{Msg: "invalid event recipe row when deleting event recipe"}
	}

	return m, nil
}

// GetEventRecipe retrieves an event recipe
func (repo *Client) GetEventRecipe(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.EventRecipeId, fields []string) (cmodel.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("eventRecipeId", id.EventRecipeId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.EventRecipe{}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.EventRecipeFieldMasker.Convert(fields)).
		Where("event_recipe.event_recipe_id = ?", id.EventRecipeId)

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get event recipe row")
		return cmodel.EventRecipe{}, ConvertGormError(err)
	}

	m, err := convert.EventRecipeToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid event recipe row when getting event recipe")
		return cmodel.EventRecipe{}, repository.ErrInternal{Msg: "invalid event recipe row when getting event recipe"}
	}

	return m, nil
}

// ListEventRecipes lists event recipes
func (repo *Client) ListEventRecipes(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.EventRecipeParent, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("calendarId", parent.CalendarId).
		Int64("eventId", parent.EventId).
		Int32("pageSize", pageSize).
		Int64("offset", offset).
		Str("filter", filter).
		Strs("fields", fields).
		Logger()

	var gms []gmodel.EventRecipe

	tx := repo.db.WithContext(ctx).
		Select(gmodel.EventRecipeFieldMasker.Convert(fields)).
		Where("event_recipe.event_id = ?", parent.EventId)

	if pageSize > 0 {
		tx = tx.Limit(int(pageSize))
	}
	if offset > 0 {
		tx = tx.Offset(int(offset))
	}

	err := tx.Find(&gms).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list event recipe rows")
		return nil, repository.ErrInternal{Msg: "unable to list event recipe rows"}
	}

	ms := make([]cmodel.EventRecipe, len(gms))
	for i, gm := range gms {
		m, err := convert.EventRecipeToCoreModel(gm)
		if err != nil {
			log.Error().Err(err).Msg("invalid event recipe row when listing event recipes")
			return nil, repository.ErrInternal{Msg: "invalid event recipe row when listing event recipes"}
		}
		ms[i] = m
	}

	return ms, nil
}
