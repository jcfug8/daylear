package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// ListLists lists lists.
func (repo *Client) ListLists(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, pageOffset int32, filter string, fields []string) ([]cmodel.List, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("filter", filter).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int("pageOffset", int(pageOffset)).
		Logger()

	if authAccount.PermissionLevel < types.PermissionLevel_PERMISSION_LEVEL_PUBLIC {
		return nil, repository.ErrInvalidArgument{Msg: "invalid permission level"}
	}

	dbLists := []gmodel.List{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "list.list_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.ListFieldMasker.Convert(fields)).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	if authAccount.CircleId != 0 {
		maxVisibilityLevel := types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC
		if authAccount.PermissionLevel > types.PermissionLevel_PERMISSION_LEVEL_PUBLIC {
			maxVisibilityLevel = types.VisibilityLevel_VISIBILITY_LEVEL_PRIVATE
		}
		tx = tx.Joins("LEFT JOIN list_access ON list.list_id = list_access.list_id AND list_access.recipient_circle_id = ? AND (list.visibility_level != ? OR list_access.permission_level = ?)", authAccount.CircleId, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN list_access as la ON list.list_id = la.list_id AND la.recipient_user_id = ? AND (list.visibility_level != ? OR la.permission_level = ?)", authAccount.AuthUserId, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN list_favorite ON list.list_id = list_favorite.list_id AND list_favorite.circle_id = ?", authAccount.CircleId).
			Where("(list.visibility_level <= ? OR (list_access.recipient_circle_id = ? AND la.state = ?))",
				maxVisibilityLevel, authAccount.CircleId, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else if authAccount.UserId != 0 && authAccount.UserId != authAccount.AuthUserId {
		maxVisibilityLevel := types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC
		if authAccount.PermissionLevel > types.PermissionLevel_PERMISSION_LEVEL_PUBLIC {
			maxVisibilityLevel = types.VisibilityLevel_VISIBILITY_LEVEL_RESTRICTED
		}
		tx = tx.Joins("LEFT JOIN list_access ON list.list_id = list_access.list_id AND list_access.recipient_user_id = ? AND (list.visibility_level != ? OR list_access.permission_level = ?)", authAccount.UserId, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN list_access as la ON list.list_id = la.list_id AND la.recipient_user_id = ? AND (list.visibility_level != ? OR la.permission_level = ?)", authAccount.AuthUserId, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN list_favorite ON list.list_id = list_favorite.list_id AND list_favorite.user_id = ?", authAccount.UserId).
			Where("(list.visibility_level <= ? OR (list_access.recipient_user_id = ? AND la.state = ?))",
				maxVisibilityLevel, authAccount.UserId, types.AccessState_ACCESS_STATE_ACCEPTED)
	} else {
		tx = tx.Joins("LEFT JOIN list_access ON list.list_id = list_access.list_id AND list_access.recipient_user_id = ? AND (list.visibility_level != ? OR list_access.permission_level = ?)", authAccount.AuthUserId, types.VisibilityLevel_VISIBILITY_LEVEL_HIDDEN, types.PermissionLevel_PERMISSION_LEVEL_ADMIN).
			Joins("LEFT JOIN list_favorite ON list.list_id = list_favorite.list_id AND list_favorite.user_id = ?", authAccount.AuthUserId).
			Where("(list.visibility_level = ? OR list_access.recipient_user_id = ?)", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.AuthUserId)
	}

	conversion, err := gmodel.ListSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing list rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbLists).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list list rows")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.List, len(dbLists))
	for i, m := range dbLists {
		res[i], err = convert.ListToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("invalid list row when listing lists")
			return nil, repository.ErrInternal{Msg: "invalid list row when listing lists"}
		}
	}

	return res, nil
}

// CreateList creates a new list.
func (repo *Client) CreateList(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.List) (cmodel.List, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Logger()

	gm, err := convert.ListFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid list when creating list row")
		return cmodel.List{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid list: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.ListFieldMasker.Convert([]string{})).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create list row")
		return cmodel.List{}, ConvertGormError(err)
	}

	m, err = convert.ListToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list row when creating list")
		return cmodel.List{}, fmt.Errorf("unable to read list: %v", err)
	}

	return m, nil
}

