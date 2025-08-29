package gorm

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/clients/gorm/convert"
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	"github.com/jcfug8/daylear/server/core/logutil"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/repository"
	"gorm.io/gorm/clause"
)

// CreateAccessKey creates a new access key
func (repo *Client) CreateAccessKey(ctx context.Context, m cmodel.AccessKey, fields []string) (cmodel.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Strs("fields", fields).
		Logger()

	gm, err := convert.AccessKeyFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key when creating access key row")
		return cmodel.AccessKey{}, repository.ErrInvalidArgument{Msg: "invalid access key"}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.AccessKeyFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Create(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to create access key row")
		return cmodel.AccessKey{}, ConvertGormError(err)
	}

	m, err = convert.AccessKeyToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key row when creating access key")
		return cmodel.AccessKey{}, repository.ErrInternal{Msg: "invalid access key row when creating access key"}
	}

	return m, nil
}

// DeleteAccessKey deletes an access key
func (repo *Client) DeleteAccessKey(ctx context.Context, id cmodel.AccessKeyId) (cmodel.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("accessKeyId", id.AccessKeyId).
		Logger()

	gm := gmodel.AccessKey{AccessKeyId: id.AccessKeyId}

	err := repo.db.WithContext(ctx).
		Select(gmodel.AccessKeyFieldMasker.Get()).
		Clauses(clause.Returning{}).
		Delete(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to delete access key row")
		return cmodel.AccessKey{}, ConvertGormError(err)
	}

	m, err := convert.AccessKeyToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key row when deleting access key")
		return cmodel.AccessKey{}, repository.ErrInternal{Msg: "invalid access key row when deleting access key"}
	}

	return m, nil
}

// GetAccessKey retrieves an access key
func (repo *Client) GetAccessKey(ctx context.Context, authAccount cmodel.AuthAccount, id cmodel.AccessKeyId, fields []string) (cmodel.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("accessKeyId", id.AccessKeyId).
		Strs("fields", fields).
		Logger()

	gm := gmodel.AccessKey{}

	err := repo.db.WithContext(ctx).
		Select(gmodel.AccessKeyFieldMasker.Convert(fields)).
		Where("access_key_id = ?", id.AccessKeyId).
		First(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to get access key row")
		return cmodel.AccessKey{}, ConvertGormError(err)
	}

	m, err := convert.AccessKeyToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key row when getting access key")
		return cmodel.AccessKey{}, repository.ErrInternal{Msg: "invalid access key row when getting access key"}
	}

	return m, nil
}

// ListAccessKeys lists access keys for a user
func (repo *Client) ListAccessKeys(ctx context.Context, authAccount cmodel.AuthAccount, userId int64, pageSize int32, offset int64, filter string, fields []string) ([]cmodel.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("userId", userId).
		Int32("pageSize", pageSize).
		Int64("offset", offset).
		Str("filter", filter).
		Strs("fields", fields).
		Logger()

	var gAccessKeys []gmodel.AccessKey

	tx := repo.db.WithContext(ctx).
		Select(gmodel.AccessKeyFieldMasker.Convert(fields)).
		Where("user_id = ?", userId).
		Order("create_time DESC").
		Limit(int(pageSize)).
		Offset(int(offset))

	err := tx.Find(&gAccessKeys).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to list access key rows")
		return nil, ConvertGormError(err)
	}

	accessKeys := make([]cmodel.AccessKey, len(gAccessKeys))
	for i, gAccessKey := range gAccessKeys {
		accessKey, err := convert.AccessKeyToCoreModel(gAccessKey)
		if err != nil {
			log.Error().Err(err).Msg("invalid access key row when listing access keys")
			return nil, repository.ErrInternal{Msg: "invalid access key row when listing access keys"}
		}
		accessKeys[i] = accessKey
	}

	return accessKeys, nil
}

// UpdateAccessKey updates an access key
func (repo *Client) UpdateAccessKey(ctx context.Context, authAccount cmodel.AuthAccount, m cmodel.AccessKey, fields []string) (cmodel.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(repo.log, ctx).With().
		Int64("accessKeyId", m.AccessKeyId.AccessKeyId).
		Strs("fields", fields).
		Logger()

	gm, err := convert.AccessKeyFromCoreModel(m)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key when updating access key row")
		return cmodel.AccessKey{}, repository.ErrInvalidArgument{Msg: "invalid access key"}
	}

	err = repo.db.WithContext(ctx).
		Select(gmodel.AccessKeyFieldMasker.Convert(fields)).
		Clauses(clause.Returning{}).
		Where("access_key_id = ?", gm.AccessKeyId).
		Updates(&gm).Error
	if err != nil {
		log.Error().Err(err).Msg("unable to update access key row")
		return cmodel.AccessKey{}, ConvertGormError(err)
	}

	m, err = convert.AccessKeyToCoreModel(gm)
	if err != nil {
		log.Error().Err(err).Msg("invalid access key row when updating access key")
		return cmodel.AccessKey{}, repository.ErrInternal{Msg: "invalid access key row when updating access key"}
	}

	return m, nil
}
