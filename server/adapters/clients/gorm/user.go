package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/filtering"
	"github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START userCreateImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// CreateUser creates a new user.
func (repo *Client) CreateUser(ctx context.Context, m cmodel.User) (cmodel.User, error) {
	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	fields := masks.RemovePaths(gmodel.UserFields.Mask())

	err = repo.db.WithContext(ctx).
		Select(fields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// GetUser gets a user.
func (repo *Client) GetUser(ctx context.Context, id cmodel.UserId, fields []string) (cmodel.User, error) {
	gm := gmodel.User{UserId: id.UserId}

	mask := masks.Map(fields, model.UserMap)
	if len(mask) == 0 {
		mask = gmodel.UserFields.Mask()
	}

	err := repo.db.WithContext(ctx).
		Select(mask).
		Clauses(clause.Returning{}).
		First(&gm).Error
	if err != nil {
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err := convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// ListUsers lists users.
func (repo *Client) ListUsers(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, pageOffset int64, filter string, fields []string) ([]cmodel.User, error) {
	args := make([]any, 0, 1)

	fields = masks.Map(fields, gmodel.UserMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		for i, field := range fields {
			fields[i] = fmt.Sprintf("u.%s", field)
		}
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			cmodel.UserFields.GoogleId:   filtering.NewSQLField[string](gmodel.UserFields.GoogleId, "="),
			cmodel.UserFields.FacebookId: filtering.NewSQLField[string](gmodel.UserFields.FacebookId, "="),
			cmodel.UserFields.AmazonId:   filtering.NewSQLField[string](gmodel.UserFields.AmazonId, "="),
			cmodel.UserFields.Username:   filtering.NewSQLField[string](gmodel.UserFields.Username, "="),
		})

	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid filter: %v", err)}
	}

	tx = tx.Clauses(filterClause)
	if len(args) > 0 {
		tx = tx.Where(gmodel.User{}, args...)
	}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_id"},
		Desc:   true,
	}}

	tx = tx.Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	var mods []gmodel.User
	tableAlias := fmt.Sprintf("%s AS u", gmodel.User{}.TableName())
	err = tx.Table(tableAlias).Find(&mods).Error
	if err != nil {
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.User, len(mods))
	for i, m := range mods {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			return nil, fmt.Errorf("unable to read user: %v", err)
		}
	}

	return res, nil
}

// UpdateUser updates a user.
func (repo *Client) UpdateUser(ctx context.Context, m cmodel.User, fields []string) (cmodel.User, error) {
	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	mask := masks.Map(fields, gmodel.UserMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

func (repo *Client) DeleteUser(ctx context.Context, id cmodel.UserId) (cmodel.User, error) {
	return cmodel.User{}, nil
}
