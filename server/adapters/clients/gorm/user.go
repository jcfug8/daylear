package gorm

import (
	"context"
	"fmt"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/masks"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
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
	"visibility":  gmodel.UserFields.Visibility,
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

// GetUser gets a user.
func (repo *Client) GetUser(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.UserId) (cmodel.User, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx)
	log.Info().Msg("GORM GetUser called")

	gm := gmodel.User{UserId: id.UserId}

	err := repo.db.WithContext(ctx).
		Select("daylear_user.*", "user_access.permission_level", "user_access.state", "user_access.user_access_id").
		Joins("LEFT JOIN user_access ON daylear_user.user_id = user_access.recipient_user_id AND user_access.recipient_user_id = ?", authAccount.UserId).
		Where("daylear_user.user_id = ? AND (daylear_user.visibility = ? OR user_access.recipient_user_id = ?)", id.UserId, types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId).
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

	orders := []clause.OrderByColumn{{
		Column: clause.Column{Name: "user_id"},
		Desc:   true,
	}}

	tx := repo.db.WithContext(ctx).
		Select("daylear_user.*", "user_access.permission_level", "user_access.state", "user_access.user_access_id").
		Joins("LEFT JOIN user_access ON daylear_user.user_id = user_access.recipient_user_id AND user_access.recipient_user_id = ?", authAccount.UserId).
		Where("daylear_user.visibility = ? OR user_access.recipient_user_id = ?", types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC, authAccount.UserId).
		Order(clause.OrderBy{Columns: orders}).
		Limit(int(pageSize)).
		Offset(int(pageOffset))

	conversion, err := repo.userSQLConverter.Convert(filter)
	if err != nil {
		log.Error().Err(err).Msg("userSQLConverter.Convert failed")
		return nil, repository.ErrInvalidArgument{Msg: "invalid filter: " + err.Error()}
	}

	if conversion.WhereClause != "" {
		tx = tx.Where(conversion.WhereClause, conversion.Params...)
	}

	err = tx.Find(&dbUsers).Error
	if err != nil {
		log.Error().Err(err).Msg("db.Find failed")
		return nil, ConvertGormError(err)
	}

	res := make([]cmodel.User, len(dbUsers))
	for i, m := range dbUsers {
		res[i], err = convert.UserToCoreModel(m)
		if err != nil {
			log.Error().Err(err).Msg("unable to read user")
			return nil, fmt.Errorf("unable to read user: %v", err)
		}
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
