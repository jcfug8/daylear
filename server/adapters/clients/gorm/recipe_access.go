package gorm

import (
	"context"
	"errors"
	"slices"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (repo *Client) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess, fields []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", access.RecipeId.RecipeId).
		Int64("recipientUserId", access.Recipient.UserId).
		Int64("recipientCircleId", access.Recipient.CircleId).
		Logger()

	// Validate that exactly one recipient type is set
	if (access.Recipient.UserId != 0) == (access.Recipient.CircleId != 0) {
		log.Error().Msg("exactly one recipient (user or circle) is required to create recipe access row")
		return model.RecipeAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	recipeAccess := convert.CoreRecipeAccessToRecipeAccess(access)
	res := repo.db.WithContext(ctx).
		Select(dbModel.RecipeAccessFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			log.Error().Err(res.Error).Msg("unable to create recipe access row: already exists")
			return model.RecipeAccess{}, repository.ErrNewAlreadyExists{}
		}
		log.Error().Err(res.Error).Msg("unable to create recipe access row")
		return model.RecipeAccess{}, res.Error
	}

	access.RecipeAccessId.RecipeAccessId = recipeAccess.RecipeAccessId
	return access, nil
}

func (repo *Client) DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", parent.RecipeId.RecipeId).
		Int64("recipeAccessId", id.RecipeAccessId).
		Logger()

	if parent.RecipeId.RecipeId == 0 {
		log.Error().Msg("recipe id is required to delete recipe access row")
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	if id.RecipeAccessId == 0 {
		log.Error().Msg("recipe access id is required to delete recipe access row")
		return repository.ErrInvalidArgument{Msg: "recipe access id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ? AND recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to delete recipe access row")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no recipe access row deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent) error {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", parent.RecipeId.RecipeId).
		Logger()

	if parent.RecipeId.RecipeId == 0 {
		log.Error().Msg("recipe id is required to bulk delete recipe access rows")
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ?", parent.RecipeId.RecipeId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		log.Error().Err(res.Error).Msg("unable to bulk delete recipe access rows")
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		log.Warn().Msg("no recipe access rows deleted")
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId, fields []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", parent.RecipeId.RecipeId).
		Int64("recipeAccessId", id.RecipeAccessId).
		Logger()

	var recipeAccess dbModel.RecipeAccess
	tx := repo.db.WithContext(ctx).
		Select(dbModel.RecipeAccessFieldMasker.Convert(fields)).
		Where("recipe_access.recipe_id = ? AND recipe_access.recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId)

	if slices.Contains(fields, dbModel.RecipeAccessFields_RecipientUserId) || slices.Contains(fields, dbModel.RecipeAccessFields_RecipientCircleId) {
		tx = tx.Joins(`LEFT JOIN daylear_user ON recipe_access.recipient_user_id = daylear_user.user_id`).
			Joins(`LEFT JOIN circle ON recipe_access.recipient_circle_id = circle.circle_id`)
	}

	res := tx.First(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("recipe access row not found")
			return model.RecipeAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to get recipe access row")
		return model.RecipeAccess{}, res.Error
	}

	return convert.RecipeAccessToCoreRecipeAccess(recipeAccess), nil
}

func (repo *Client) ListRecipeAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filterStr string, fields []string) ([]model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", parent.RecipeId.RecipeId).
		Str("filter", filterStr).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).
		Logger()

	if authAccount.AuthUserId == 0 && authAccount.CircleId == 0 {
		log.Error().Msg("user id or circle id is required to list recipe access rows")
		return nil, repository.ErrInvalidArgument{Msg: "user id or circle id is required"}
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "recipe_access.recipe_access_id"},
		Desc:   true,
	}}

	var recipeAccesses []dbModel.RecipeAccess
	// Start building the query
	db := repo.db.WithContext(ctx).
		Select(dbModel.RecipeAccessFieldMasker.Convert(fields)).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	if slices.Contains(fields, dbModel.RecipeAccessFields_RecipientUserId) || slices.Contains(fields, dbModel.RecipeAccessFields_RecipientCircleId) {
		db = db.Joins(`LEFT JOIN daylear_user u ON recipe_access.recipient_user_id = u.user_id`).
			Joins(`LEFT JOIN circle c ON recipe_access.recipient_circle_id = c.circle_id`)
	}

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

	conversion, err := dbModel.RecipeAccessSQLConverter.Convert(filterStr)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing recipe access rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&recipeAccesses).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list recipe access rows")
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.RecipeAccess, len(recipeAccesses))
	for i, access := range recipeAccesses {
		accesses[i] = convert.RecipeAccessToCoreRecipeAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess, fields []string) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeAccessId", access.RecipeAccessId.RecipeAccessId).
		Strs("fields", fields).
		Logger()

	dbAccess := convert.CoreRecipeAccessToRecipeAccess(access)

	db := repo.db.WithContext(ctx).
		Select(dbModel.UpdateRecipeAccessFieldMasker.Convert(fields)).
		Clauses(&clause.Returning{})

	err := db.Where("recipe_access_id = ?", access.RecipeAccessId.RecipeAccessId).Updates(&dbAccess).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update recipe access row")
		return model.RecipeAccess{}, ConvertGormError(err)
	}

	return convert.RecipeAccessToCoreRecipeAccess(dbAccess), nil
}

