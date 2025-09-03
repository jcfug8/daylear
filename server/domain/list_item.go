package domain

import (
	"context"
	"slices"

	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

// CreateListItem creates a new list item
func (d *Domain) CreateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem) (dbListItem model.ListItem, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when creating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if listItem.Parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required when creating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	if listItem.Title == "" {
		log.Error().Msg("title is required when creating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "title is required"}
	}

	// Check access to the parent list
	_, err = d.determineListAccess(ctx, authAccount, listItem.Parent.ListId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when creating list item")
		return model.ListItem{}, err
	}

	// Set ID to 0 for new items
	listItem.Id.ListItemId = 0

	dbListItem, err = d.repo.CreateListItem(ctx, authAccount, listItem)
	if err != nil {
		log.Error().Err(err).Msg("unable to create list item")
		return model.ListItem{}, domain.ErrInternal{Msg: "unable to create list item"}
	}

	return dbListItem, nil
}

// GetListItem retrieves a list item
func (d *Domain) GetListItem(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, id model.ListItemId, fields []string) (dbListItem model.ListItem, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when getting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required when getting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	if id.ListItemId == 0 {
		log.Error().Msg("list item id is required when getting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list item id is required"}
	}

	dbList, err := d.repo.GetList(ctx, authAccount, parent.ListId, []string{model.ListField_VisibilityLevel})
	if err != nil {
		log.Error().Err(err).Msg("unable to get list")
		return model.ListItem{}, domain.ErrInternal{Msg: "unable to get list"}
	}

	// Check access to the parent list
	_, err = d.determineListAccess(
		ctx, authAccount, parent.ListId,
		withResourceVisibilityLevel(dbList.VisibilityLevel),
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when getting list item")
		return model.ListItem{}, err
	}

	dbListItem, err = d.repo.GetListItem(ctx, authAccount, id, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to get list item")
		return model.ListItem{}, domain.ErrInternal{Msg: "unable to get list item"}
	}

	return dbListItem, nil
}

// ListListItems lists list items with pagination and filtering
func (d *Domain) ListListItems(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, pageSize int32, pageOffset int32, filter string, fields []string) (dbListItems []model.ListItem, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when listing list items")
		return []model.ListItem{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required when listing list items")
		return []model.ListItem{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	dbList, err := d.repo.GetList(ctx, authAccount, parent.ListId, []string{model.ListField_VisibilityLevel})
	if err != nil {
		log.Error().Err(err).Msg("unable to get list")
		return []model.ListItem{}, domain.ErrInternal{Msg: "unable to get list"}
	}

	// Check access to the parent list
	_, err = d.determineListAccess(
		ctx, authAccount, parent.ListId,
		withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_READ),
		withResourceVisibilityLevel(dbList.VisibilityLevel),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when listing list items")
		return []model.ListItem{}, err
	}

	dbListItems, err = d.repo.ListListItems(ctx, authAccount, parent, pageSize, pageOffset, filter, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to list list items")
		return []model.ListItem{}, domain.ErrInternal{Msg: "unable to list list items"}
	}

	return dbListItems, nil
}

// UpdateListItem updates an existing list item
func (d *Domain) UpdateListItem(ctx context.Context, authAccount model.AuthAccount, listItem model.ListItem, fields []string) (dbListItem model.ListItem, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when updating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if listItem.Parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required when updating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	if listItem.Id.ListItemId == 0 {
		log.Error().Msg("list item id is required when updating list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list item id is required"}
	}

	// Check access to the parent list
	_, err = d.determineListAccess(ctx, authAccount, listItem.Parent.ListId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when updating list item")
		return model.ListItem{}, err
	}

	// Validate title if it's being updated
	if len(fields) == 0 || slices.Contains(fields, model.ListItemField_Title) {
		if listItem.Title == "" {
			log.Error().Msg("title is required when updating list item")
			return model.ListItem{}, domain.ErrInvalidArgument{Msg: "title is required"}
		}
	}

	dbListItem, err = d.repo.UpdateListItem(ctx, authAccount, listItem, fields)
	if err != nil {
		log.Error().Err(err).Msg("unable to update list item")
		return model.ListItem{}, domain.ErrInternal{Msg: "unable to update list item"}
	}

	return dbListItem, nil
}

// DeleteListItem deletes a list item
func (d *Domain) DeleteListItem(ctx context.Context, authAccount model.AuthAccount, parent model.ListItemParent, id model.ListItemId) (dbListItem model.ListItem, err error) {
	log := logutil.EnrichLoggerWithContext(d.log, ctx)

	if authAccount.AuthUserId == 0 {
		log.Error().Msg("user id is required when deleting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "user id is required"}
	}

	if parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required when deleting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list id is required"}
	}

	if id.ListItemId == 0 {
		log.Error().Msg("list item id is required when deleting list item")
		return model.ListItem{}, domain.ErrInvalidArgument{Msg: "list item id is required"}
	}

	// Check access to the parent list
	_, err = d.determineListAccess(ctx, authAccount, parent.ListId, withMinimumPermissionLevel(types.PermissionLevel_PERMISSION_LEVEL_WRITE))
	if err != nil {
		log.Error().Err(err).Msg("unable to determine access when deleting list item")
		return model.ListItem{}, err
	}

	dbListItem, err = d.repo.DeleteListItem(ctx, authAccount, id)
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list item")
		return model.ListItem{}, domain.ErrInternal{Msg: "unable to delete list item"}
	}

	return dbListItem, nil
}
