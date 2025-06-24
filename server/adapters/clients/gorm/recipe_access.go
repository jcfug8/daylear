package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
)

// RecipeAccessUserMap maps the core model fields to the database model fields.
var RecipeAccessUserMap = map[string]string{
	model.RecipeAccessFields.Level:         dbModel.RecipeUserFields.PermissionLevel,
	model.RecipeAccessFields.State:         dbModel.RecipeUserFields.State,
	model.RecipeAccessFields.RecipientUser: dbModel.RecipeUserFields.UserId,
}

// RecipeAccessCircleMap maps the core model fields to the database model fields.
var RecipeAccessCircleMap = map[string]string{
	model.RecipeAccessFields.Level:           dbModel.RecipeCircleFields.PermissionLevel,
	model.RecipeAccessFields.State:           dbModel.RecipeCircleFields.State,
	model.RecipeAccessFields.RecipientCircle: dbModel.RecipeCircleFields.CircleId,
}

func (repo *Client) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	db := repo.db.WithContext(ctx)

	if access.RecipeAccessParent.Recipient.UserId != 0 {
		recipeUser := convert.CoreRecipeAccessToRecipeUser(access)
		res := db.Create(&recipeUser)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
				return model.RecipeAccess{}, repository.ErrNewAlreadyExists{}
			}
			return model.RecipeAccess{}, res.Error
		}
		access.RecipeAccessId.RecipeAccessId = recipeUser.RecipeUserId
	} else if access.RecipeAccessParent.Recipient.CircleId != 0 {
		recipeCircle := convert.CoreRecipeAccessToRecipeCircle(access)
		res := db.Create(&recipeCircle)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
				return model.RecipeAccess{}, repository.ErrNewAlreadyExists{}
			}
			return model.RecipeAccess{}, res.Error
		}
		access.RecipeAccessId.RecipeAccessId = recipeCircle.RecipeCircleId
	} else {
		return model.RecipeAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	return access, nil
}

func (repo *Client) DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	db := repo.db.WithContext(ctx)

	if parent.Recipient.UserId != 0 {
		res := db.Delete(&dbModel.RecipeUser{}, id.RecipeAccessId)
		if res.Error != nil {
			return ConvertGormError(res.Error)
		}
		if res.RowsAffected == 0 {
			return repository.ErrNotFound{}
		}
	} else if parent.Recipient.CircleId != 0 {
		res := db.Delete(&dbModel.RecipeCircle{}, id.RecipeAccessId)
		if res.Error != nil {
			return ConvertGormError(res.Error)
		}
		if res.RowsAffected == 0 {
			return repository.ErrNotFound{}
		}
	} else {
		return repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	return nil
}

func (repo *Client) GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	db := repo.db.WithContext(ctx)

	var recipeUser dbModel.RecipeUser
	res := db.Where("recipe_id = ? AND recipe_user_id = ?", parent.RecipeId.RecipeId, id.RecipeAccessId).First(&recipeUser)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return model.RecipeAccess{}, repository.ErrNotFound{}
		}
		return model.RecipeAccess{}, res.Error
	}

	return model.RecipeAccess{
		Level: convert.PermissionLevelToRecipeAccessLevel(recipeUser.PermissionLevel),
	}, nil
}

func (repo *Client) ListRecipeAccesses(ctx context.Context, parent model.RecipeAccessParent, pageSize int64, pageOffset int64, filterStr string) ([]model.RecipeAccess, error) {
	userConversion, err := repo.recipeAccessUserSQLConverter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	circleConversion, err := repo.recipeAccessCircleSQLConverter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	queryUsers := true
	queryCircles := true

	_, userFieldOk := userConversion.UsedColumns[dbModel.RecipeUserFields.UserId]
	_, circleFieldOk := circleConversion.UsedColumns[dbModel.RecipeCircleFields.CircleId]

	if userFieldOk && !circleFieldOk {
		queryCircles = false
	}

	if circleFieldOk && !userFieldOk {
		queryUsers = false
	}

	if userFieldOk && circleFieldOk {
		return nil, repository.ErrInvalidArgument{Msg: "both user and circle filters are not supported"}
	}

	accesses := make([]model.RecipeAccess, 0)

	if queryUsers {
		var recipeUsers []dbModel.RecipeUser
		dbForUsers := repo.db.WithContext(ctx).Model(&dbModel.RecipeUser{})
		if userConversion.WhereClause != "" {
			dbForUsers = dbForUsers.Where(userConversion.WhereClause, userConversion.Params...)
		}
		if parent.RecipeId.RecipeId != 0 {
			dbForUsers.Where("recipe_user.recipe_id = ?", parent.RecipeId.RecipeId)
		}
		if parent.Issuer.UserId != 0 {
			dbForUsers.Joins("JOIN recipe_user as ru ON ru.user_id = ? AND (ru.permission_level >= ? OR ru.user_id = recipe_user.user_id) AND ru.recipe_id = recipe_user.recipe_id", parent.Issuer.UserId, pb.Access_LEVEL_WRITE)
		} else if parent.Issuer.CircleId != 0 {
			dbForUsers.Joins("JOIN recipe_circle as rc ON rc.circle_id = ? AND (rc.permission_level >= ? OR rc.circle_id = recipe_circle.circle_id) AND rc.recipe_id = recipe_circle.recipe_id", parent.Issuer.CircleId, pb.Access_LEVEL_WRITE)
		} else {
			return nil, repository.ErrInvalidArgument{Msg: "issuer is required"}
		}
		err = dbForUsers.Limit(int(pageSize)).
			Offset(int(pageOffset)).
			Find(&recipeUsers).Error
		if err != nil {
			return nil, ConvertGormError(err)
		}
		for _, user := range recipeUsers {
			accesses = append(accesses, convert.RecipeUserToCoreRecipeAccess(user))
		}
	}

	if queryCircles {
		var recipeCircles []dbModel.RecipeCircle
		dbForCircles := repo.db.WithContext(ctx).Model(&dbModel.RecipeCircle{})
		if circleConversion.WhereClause != "" {
			dbForCircles = dbForCircles.Where(circleConversion.WhereClause, circleConversion.Params...)
		}
		if parent.Issuer.UserId != 0 {
			dbForCircles.Joins("JOIN recipe_user as ru ON ru.user_id = ? AND (ru.permission_level >= ? OR ru.user_id = recipe_circle.circle_id) AND ru.recipe_id = recipe_circle.recipe_id", parent.Issuer.UserId, pb.Access_LEVEL_WRITE)
		} else if parent.Issuer.CircleId != 0 {
			dbForCircles.Joins("JOIN recipe_circle as rc ON rc.circle_id = ? AND (rc.permission_level >= ? OR rc.circle_id = recipe_circle.circle_id) AND rc.recipe_id = recipe_circle.recipe_id", parent.Issuer.CircleId, pb.Access_LEVEL_WRITE)
		} else {
			return nil, repository.ErrInvalidArgument{Msg: "issuer is required"}
		}
		err = dbForCircles.Where("recipe_circle.recipe_id = ?", parent.RecipeId.RecipeId).
			Limit(int(pageSize)).
			Offset(int(pageOffset)).
			Find(&recipeCircles).Error
		if err != nil {
			return nil, ConvertGormError(err)
		}
		for _, circle := range recipeCircles {
			accesses = append(accesses, convert.RecipeCircleToCoreRecipeAccess(circle))
		}
	}

	return accesses, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
