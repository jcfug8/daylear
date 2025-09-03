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

// CreateListItem creates a new list item
func (repo *Client) CreateListItem(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.ListItem) (cmodel.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", m.Parent.ListId.ListId).
		Str("title", m.Title).
		Logger()

	gm, err := convert.ListItemFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item when creating list item row")
		return cmodel.ListItem{}, repository.ErrInvalidArgument{Msg: "invalid list item"}
	}

	err = repo.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create list item row")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	m, err = convert.ListItemToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item row when creating list item")
		return cmodel.ListItem{}, repository.ErrInternal{Msg: "invalid list item row when creating list item"}
	}

	return m, nil
}

// DeleteListItem deletes a list item
func (repo *Client) DeleteListItem(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListItemId) (cmodel.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listItemId", id.ListItemId).
		Logger()

	gm := gmodel.ListItem{ListItemId: id.ListItemId}

	// First get the item to return it
	err := repo.db.WithContext(ctx).First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get list item for deletion")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	// Convert to core model for return
	m, err := convert.ListItemToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item row when deleting list item")
		return cmodel.ListItem{}, repository.ErrInternal{Msg: "invalid list item row when deleting list item"}
	}

	// Delete the item
	err = repo.db.WithContext(ctx).Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list item row")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	return m, nil
}

// GetListItem retrieves a list item
func (repo *Client) GetListItem(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListItemId, fields []string) (cmodel.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listItemId", id.ListItemId).
		Strs("fields", fields).
		Logger()

	var gm gmodel.ListItem

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		tx = tx.Select(gmodel.ListItemFieldMasker.Convert(fields))
	}

	err := tx.First(&gm, id.ListItemId).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get list item row")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	m, err := convert.ListItemToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item row when getting list item")
		return cmodel.ListItem{}, repository.ErrInternal{Msg: "invalid list item row when getting list item"}
	}

	return m, nil
}

// ListListItems lists list items with pagination and filtering
func (repo *Client) ListListItems(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.ListItemParent, pageSize int32, pageOffset int32, filter string, fields []string) ([]cmodel.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Int32("pageSize", pageSize).
		Int32("pageOffset", pageOffset).
		Str("filter", filter).
		Strs("fields", fields).
		Logger()

	var gms []gmodel.ListItem

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "list_item.list_item_id"},
		Desc:   false,
	}}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.ListItemFieldMasker.Convert(fields)).
		Where("list_item.list_id = ?", parent.ListId.ListId).
		Order(clause.OrderBy{Columns: orders})

	if pageSize > 0 {
		tx = tx.Limit(int(pageSize))
	}
	if pageOffset > 0 {
		tx = tx.Offset(int(pageOffset))
	}

	// Apply filter if provided
	if filter != "" {
		conversion, err := gmodel.ListItemSQLConverter.Convert(filter)
		if err != nil {
			log.Error().Err(err).Msg("invalid filter string when listing list item rows")
			return []cmodel.ListItem{}, repository.ErrInvalidArgument{Msg: "invalid filter"}
		}

		if conversion.WhereClause != "" {
			tx = tx.Where(conversion.WhereClause, conversion.Params...)
		}
	}

	err := tx.Find(&gms).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list list item rows")
		return []cmodel.ListItem{}, ConvertGormError(err)
	}

	ms := make([]cmodel.ListItem, len(gms))
	for i, gm := range gms {
		m, err := convert.ListItemToCoreModel(gm)
		if err != nil {
			log.Error().Err(err).Msg("invalid list item row when listing list items")
			return []cmodel.ListItem{}, repository.ErrInternal{Msg: "invalid list item row when listing list items"}
		}
		ms[i] = m
	}

	return ms, nil
}

// UpdateListItem updates an existing list item
func (repo *Client) UpdateListItem(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.ListItem, fields []string) (cmodel.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listItemId", m.Id.ListItemId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.ListItemFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item when updating list item row")
		return cmodel.ListItem{}, repository.ErrInvalidArgument{Msg: "invalid list item"}
	}

	tx := repo.db.WithContext(ctx).
		Clauses(clause.Returning{})

	if len(fields) > 0 {
		tx = tx.Select(gmodel.ListItemFieldMasker.Convert(fields))
	}

	err = tx.Where("list_item_id = ?", m.Id.ListItemId).Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update list item row")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	// Get the updated item
	err = repo.db.WithContext(ctx).First(&gm, m.Id.ListItemId).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get updated list item row")
		return cmodel.ListItem{}, ConvertGormError(err)
	}

	m, err = convert.ListItemToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list item row when updating list item")
		return cmodel.ListItem{}, repository.ErrInternal{Msg: "invalid list item row when updating list item"}
	}

	return m, nil
}

// BulkDeleteListItems deletes all list items for a given parent list
func (repo *Client) BulkDeleteListItems(ctx context.Context, parent cmodel.ListItemParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Logger()

	err := repo.db.WithContext(ctx).
		Where("list_id = ?", parent.ListId.ListId).
		Delete(&gmodel.ListItem{}).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to bulk delete list item rows")
		return ConvertGormError(err)
	}

	return nil
}
