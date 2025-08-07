package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

var UserMap = map[string]string{
	"google_id":   gmodel.UserFields.GoogleId,
	"facebook_id": gmodel.UserFields.FacebookId,
	"amazon_id":   gmodel.UserFields.AmazonId,
	"permission":  gmodel.UserFields.Permission,
	"state":       gmodel.UserFields.State,
	"username":    gmodel.UserFields.Username,
}

var UserCircleMap = map[string]string{
	"google_id":   gmodel.UserFields.GoogleId,
	"facebook_id": gmodel.UserFields.FacebookId,
	"amazon_id":   gmodel.UserFields.AmazonId,
	"username":    gmodel.UserFields.Username,
	"permission":  gmodel.UserFields.Permission,
	"state":       gmodel.UserFields.State,
}

// CreateUser creates a new user.
func (repo *Client) CreateUser(ctx context.Context, m cmodel.User) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid user model")
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	fields := masks.RemovePaths(gmodel.UserFields.Mask())

	err = repo.db.WithContext(ctx).
		Select(fields).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Create failed")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// GetUser gets a user. TODO: the WHERE clause is not correct yet.
func (repo *Client) GetUser(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.UserId) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM GetUser called")

	gm := gmodel.User{}

	err := repo.db.WithContext(ctx).
		Select("daylear_user.*", "user_access.permission_level", "user_access.state", "user_access.user_access_id", "user_access.requester_user_id").
		Joins("LEFT JOIN user_access ON daylear_user.user_id = user_access.user_id AND user_access.recipient_user_id = ?", authAccount.AuthUserId).
		Where("daylear_user.user_id = ?", id.UserId).
		First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.First failed")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err := convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

// ListUsers lists users.
func (repo *Client) ListUsers(ctx context.Context, authAccount cmodel.AuthAccount, pageSize int32, pageOffset int64, filter string, fields []string) ([]cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM ListUsers called")
	dbUsers := []gmodel.User{}

	columns := []string{"daylear_user.*"}
	joins := []string{}
	joinParams := [][]any{}
	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_id"},
		Desc:   true,
	}}
	converter := repo.userSQLConverter

	switch {
	case authAccount.CircleId != 0:
		converter = repo.userCircleSQLConverter
		columns = append(columns, "circle_access.permission_level", "circle_access.state", "circle_access.circle_access_id", "circle_access.requester_user_id")

		joins = append(joins, "LEFT JOIN circle_access ON daylear_user.user_id = circle_access.recipient_user_id AND circle_access.circle_id = ?")
		joinParams = append(joinParams, []any{authAccount.CircleId})

		joins = append(joins, "LEFT JOIN user_access as user_access_auth ON daylear_user.user_id = user_access_auth.user_id AND user_access_auth.recipient_user_id = ?")
		joinParams = append(joinParams, []any{authAccount.AuthUserId})
	case authAccount.AuthUserId != authAccount.UserId:
		columns = append(columns, "user_access_auth.permission_level", "user_access_auth.state", "user_access_auth.user_access_id", "user_access_auth.requester_user_id")

		joins = append(joins, "LEFT JOIN user_access ON daylear_user.user_id = user_access.user_id AND user_access.recipient_user_id = ?")
		joinParams = append(joinParams, []any{authAccount.UserId})

		joins = append(joins, "LEFT JOIN user_access as user_access_auth ON daylear_user.user_id = user_access_auth.user_id AND user_access_auth.recipient_user_id = ?")
		joinParams = append(joinParams, []any{authAccount.AuthUserId})
	case authAccount.UserId == authAccount.AuthUserId && authAccount.AuthUserId != 0:
		columns = append(columns, "user_access.permission_level", "user_access.state", "user_access.user_access_id", "user_access.requester_user_id")

		joins = append(joins, "LEFT JOIN user_access ON daylear_user.user_id = user_access.user_id AND user_access.recipient_user_id = ?")
		joinParams = append(joinParams, []any{authAccount.UserId})
	}

	tx := repo.db.WithContext(ctx).
		Select(columns).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	for i, join := range joins {
		tx = tx.Joins(join, joinParams[i]...)
	}

	conversion, err := converter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert list users filter")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbUsers).Error
	if err != nil {
		log.Error().Err(err).Msg("list users failed")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.User, len(dbUsers))
	for i, m := range dbUsers {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("unable to read user")
			return nil, fmt.Errorf("unable to read user: %v", err)
		}
		res[i].Parent.CircleId = authAccount.CircleId
	}

	log.Info().Msg("GORM ListUsers returning successfully")
	return res, nil
}

// UpdateUser updates a user.
func (repo *Client) UpdateUser(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.User, fields []string) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM UpdateUser called")

	gm, err := convert.UserFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid user model")
		return cmodel.User{}, repository.ErrInvalidArgument{Msg: fmt.Sprintf("invalid user: %v", err)}
	}

	mask := masks.Map(fields, gmodel.UserMap)

	err = repo.db.WithContext(ctx).
		Select(mask).
		Clauses(&clause.Returning{}).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Updates failed")
		return cmodel.User{}, ConvertGormError(err)
	}

	m, err = convert.UserToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("unable to read user")
		return cmodel.User{}, fmt.Errorf("unable to read user: %v", err)
	}

	return m, nil
}

func (repo *Client) DeleteUser(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.UserId) (cmodel.User, error) {
	return cmodel.User{}, nil
}
