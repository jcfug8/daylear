package gorm

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	dbModel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	model "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CircleAccessMap maps the core model fields to the database model fields for the unified CircleAccess model.
var CircleAccessMap = map[string]string{
	model.CircleAccessFields.Level:         dbModel.CircleAccessFields.PermissionLevel,
	model.CircleAccessFields.State:         dbModel.CircleAccessFields.State,
	model.CircleAccessFields.RecipientUser: dbModel.CircleAccessFields.RecipientUserId,
}

func (repo *Client) CreateCircleAccess(ctx context.Context, access model.CircleAccess) (model.CircleAccess, error) {
	db := repo.db.WithContext(ctx)

	// Validate that exactly one recipient type is set
	if access.Recipient == 0 {
		return model.CircleAccess{}, repository.ErrInvalidArgument{Msg: "recipient is required"}
	}

	circleAccess := convert.CoreCircleAccessToCircleAccess(access)
	res := db.Create(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return model.CircleAccess{}, repository.ErrNewAlreadyExists{}
		}
		return model.CircleAccess{}, res.Error
	}

	access.CircleAccessId.CircleAccessId = circleAccess.CircleAccessId
	return access, nil
}

func (repo *Client) DeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) error {
	db := repo.db.WithContext(ctx)

	res := db.Delete(&dbModel.CircleAccess{}, id.CircleAccessId)
	if res.Error != nil {
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) BulkDeleteCircleAccess(ctx context.Context, parent model.CircleAccessParent) error {
	db := repo.db.WithContext(ctx)

	res := db.Where("circle_id = ?", parent.CircleId.CircleId).Delete(&dbModel.CircleAccess{})
	if res.Error != nil {
		return ConvertGormError(res.Error)
	}
	if res.RowsAffected == 0 {
		return repository.ErrNotFound{}
	}

	return nil
}

func (repo *Client) GetCircleAccess(ctx context.Context, parent model.CircleAccessParent, id model.CircleAccessId) (model.CircleAccess, error) {
	db := repo.db.WithContext(ctx)

	var circleAccess dbModel.CircleAccess
	res := db.Where("circle_id = ? AND circle_access_id = ?", parent.CircleId.CircleId, id.CircleAccessId).First(&circleAccess)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return model.CircleAccess{}, repository.ErrNotFound{}
		}
		return model.CircleAccess{}, res.Error
	}

	return convert.CircleAccessToCoreCircleAccess(circleAccess), nil
}

func (repo *Client) ListCircleAccesses(ctx context.Context, authAccount cmodel.AuthAccount, parent model.CircleAccessParent, pageSize int32, pageOffset int64, filterStr string) ([]model.CircleAccess, error) {
	if authAccount.UserId == 0 {
		return nil, repository.ErrInvalidArgument{Msg: "user id is required"}
	}

	conversion, err := repo.circleAccessSQLConverter.Convert(filterStr)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	var circleAccesses []dbModel.CircleAccess
	db := repo.db.WithContext(ctx).Model(&dbModel.CircleAccess{})

	if conversion.WhereClause != "" {
		db = db.Where(conversion.WhereClause, conversion.Params...)
	}

	// Filter by circle ID if provided
	if parent.CircleId.CircleId != 0 {
		db = db.Where("circle_access.circle_id = ?", parent.CircleId.CircleId)
	}

	db = db.Where(
		"circle_access.recipient_user_id = ? OR circle_access.circle_id IN (SELECT circle_id FROM circle_access WHERE recipient_user_id = ? AND permission_level >= ?)",
		authAccount.UserId, authAccount.UserId, types.PermissionLevel_PERMISSION_LEVEL_WRITE,
	)

	err = db.Limit(int(pageSize)).
		Offset(int(pageOffset)).
		Find(&circleAccesses).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	accesses := make([]model.CircleAccess, len(circleAccesses))
	for i, access := range circleAccesses {
		accesses[i] = convert.CircleAccessToCoreCircleAccess(access)
	}

	return accesses, nil
}

func (repo *Client) UpdateCircleAccess(ctx context.Context, access model.CircleAccess) (model.CircleAccess, error) {
	dbAccess := convert.CoreCircleAccessToCircleAccess(access)

	db := repo.db.WithContext(ctx).Select("permission_level", "state").Clauses(clause.Returning{})

	err := db.Where("circle_access_id = ?", access.CircleAccessId.CircleAccessId).Updates(&dbAccess).Error
	if err != nil {
		return model.CircleAccess{}, ConvertGormError(err)
	}

	return convert.CircleAccessToCoreCircleAccess(dbAccess), nil
}
