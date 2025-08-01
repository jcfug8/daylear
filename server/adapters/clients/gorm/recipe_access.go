package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UpdateRecipeAccessMap maps the updatable core model fields to the database model fields for the RecipeAccess model.
var UpdateRecipeAccessMap = map[string][]string{
	model.RecipeAccessFields.Level: []string{
		dbModel.RecipeAccessFields.PermissionLevel,
	},
	model.RecipeAccessFields.State: []string{
		dbModel.RecipeAccessFields.State,
	},
}

var UpdateRecipeAccessFieldMasker = fieldmask.NewFieldMasker(UpdateRecipeAccessMap)

// RecipeAccessMap maps the core model fields to the database model fields for the unified RecipeAccess model.
var RecipeAccessMap = map[string]string{
	model.RecipeAccessFields.Level:           dbModel.RecipeAccessFields.PermissionLevel,
	model.RecipeAccessFields.State:           dbModel.RecipeAccessFields.State,
	model.RecipeAccessFields.RecipientUser:   dbModel.RecipeAccessFields.RecipientUserId,
	model.RecipeAccessFields.RecipientCircle: dbModel.RecipeAccessFields.RecipientCircleId,
}

func (repo *Client) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	db := repo.db.WithContext(ctx)

	// Validate that exactly one recipient type is set
	if (access.Recipient.UserId != 0) == (access.Recipient.CircleId != 0) {
		log.Error().Msg("exactly one recipient (user or circle) is required")
		return model.RecipeAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	recipeAccess := convert.CoreRecipeAccessToRecipeAccess(access)
	res := db.Clauses(clause.Returning{}).Create(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("duplicate key error")
			return model.RecipeAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("db.Create failed")
		return model.RecipeAccess{}, res.Error
	}

	access.RecipeAccessId.RecipeAccessId = recipeAccess.RecipeAccessId
	return access, nil
}

func (repo *Client) DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	if parent.RecipeId.RecipeId == 0 {
		log.Error().Msg("recipe id is required")
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	if id.RecipeAccessId == 0 {
		log.Error().Msg("recipe access id is required")
		return repository.ErrInvalidArgument{Msg: "recipe access id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ? AND recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("db.Delete failed")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no rows affected (not found)")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	if parent.RecipeId.RecipeId == 0 {
		log.Error().Msg("recipe id is required")
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ?", parent.RecipeId.RecipeId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("db.Delete failed")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no rows affected (not found)")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	db := repo.db.WithContext(ctx)

	var recipeAccess dbModel.RecipeAccess
	res := db.Table("recipe_access").
		Select(`recipe_access.*, u.username as recipient_username, u.given_name as recipient_given_name, u.family_name as recipient_family_name, c.title as recipient_circle_title, c.handle as recipient_circle_handle`).
		Joins(`LEFT JOIN daylear_user u ON recipe_access.recipient_user_id = u.user_id`).
		Joins(`LEFT JOIN circle c ON recipe_access.recipient_circle_id = c.circle_id`).
		Where("recipe_access.recipe_id = ? AND recipe_access.recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).
		First(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("record not found")
			return model.RecipeAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("db.First failed")
		return model.RecipeAccess{}, res.Error
	}

	return convert.RecipeAccessToCoreRecipeAccess(recipeAccess), nil
}

func (repo *Client) ListRecipeAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		return nil, repository.ErrInvalidArgument{Msg: "user id or circle id is required"}
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "recipe_access.recipe_access_id"},
		Desc:   true,
	}}

	var recipeAccesses []dbModel.RecipeAccess
	// Start building the query
	db := repo.db.WithContext(ctx).
		Table("recipe_access").
		Select(`recipe_access.*, u.username as recipient_username, u.given_name as recipient_given_name, u.family_name as recipient_family_name, c.title as recipient_circle_title, c.handle as recipient_circle_handle`).
		Joins(`LEFT JOIN daylear_user u ON recipe_access.recipient_user_id = u.user_id`).
		Joins(`LEFT JOIN circle c ON recipe_access.recipient_circle_id = c.circle_id`).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	// Filter by recipe ID if provided
	if parent.RecipeId.RecipeId != 0 {
		db = db.Where("recipe_access.recipe_id = ?", parent.RecipeId.RecipeId)
	}

	if authAccount.CircleId != 0 {
		db = db.Where(
			"recipe_access.recipient_circle_id = ? OR recipe_access.recipe_id IN (SELECT recipe_id FROM recipe_access WHERE recipient_circle_id = ? AND permission_level >= ?)",
			authAccount.CircleId, authAccount.CircleId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	} else if authAccount.AuthUserId != 0 {
		db = db.Where(
			"recipe_access.recipient_user_id = ? OR recipe_access.recipe_id IN (SELECT recipe_id FROM recipe_access WHERE recipient_user_id = ? AND permission_level >= ?)",
			authAccount.AuthUserId, authAccount.AuthUserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
		)
	}

	conversion, err := repo.recipeAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&recipeAccesses).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.RecipeAccess, len(recipeAccesses))
	for i, access := range recipeAccesses {
		accesses[i] = convert.RecipeAccessToCoreRecipeAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess, updateMask []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	dbAccess := convert.CoreRecipeAccessToRecipeAccess(access)

	columns := UpdateRecipeAccessFieldMasker.Convert(updateMask)

	db := repo.db.WithContext(ctx).Select(columns).Clauses(&clause.Returning{})

	err := db.Where("recipe_access_id = ?", access.RecipeAccessId.RecipeAccessId).Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return model.RecipeAccess{}, ConvertGormError(err)
	}

	return convert.RecipeAccessToCoreRecipeAccess(dbAccess), nil
}