// DeleteList deletes a list.
func (repo *Client) DeleteList(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) (cmodel.List, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Logger()

	gm := gmodel.List{ListId: id.ListId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.ListFieldMasker.Get()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list row")
		return cmodel.List{}, ConvertGormError(err)
	}

	m, err := convert.ListToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list row when deleting list")
		return cmodel.List{}, fmt.Errorf("unable to read list: %v", err)
	}

	return m, nil
}

// GetList gets a list.
func (repo *Client) GetList(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId, fields []string) (cmodel.List, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.List{}

	tx := repo.db.WithContext(ctx).
		Select(gmodel.ListFieldMasker.Convert(
			fields,
			fieldmask.ExcludeKeys(
				cmodel.ListField_ListAccess,
			),
		)).
		Joins("LEFT JOIN list_favorite ON list.list_id = list_favorite.list_id AND list_favorite.user_id = ?", authAccount.AuthUserId).
		Where("list.list_id = ?", id.ListId)

	err := tx.First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get list row")
		return cmodel.List{}, ConvertGormError(err)
	}

	m, err := convert.ListToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list row when getting list")
		return cmodel.List{}, fmt.Errorf("unable to read list: %v", err)
	}

	return m, nil
}

// UpdateList updates a list.
func (repo *Client) UpdateList(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.List, fields []string) (cmodel.List, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", m.Id.ListId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.ListFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid list when updating list row")
		return cmodel.List{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("error reading list: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.ListFieldMasker.Convert(fields)).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update list row")
		return cmodel.List{}, ConvertGormError(err)
	}

	m, err = convert.ListToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid list row when updating list")
		return cmodel.List{}, fmt.Errorf("unable to read list: %v", err)
	}

	return m, nil
}

// CreateListFavorite creates a list favorite for a user.
func (repo *Client) CreateListFavorite(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Interface("authAccount", authAccount).
		Int64("listId", id.ListId).
		Logger()

	favorite := gmodel.ListFavorite{
		ListId: id.ListId,
	}

	if authAccount.CircleId != 0 {
		favorite.CircleId = authAccount.CircleId
	} else if authAccount.UserId != 0 {
		favorite.UserId = authAccount.UserId
	} else {
		favorite.UserId = authAccount.AuthUserId
	}

	err := repo.db.WithContext(ctx).Create(&favorite).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create list favorite")
		return ConvertGormError(err)
	}

	log.Info().Msg("list favorite created successfully")
	return nil
}

// DeleteListFavorite removes a list favorite for a user.
func (repo *Client) DeleteListFavorite(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Interface("authAccount", authAccount).
		Int64("listId", id.ListId).
		Logger()

	tx := repo.db.WithContext(ctx)

	if authAccount.CircleId != 0 {
		tx = tx.Where("circle_id = ? AND list_id = ?", authAccount.CircleId, id.ListId)
	} else if authAccount.UserId != 0 {
		tx = tx.Where("user_id = ? AND list_id = ?", authAccount.UserId, id.ListId)
	} else {
		tx = tx.Where("user_id = ? AND list_id = ?", authAccount.AuthUserId, id.ListId)
	}

	err := tx.Delete(&gmodel.ListFavorite{}).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete list favorite")
		return ConvertGormError(err)
	}

	log.Info().Msg("list favorite deleted successfully")
	return nil
}

// BulkDeleteListFavorites removes all list favorites for a list.
func (repo *Client) BulkDeleteListFavorites(ctx context.Context, id cmodel.ListId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Logger()

	err := repo.db.WithContext(ctx).
		Where("list_id = ?", id.ListId).
		Delete(&gmodel.ListFavorite{}).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to bulk delete list favorite")
		return ConvertGormError(err)
	}

	log.Info().Msg("list favorites deleted successfully")
	return nil
}
