package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// AccessKeyFromCoreModel converts a core AccessKey to a GORM AccessKey
func AccessKeyFromCoreModel(accessKey cmodel.AccessKey) (gmodel.AccessKey, error) {
	return gmodel.AccessKey{
		AccessKeyId:        accessKey.AccessKeyId.AccessKeyId,
		UserId:             accessKey.Parent.UserId,
		Title:              accessKey.Title,
		Description:        accessKey.Description,
		EncryptedAccessKey: accessKey.EncryptedAccessKey,
		CreateTime:         accessKey.CreateTime,
		UpdateTime:         accessKey.UpdateTime,
	}, nil
}

// AccessKeyToCoreModel converts a GORM AccessKey to a core AccessKey
func AccessKeyToCoreModel(gormAccessKey gmodel.AccessKey) (cmodel.AccessKey, error) {
	accessKey := cmodel.AccessKey{
		AccessKeyId: cmodel.AccessKeyId{AccessKeyId: gormAccessKey.AccessKeyId},
		Parent: cmodel.AccessKeyParent{
			UserId: gormAccessKey.UserId,
		},
		Title:              gormAccessKey.Title,
		Description:        gormAccessKey.Description,
		EncryptedAccessKey: gormAccessKey.EncryptedAccessKey,
		CreateTime:         gormAccessKey.CreateTime,
		UpdateTime:         gormAccessKey.UpdateTime,
	}

	return accessKey, nil
}
