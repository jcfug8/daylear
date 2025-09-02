package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Client) CreateListAccess(ctx context.Context, access cmodel.ListAccess, fields []string) (cmodel.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", access.ListId.ListId).
		Int64("recipientUserId", access.Recipient.UserId).
		Int64("recipientCircleId", access.Recipient.CircleId).
		Logger()

	// Validate that exactly one recipient type is set
	if (access.Recipient.UserId != 0) == (access.Recipient.CircleId != 0) {
		log.Error().Msg("exactly one recipient (user or circle) is required to create list access row")
		return cmodel.ListAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	listAccess := convert.ListAccessFromCoreModel(access)
	res := repo.db.WithContext(ctx).
		Select(dbModel.ListAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&listAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("unable to create list access row: already exists")
			return cmodel.ListAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("unable to create list access row")
		return cmodel.ListAccess{}, ConvertGormError(res.Error)
	}

	access.ListAccessId.ListAccessId = listAccess.ListAccessId
	return access, nil
}

func (repo *Client) DeleteListAccess(ctx context.Context, parent cmodel.ListAccessParent, id cmodel.ListAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Int64("listAccessId", id.ListAccessId).
		Logger()

	if parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required to delete list access row")
		return repository.ErrInvalidArgument{Msg: "list id is required"}
	}

	if id.ListAccessId == 0 {
		log.Error().Msg("list access id is required to delete list access row")
		return repository.ErrInvalidArgument{Msg: "list access id is required"}
	}

	res := repo.db.WithContext(ctx).
		Where("list_id = ? AND list_access_id = ?", parent.ListId.ListId, id.ListAccessId).
		Delete(&dbModel.ListAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to delete list access row")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no list access row deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteListAccess(ctx context.Context, parent cmodel.ListAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Logger()

	if parent.ListId.ListId == 0 {
		log.Error().Msg("list id is required to bulk delete list access rows")
		return repository.ErrInvalidArgument{Msg: "list id is required"}
	}

	res := repo.db.WithContext(ctx).
		Where("list_id = ?", parent.ListId.ListId).
		Delete(&dbModel.ListAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to bulk delete list access rows")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no list access rows deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetListAccess(ctx context.Context, parent cmodel.ListAccessParent, id cmodel.ListAccessId, fields []string) (cmodel.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Int64("listAccessId", id.ListAccessId).
		Logger()

	dbFields := dbModel.ListAccessFieldMasker.Convert(fields)

	var listAccess dbModel.ListAccess
	tx := repo.db.WithContext(ctx).
		Select(dbFields).
		Where("list_access.list_id = ? AND list_access.list_access_id = ?", parent.ListId.ListId, id.ListAccessId)

	if fieldmask.ContainsAny(dbFields, dbModel.ListAccessFields_RecipientUserId, dbModel.ListAccessFields_RecipientCircleId) {
		tx = tx.Joins(`LEFT JOIN daylear_user ON list_access.recipient_user_id = daylear_user.user_id`).
			Joins(`LEFT JOIN circle ON list_access.recipient_circle_id = circle.circle_id`)
	}

	res := tx.First(&listAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("list access row not found")
			return cmodel.ListAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to get list access row")
		return cmodel.ListAccess{}, ConvertGormError(res.Error)
	}

	return convert.ListAccessToCoreModel(listAccess), nil
}

func (repo *Client) ListListAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent cmodel.ListAccessParent, pageSize int32, pageOffset int64, filterStr string, fields []string) ([]cmodel.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", parent.ListId.ListId).
		Str("filter", filterStr).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).
		Logger()

	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Error().Msg("user id or circle id is required to list list access rows")
		return nil, repository.ErrInvalidArgument{Msg: "user id or circle id is required"}
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "list_access.list_access_id"},
		Desc:   true,
	}}

	dbFields := dbModel.ListAccessFieldMasker.Convert(fields)

	var listAccesses []dbModel.ListAccess
	// Start building the query
	tx := repo.db.WithContext(ctx).
		Select(dbFields).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	if fieldmask.ContainsAny(dbFields, dbModel.ListAccessFields_RecipientUserId, dbModel.ListAccessFields_RecipientCircleId) {
		tx = tx.Joins(`LEFT JOIN daylear_user ON list_access.recipient_user_id = daylear_user.user_id`).
			Joins(`LEFT JOIN circle ON list_access.recipient_circle_id = circle.circle_id`)
	}

	// Filter by list ID if provided
	if parent.ListId.ListId != 0 {
		tx = tx.Where("list_access.list_id = ?", parent.ListId.ListId)
	}

	if authAccount.CircleId != 0 {
		tx = tx.Where(
			"list_access.recipient_circle_id = ? OR list_access.list_id IN (SELECT list_id FROM list_access WHERE recipient_circle_id = ? AND permission_level >= ?)",
			authAccount.CircleId, authAccount.CircleId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	} else if authAccount.AuthUserId != 0 {
		tx = tx.Where(
			"list_access.recipient_user_id = ? OR list_access.list_id IN (SELECT list_id FROM list_access WHERE recipient_user_id = ? AND permission_level >= ?)",
			authAccount.AuthUserId, authAccount.AuthUserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	}

	conversion, err := dbModel.ListAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing list access rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	res := tx.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&listAccesses)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to list list access rows")
		return nil, ConvertGormError(res.Error)
	}

	accesses := make([]cmodel.ListAccess, len(listAccesses))
	for i, access := range listAccesses {
		accesses[i] = convert.ListAccessToCoreModel(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateListAccess(ctx context.Context, access cmodel.ListAccess, fields []string) (cmodel.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listAccessId", access.ListAccessId.ListAccessId).
		Strs("fields", fields).
		Logger()

	dbAccess := convert.ListAccessFromCoreModel(access)

	res := repo.db.WithContext(ctx).
		Select(dbModel.ListAccessFieldMasker.Convert(fields, fieldmask.OnlyUpdatable())).
		Clauses(&clause.Returning{}).
		Where("list_access_id = ?", access.ListAccessId.ListAccessId).
		Updates(&dbAccess)
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to update list access row")
		return cmodel.ListAccess{}, ConvertGormError(res.Error)
	}

	return convert.ListAccessToCoreModel(dbAccess), nil
}

func (repo *Client) FindStandardUserListAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) (cmodel.ListAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	var listAccess dbModel.ListAccess
	res := repo.db.WithContext(ctx).
		Where("list_id = ? AND recipient_user_id = ?", id.ListId, authAccount.AuthUserId).
		First(&listAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("standard user list access not found")
			return cmodel.ListAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find standard user list access")
		return cmodel.ListAccess{}, ConvertGormError(res.Error)
	}
	return convert.ListAccessToCoreModel(listAccess), nil
}

func (repo *Client) FindDelegatedCircleListAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) (cmodel.ListAccess, cmodel.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.ListAccess
		dbModel.CircleAccess
	}
	var result Result
	res := repo.db.WithContext(ctx).
		Select("list_access.*, circle_access.*").
		Table("list_access").
		Joins("JOIN circle_access ON circle_access.circle_id = list_access.recipient_circle_id").
		Where("list_access.list_id = ? AND circle_access.recipient_user_id = ? AND list_access.state = ?", id.ListId, authAccount.AuthUserId, types.AccessState_ACCESS_STATE_ACCEPTED).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated circle list access not found")
			return cmodel.ListAccess{}, cmodel.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated circle list access")
		return cmodel.ListAccess{}, cmodel.CircleAccess{}, ConvertGormError(res.Error)
	}
	return convert.ListAccessToCoreModel(result.ListAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserListAccess(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.ListId) (cmodel.ListAccess, cmodel.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("listId", id.ListId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.ListAccess
		dbModel.UserAccess
	}
	var result Result

	res := repo.db.WithContext(ctx).
		Select("list_access.*, user_access.*").
		Table("list_access").
		Joins("JOIN user_access ON user_access.user_id = list_access.recipient_user_id").
		Where("list_access.list_id = ? AND user_access.recipient_user_id = ? AND list_access.state = ?", id.ListId, authAccount.AuthUserId, types.AccessState_ACCESS_STATE_ACCEPTED).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated user list access not found")
			return cmodel.ListAccess{}, cmodel.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated user list access")
		return cmodel.ListAccess{}, cmodel.UserAccess{}, ConvertGormError(res.Error)
	}
	return convert.ListAccessToCoreModel(result.ListAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
