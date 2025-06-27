package gorm

import (
	"context"

	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
)

// UserAccessMap maps the core model fields to the database model fields for the unified UserAccess model.
var UserAccessMap = map[string]string{
	model.UserAccessFields.Level:         dbModel.UserAccessFields.PermissionLevel,
	model.UserAccessFields.State:         dbModel.UserAccessFields.State,
	model.UserAccessFields.RecipientUser: dbModel.UserAccessFields.RecipientUserId,
}

func (repo *Client) CreateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error) {
	// db := repo.db.WithContext(ctx)

	// // Validate that exactly one recipient type is set
	// if access.Recipient != 0 {
	// 	return model.UserAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	// }

	// userAccess := convert.UserAccessFromCoreModel(access)
	// res := db.Create(&userAccess)
	// if res.Error != nil {
	// 	if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
	// 		return model.UserAccess{}, repository.ErrNewAlreadyExists{}
	// 	}
	// 	return model.UserAccess{}, res.Error
	// }

	// access.UserAccessId.UserAccessId = userAccess.UserAccessId
	// return access, nil
	return model.UserAccess{}, nil
}

func (repo *Client) DeleteUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) error {
	db := repo.db.WithContext(ctx)

	res := db.Delete(&dbModel.UserAccess{}, id.UserAccessId)
	if res.Error != nil {
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetUserAccess(ctx context.Context, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	// db := repo.db.WithContext(ctx)

	// var userAccess dbModel.UserAccess
	// res := db.Where("user_id = ? AND user_access_id = ?", parent.UserId.UserId, id.UserAccessId).First(&userAccess)
	// if res.Error != nil {
	// 	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 		return model.UserAccess{}, repository.ErrNotFound{}
	// 	}
	// 	return model.UserAccess{}, res.Error
	// }

	// return convert.UserAccessToCoreUserAccess(userAccess), nil
	return model.UserAccess{}, nil
}

func (repo *Client) ListUserAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.UserAccess, error) {
	// conversion, err := repo.userAccessSQLConverter.Convert(filterStr)
	// if err != nil {
	// 	return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	// }

	// var userAccesses []dbModel.UserAccess
	// db := repo.db.WithContext(ctx).Model(&dbModel.UserAccess{})

	// if conversion.WhereClause != "" {
	// 	db = db.Where(conversion.WhereClause, conversion.Params...)
	// }

	// // Filter by user ID
	// if parent.UserId.UserId != 0 {
	// 	db = db.Where("user_access.user_id = ?", parent.UserId.UserId)
	// }

	// // Add authorization check - only allow access if the requester has write permission or is the recipient
	// if parent.Requester.UserId != 0 {
	// 	db = db.Where(`
	// 		EXISTS (
	// 			SELECT 1 FROM user_access ra
	// 			WHERE ra.user_id = ?
	// 			AND ra.permission_level >= ?
	// 			AND ra.user_id = user_access.user_id
	// 		) OR user_access.user_id = ?`,
	// 		parent.Requester.UserId, permPb.PermissionLevel_PERMISSION_LEVEL_WRITE, parent.Requester.UserId)
	// } else if parent.Requester.UserId != 0 {
	// 	db = db.Where(`
	// 		EXISTS (
	// 			SELECT 1 FROM user_access ra
	// 			WHERE ra.user_id = ?
	// 			AND ra.permission_level >= ?
	// 			AND ra.user_id = user_access.user_id
	// 		) OR user_access.user_id = ?`,
	// 		parent.Requester.UserId, permPb.PermissionLevel_PERMISSION_LEVEL_WRITE, parent.Requester.UserId)
	// } else {
	// 	return nil, repository.ErrInvalidArgument{Msg: "requester is required"}
	// }

	// err = db.Limit(int(pageSize)).
	// 	Offset(int(pageOffset)).
	// 	Find(&userAccesses).Error
	// if err != nil {
	// 	return nil, ConvertGormError(err)
	// }

	// accesses := make([]model.UserAccess, len(userAccesses))
	// for i, access := range userAccesses {
	// 	accesses[i] = convert.UserAccessToCoreUserAccess(access)
	// }

	// return accesses, nil
	return nil, nil
}

func (repo *Client) UpdateUserAccess(ctx context.Context, access model.UserAccess) (model.UserAccess, error) {
	// dbAccess := convert.CoreUserAccessToUserAccess(access)

	// db := repo.db.WithContext(ctx).Select("state").Clauses(&clause.Returning{})

	// err := db.Where("user_access_id = ?", access.UserAccessId.UserAccessId).Updates(&dbAccess).Error
	// if err != nil {
	// 	return model.UserAccess{}, ConvertGormError(err)
	// }

	// return convert.UserAccessToCoreUserAccess(dbAccess), nil
	return model.UserAccess{}, nil
}
