package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
)

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
	// User Accesses
	userFilter := filter.NewSQLConverter(model.RecipeAccessMap.ToStringMap())
	userWhere, err := userFilter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	var recipeUsers []dbModel.RecipeUser
	dbForUsers := repo.db.WithContext(ctx).Model(&dbModel.RecipeUser{})
	if userWhere != "" {
		dbForUsers = dbForUsers.Where(userWhere, userFilter.Params...)
	}
	err = dbForUsers.Where("recipe_id = ?", parent.RecipeId.RecipeId).
		Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&recipeUsers).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	// Circle Accesses
	circleFilter := filter.NewSQLConverter(model.RecipeAccessMap.ToStringMap())
	circleWhere, err := circleFilter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter"}
	}

	var recipeCircles []dbModel.RecipeCircle
	dbForCircles := repo.db.WithContext(ctx).Model(&dbModel.RecipeCircle{})
	if circleWhere != "" {
		dbForCircles = dbForCircles.Where(circleWhere, circleFilter.Params...)
	}
	err = dbForCircles.Where("recipe_id = ?", parent.RecipeId.RecipeId).
		Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&recipeCircles).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	// Convert and combine
	accesses := make([]model.RecipeAccess, 0, len(recipeUsers)+len(recipeCircles))
	for _, user := range recipeUsers {
		accesses = append(accesses, convert.RecipeUserToCoreRecipeAccess(user))
	}
	for _, circle := range recipeCircles {
		accesses = append(accesses, convert.RecipeCircleToCoreRecipeAccess(circle))
	}

	return accesses, nil
}

func (repo *Client) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