func (repo *Client) FindStandardUserRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", id.RecipeId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	var recipeAccess dbModel.RecipeAccess
	res := repo.db.WithContext(ctx).
		Where("recipe_id = ? AND recipient_user_id = ?", id.RecipeId, authAccount.AuthUserId).
		First(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("standard user recipe access not found")
			return model.RecipeAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find standard user recipe access")
		return model.RecipeAccess{}, res.Error
	}
	return convert.RecipeAccessToCoreRecipeAccess(recipeAccess), nil
}

func (repo *Client) FindDelegatedCircleRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, model.CircleAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", id.RecipeId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.RecipeAccess
		dbModel.CircleAccess
	}
	var result Result
	res := repo.db.WithContext(ctx).
		Select("recipe_access.*, circle_access.*").
		Table("recipe_access").
		Joins("JOIN circle_access ON circle_access.circle_id = recipe_access.recipient_circle_id").
		Where("recipe_access.recipe_id = ? AND circle_access.recipient_user_id = ?", id.RecipeId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated circle recipe access not found")
			return model.RecipeAccess{}, model.CircleAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated circle recipe access")
		return model.RecipeAccess{}, model.CircleAccess{}, res.Error
	}
	return convert.RecipeAccessToCoreRecipeAccess(result.RecipeAccess), convert.CircleAccessToCoreCircleAccess(result.CircleAccess), nil
}

func (repo *Client) FindDelegatedUserRecipeAccess(ctx context.Context, authAccount model.AuthAccount, id model.RecipeId) (model.RecipeAccess, model.UserAccess, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("recipeId", id.RecipeId).
		Int64("authUserId", authAccount.AuthUserId).
		Logger()

	type Result struct {
		dbModel.RecipeAccess
		dbModel.UserAccess
	}
	var result Result

	res := repo.db.WithContext(ctx).
		Select("recipe_access.*, user_access.*").
		Table("recipe_access").
		Joins("JOIN user_access ON user_access.user_id = recipe_access.recipient_user_id").
		Where("recipe_access.recipe_id = ? AND user_access.recipient_user_id = ?", id.RecipeId, authAccount.AuthUserId).
		First(&result)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Warn().Err(res.Error).Msg("delegated user recipe access not found")
			return model.RecipeAccess{}, model.UserAccess{}, repository.ErrNotFound{}
		}
		log.Error().Err(res.Error).Msg("unable to find delegated user recipe access")
		return model.RecipeAccess{}, model.UserAccess{}, res.Error
	}
	return convert.RecipeAccessToCoreRecipeAccess(result.RecipeAccess), convert.UserAccessToCoreUserAccess(result.UserAccess), nil
}
