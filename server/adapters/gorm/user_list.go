package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/gorm/convert"
	"github.com/jcfug8/daylear/server/adapters/gorm/filtering"
	gmodel "github.com/jcfug8/daylear/server/adapters/gorm/model"
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"gorm.io/gorm/clause"
	// IRIOMO:CUSTOM_CODE_SLOT_START userListImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ListUsers lists users.
func (repo *Client) ListUsers(ctx context.Context, page *cmodel.PageToken[cmodel.User], filter string, fields []string) ([]cmodel.User, error) {
	errz := errz.Context("repository.list_users")

	queryModel := gmodel.User{}

	// IRIOMO:CUSTOM_CODE_SLOT_START listUsersNameChecks
	var args []any
	// IRIOMO:CUSTOM_CODE_SLOT_END

	fields = masks.Map(fields, gmodel.UserMap)

	tx := repo.db.WithContext(ctx)
	if len(fields) > 0 {
		tx = tx.Select(fields)
	}

	t := filtering.NewSQLTranspiler(
		map[string]filtering.Field[clause.Expression]{
			cmodel.UserFields.GoogleId:   filtering.NewSQLField[string](gmodel.UserFields.GoogleId, "="),
			cmodel.UserFields.FacebookId: filtering.NewSQLField[string](gmodel.UserFields.FacebookId, "="),
			cmodel.UserFields.AmazonId:   filtering.NewSQLField[string](gmodel.UserFields.AmazonId, "="),
		})

	// IRIOMO:CUSTOM_CODE_SLOT_START listUsersFilterCustomizations
	// comment out this next line if you want to disable AND
	// t.And = nil
	// comment out this next line if you want to disable OR
	// t.Or = nil
	// IRIOMO:CUSTOM_CODE_SLOT_END

	// IRIOMO:CUSTOM_CODE_SLOT_START listUsersFilterInfo
	filterClause, _ /* info */, err := t.Transpile(filter)
	if err != nil {
		return nil, errz.NewInvalidArgument("invalid filter: %v", err)
	}
	// IRIOMO:CUSTOM_CODE_SLOT_END

	tx = tx.Clauses(filterClause)
	if len(args) > 0 {
		tx = tx.Where(queryModel, args...)
	}

	if page != nil {
		var orders []clause.OrderByColumn
		// IRIOMO:CUSTOM_CODE_SLOT_START listUsersOrder

		orders = []clause.OrderByColumn{{
			Column: clause.Column{Name: "user_id"},
			Desc:   true,
		}}

		// IRIOMO:CUSTOM_CODE_SLOT_END

		tx = tx.Order(clause.OrderBy{Columns: orders}).
			Limit(int(page.PageSize)).
			Offset(int(page.Skip))

		if page.Tail != nil {
			tail, err := convert.UserFromCoreModel(*page.Tail)
			if err != nil {
				return nil, errz.Wrapf("unable to read tail: %v", err)
			}

			tx = tx.Where(
				Seek(orders, gmodel.UserFields.Map(tail)))
		}
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START listUsersBefore
	// IRIOMO:CUSTOM_CODE_SLOT_END

	var mods []gmodel.User
	if err = tx.Find(&mods).Error; err != nil {
		return nil, ErrzError(errz, "", err)
	}

	res := make([]cmodel.User, len(mods))
	for i, m := range mods {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			return nil, errz.Wrapf("unable to read user: %v", err)
		}
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START listUsersAfter
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return res, nil
}
