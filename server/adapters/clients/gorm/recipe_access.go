package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RecipeAccessMap maps the core model fields to the database model fields for the unified RecipeAccess model.
var RecipeAccessMap = map[string]string{
	model.RecipeAccessFields.Level:           dbModel.RecipeAccessFields.PermissionLevel,
	model.RecipeAccessFields.State:           dbModel.RecipeAccessFields.State,
	model.RecipeAccessFields.RecipientUser:   dbModel.RecipeAccessFields.RecipientUserId,
	model.RecipeAccessFields.RecipientCircle: dbModel.RecipeAccessFields.RecipientCircleId,
}

func (repo *Client) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	db := repo.db.WithContext(ctx)

	// Validate that exactly one recipient type is set
	if (access.RecipeAccessParent.Recipient.UserId != 0) == (access.RecipeAccessParent.Recipient.CircleId != 0) {
		return model.RecipeAccess{}, repository.ErrInvalidArgument{Msg: "exactly one recipient (user or circle) is required"}
	}

	recipeAccess := convert.CoreRecipeAccessToRecipeAccess(access)
	res := db.Create(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return model.RecipeAccess{}, repository.ErrNewAlreadyExists{}
		}
		return model.RecipeAccess{}, res.Error
	}

	access.RecipeAccessId.RecipeAccessId = recipeAccess.RecipeAccessId
	return access, nil
}

func (repo *Client) DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	if parent.RecipeId.RecipeId == 0 {
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	if id.RecipeAccessId == 0 {
		return repository.ErrInvalidArgument{Msg: "recipe access id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ? AND recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent) error {
	if parent.RecipeId.RecipeId == 0 {
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	db := repo.db.WithContext(ctx)

	res := db.Where("recipe_id = ?", parent.RecipeId.RecipeId).Delete(&dbModel.RecipeAccess{})
	if res.Error != nil {
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	db := repo.db.WithContext(ctx)

	var recipeAccess dbModel.RecipeAccess
	res := db.Where("recipe_id = ? AND recipe_access_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).First(&recipeAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return model.RecipeAccess{}, repository.ErrNotFound{}
		}
		return model.RecipeAccess{}, res.Error
	}

	return convert.RecipeAccessToCoreRecipeAccess(recipeAccess), nil
}

func (repo *Client) ListRecipeAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.RecipeAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.RecipeAccess, error) {
	if authAccount.UserId == 0 {
		return nil, repository.ErrInvalidArgument{Msg: "user id is required"}
	}

	conversion, err := repo.recipeAccessSQLConverter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	var recipeAccesses []dbModel.RecipeAccess
	db := repo.db.WithContext(ctx).Model(&dbModel.RecipeAccess{})

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	// Filter by recipe ID
	if parent.RecipeId.RecipeId != 0 {
		db = db.Where("recipe_access.recipe_id = ?", parent.RecipeId.RecipeId)
	}

	if authAccount.CircleId != 0 {
		db = db.Where("recipe_access.recipient_circle_id = ?", authAccount.CircleId)
	} else if authAccount.UserId != 0 {
		db = db.Where("recipe_access.recipient_user_id = ?", authAccount.UserId)
	}

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&recipeAccesses).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.RecipeAccess, len(recipeAccesses))
	for i, access := range recipeAccesses {
		accesses[i] = convert.RecipeAccessToCoreRecipeAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	dbAccess := convert.CoreRecipeAccessToRecipeAccess(access)

	db := repo.db.WithContext(ctx).Select("state").Clauses(&clause.Returning{})

	err := db.Where("recipe_access_id = ?", access.RecipeAccessId.RecipeAccessId).Updates(&dbAccess).Error
	if err != nil {
		return model.RecipeAccess{}, ConvertGormError(err)
	}

	return convert.RecipeAccessToCoreRecipeAccess(dbAccess), nil
}
