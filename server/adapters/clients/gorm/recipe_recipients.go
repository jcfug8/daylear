package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
)

// GetRecipeRecipient retrieves a single recipe recipient (either user or circle) for a given recipe.
// It returns a RecipeRecipient object containing the user or circle information along with their
// permission level and title.
func (c *Client) GetRecipeRecipient(ctx context.Context, parent cmodel.RecipeParent, id cmodel.RecipeId) (cmodel.RecipeRecipient, error) {
	var recipeUser model.RecipeUser
	var recipeCircle model.RecipeCircle
	title := ""

	if id.RecipeId == 0 {
		return cmodel.RecipeRecipient{}, repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	if parent.CircleId != 0 {
		if err := c.db.WithContext(ctx).
			Table("recipe_circle rc").
			Select("rc.*, c.title").
			Joins("Join circles c ON rc.circle_id = c.circle_id").
			Where("rc.recipe_id = ? AND rc.circle_id = ?", id.RecipeId, parent.CircleId).
			Find(&recipeCircle).Error; err != nil {
			return cmodel.RecipeRecipient{}, ConvertGormError(err)
		}
		title = recipeCircle.Title
	} else if parent.UserId != 0 {
		if err := c.db.WithContext(ctx).
			Table("recipe_user ru").
			Select("ru.*, u.username as title").
			Joins("JOIN daylear_user u ON ru.user_id = u.user_id").
			Where("ru.recipe_id = ? AND ru.user_id = ?", id.RecipeId, parent.UserId).
			Find(&recipeUser).Error; err != nil {
			return cmodel.RecipeRecipient{}, ConvertGormError(err)
		}
		title = recipeUser.Title
	} else {
		return cmodel.RecipeRecipient{}, repository.ErrInvalidArgument{Msg: "invalid parent"}
	}

	return cmodel.RecipeRecipient{
		RecipeId:        id,
		RecipeParent:    parent,
		PermissionLevel: recipeUser.PermissionLevel,
		Title:           title,
	}, nil
}

// ListRecipeRecipients retrieves all recipe recipients (both users and circles) for a given recipe.
// It returns a list of RecipeRecipient objects containing user and circle information along with their
// permission levels and titles.
func (c *Client) ListRecipeRecipients(ctx context.Context, id cmodel.RecipeId) ([]cmodel.RecipeRecipient, error) {
	var recipeUsers []model.RecipeUser

	var recipeCircles []model.RecipeCircle

	if err := c.db.WithContext(ctx).
		Table("recipe_user ru").
		Select("ru.*, u.username as title").
		Joins("JOIN daylear_user u ON ru.user_id = u.user_id").
		Where("ru.recipe_id = ?", id.RecipeId).
		Find(&recipeUsers).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	if err := c.db.WithContext(ctx).
		Table("recipe_circle rc").
		Select("rc.*, c.title").
		Joins("JOIN circle c ON rc.circle_id = c.circle_id").
		Where("rc.recipe_id = ?", id.RecipeId).
		Find(&recipeCircles).Error; err != nil {
		return nil, ConvertGormError(err)
	}

	result := make([]cmodel.RecipeRecipient, len(recipeUsers)+len(recipeCircles))
	for i, ru := range recipeUsers {
		result[i] = cmodel.RecipeRecipient{
			RecipeId: cmodel.RecipeId{
				RecipeId: ru.RecipeId,
			},
			RecipeParent: cmodel.RecipeParent{
				UserId: ru.UserId,
			},
			PermissionLevel: ru.PermissionLevel,
			Title:           ru.Title,
		}
	}

	for i, rc := range recipeCircles {
		result[i] = cmodel.RecipeRecipient{
			RecipeId: cmodel.RecipeId{
				RecipeId: rc.RecipeId,
			},
			RecipeParent: cmodel.RecipeParent{
				CircleId: rc.CircleId,
			},
			PermissionLevel: rc.PermissionLevel,
			Title:           rc.Title,
		}
	}

	return result, nil
}

// BulkCreateRecipeRecipients creates recipe recipient entries for the specified recipe.
// It takes a list of parents (users and circles) and creates their corresponding recipe recipient entries
// with the specified permission level. If a parent is invalid (neither user nor circle), it returns an error.
func (repo *Client) BulkCreateRecipeRecipients(ctx context.Context, parents []cmodel.RecipeParent, id cmodel.RecipeId, permission permPb.PermissionLevel) error {
	if id.RecipeId == 0 {
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}
	if len(parents) == 0 {
		return nil
	}

	var recipeUsers []model.RecipeUser
	var recipeCircles []model.RecipeCircle

	for _, parent := range parents {
		if parent.CircleId != 0 {
			recipeCircles = append(recipeCircles, model.RecipeCircle{
				RecipeId:        id.RecipeId,
				CircleId:        parent.CircleId,
				PermissionLevel: permission,
			})
		} else if parent.UserId != 0 {
			recipeUsers = append(recipeUsers, model.RecipeUser{
				RecipeId:        id.RecipeId,
				UserId:          parent.UserId,
				PermissionLevel: permission,
			})
		} else {
			return repository.ErrInvalidArgument{Msg: "invalid parent"}
		}
	}

	if len(recipeUsers) > 0 {
		if err := repo.db.WithContext(ctx).Create(&recipeUsers).Error; err != nil {
			return ConvertGormError(err)
		}
	}

	if len(recipeCircles) > 0 {
		if err := repo.db.WithContext(ctx).Create(&recipeCircles).Error; err != nil {
			return ConvertGormError(err)
		}
	}

	return nil
}

// BulkDeleteRecipeRecipients deletes recipe recipients (both users and circles) for the specified recipe.
// It takes a list of parents (users and circles) and deletes their corresponding recipe recipient entries
// with the specified permission level. If a parent is invalid (neither user nor circle), it returns an error.
func (repo *Client) BulkDeleteRecipeRecipients(ctx context.Context, parents []cmodel.RecipeParent, id cmodel.RecipeId) error {
	if id.RecipeId == 0 {
		return repository.ErrInvalidArgument{Msg: "recipe id is required"}
	}
	if len(parents) == 0 {
		return nil
	}

	var userIds []int64
	var circleIds []int64

	for _, parent := range parents {
		if parent.CircleId != 0 {
			circleIds = append(circleIds, parent.CircleId)
		} else if parent.UserId != 0 {
			userIds = append(userIds, parent.UserId)
		} else {
			return repository.ErrInvalidArgument{Msg: "invalid parent"}
		}
	}

	if len(userIds) > 0 {
		if err := repo.db.WithContext(ctx).Where("recipe_id = ? AND user_id IN ?", id.RecipeId, userIds).Delete(&model.RecipeUser{}).Error; err != nil {
			return ConvertGormError(err)
		}
	}

	if len(circleIds) > 0 {
		if err := repo.db.WithContext(ctx).Where("recipe_id = ? AND circle_id IN ?", id.RecipeId, circleIds).Delete(&model.RecipeCircle{}).Error; err != nil {
			return ConvertGormError(err)
		}
	}

	return nil
}
