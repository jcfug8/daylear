package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
)

const (
	AccessKeyTable = "access_key"
)

const (
	AccessKeyColumn_AccessKeyId        = "access_key_id"
	AccessKeyColumn_UserId             = "user_id"
	AccessKeyColumn_Title              = "title"
	AccessKeyColumn_Description        = "description"
	AccessKeyColumn_EncryptedAccessKey = "encrypted_access_key"
	AccessKeyColumn_CreateTime         = "create_time"
	AccessKeyColumn_UpdateTime         = "update_time"
)

var AccessKeyFieldMasker = fieldmask.NewSQLFieldMasker(AccessKey{}, map[string][]fieldmask.Field{
	cmodel.AccessKeyField_Parent:               {{Name: AccessKeyColumn_UserId, Table: AccessKeyTable}},
	cmodel.AccessKeyField_AccessKeyId:          {{Name: AccessKeyColumn_AccessKeyId, Table: AccessKeyTable}},
	cmodel.AccessKeyField_Title:                {{Name: AccessKeyColumn_Title, Table: AccessKeyTable, Updatable: true}},
	cmodel.AccessKeyField_Description:          {{Name: AccessKeyColumn_Description, Table: AccessKeyTable, Updatable: true}},
	cmodel.AccessKeyField_UnencryptedAccessKey: {{Name: AccessKeyColumn_EncryptedAccessKey, Table: AccessKeyTable}},
	cmodel.AccessKeyField_CreateTime:           {{Name: AccessKeyColumn_CreateTime, Table: AccessKeyTable}},
	cmodel.AccessKeyField_UpdateTime:           {{Name: AccessKeyColumn_UpdateTime, Table: AccessKeyTable}},
})

var AccessKeySQLConverter = filter.NewSQLConverter(map[string]filter.Field{}, true)

// AccessKey is the GORM model for an access key.
type AccessKey struct {
	AccessKeyId        int64     `gorm:"primaryKey;column:access_key_id;autoIncrement;<-:false"`
	UserId             int64     `gorm:"column:user_id;not null;index"`
	Title              string    `gorm:"column:title;not null"`
	Description        string    `gorm:"column:description"`
	EncryptedAccessKey string    `gorm:"column:encrypted_access_key;not null"`
	CreateTime         time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime         time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// TableName sets the table name for the AccessKey model.
func (AccessKey) TableName() string {
	return AccessKeyTable
}
