package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateUser creates a new user.
func (repo *Client) CreateUser(ctx context.Context, m cmodel.User, fields []string) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid user when creating user row")
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.UserFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create user row")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid user row when creating user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// GetUser gets a user. TODO: the WHERE clause is not correct yet.
func (repo *Client) GetUser(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.UserId, fields []string) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", id.UserId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.User{}

	err := repo.db.WithContext(ctx).
		Select(gmodel.UserFieldMasker.Convert(
			fields,
			fieldmask.ExcludeKeys(
				cmodel.UserField_AccessName,
				cmodel.UserField_AccessPermissionLevel,
				cmodel.UserField_AccessState,
			),
		)).
		Where("daylear_user.user_id = ?", id.UserId).
		First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get user row")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err := convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid user row when getting user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// ListUsers lists users.
func (repo *Client) ListUsers(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, pageOffset int64, filter string, fields []string) ([]cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Str("filter", filter).
		Strs("fields", fields).
		Int("pageSize", int(pageSize)).
		Int64("pageOffset", pageOffset).
		Logger()

	dbUsers := []gmodel.User{}

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	converter := gmodel.UserSQLConverter
	switch {
	case authAccount.CircleId != 0:
		converter = gmodel.UserCircleSQLConverter
		tx = tx.Select(gmodel.UserFieldMasker.Convert(fields, fieldmask.ExcludeTables("user_access"))).
			Joins("LEFT JOIN circle_access ON daylear_user.user_id = circle_access.recipient_user_id AND circle_access.circle_id = ?", authAccount.CircleId).
			Joins("LEFT JOIN user_access as user_access_auth ON daylear_user.user_id = user_access_auth.user_id AND user_access_auth.recipient_user_id = ?", authAccount.AuthUserId)
	case authAccount.AuthUserId != authAccount.UserId && authAccount.UserId != 0:
		tx = tx.Select(gmodel.UserFieldMasker.Convert(fields, fieldmask.ExcludeTables("circle_access"))).
			Joins("LEFT JOIN user_access ON daylear_user.user_id = user_access.user_id AND user_access.recipient_user_id = ?", authAccount.UserId).
			Joins("LEFT JOIN user_access as user_access_auth ON daylear_user.user_id = user_access_auth.user_id AND user_access_auth.recipient_user_id = ?", authAccount.AuthUserId)
	default:
		tx = tx.Select(gmodel.UserFieldMasker.Convert(fields, fieldmask.ExcludeTables("circle_access"))).
			Joins("LEFT JOIN user_access ON daylear_user.user_id = user_access.user_id AND user_access.recipient_user_id = ?", authAccount.AuthUserId)
	}

	conversion, err := converter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("invalid filter string when listing user rows")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbUsers).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list user rows")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.User, len(dbUsers))
	for i, m := range dbUsers {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("invalid user row when listing users")
			return nil, fmt.Errorf("unable to read user: %v", err)
		}
		res[i].Parent.CircleId = authAccount.CircleId
	}

	return res, nil
}

// UpdateUser updates a user.
func (repo *Client) UpdateUser(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.User, fields []string) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", m.Id.UserId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid user when updating user row")
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.UserFieldMasker.Convert(fields, fieldmask.OnlyUpdatable())).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update user row")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid user row when updating user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

func (repo *Client) DeleteUser(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.UserId) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", id.UserId).
		Logger()

	log.Error().Msg("delete user operation not implemented")
	return cmodel.User{}, repository.ErrNewUnimplemented{Msg: "delete user operation not implemented"}
}
