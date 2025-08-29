package model

import (
	"time"
)

// AccessKeyFields defines the access key fields.
const (
	AccessKeyField_Parent               = "parent"
	AccessKeyField_AccessKeyId          = "id"
	AccessKeyField_Title                = "title"
	AccessKeyField_Description          = "description"
	AccessKeyField_UnencryptedAccessKey = "unencrypted_access_key"
	AccessKeyField_EncryptedAccessKey   = "encrypted_access_key"
	AccessKeyField_CreateTime           = "create_time"
	AccessKeyField_UpdateTime           = "update_time"
)

// AccessKey represents an API access key for a user.
// This allows third-party systems to make API requests on behalf of the user.
type AccessKey struct {
	// Parent is the parent of the access key (always a user)
	Parent AccessKeyParent
	// AccessKeyId is the unique identifier for the access key
	AccessKeyId AccessKeyId
	// Title is the title of the access key
	Title string
	// Description is the description of the access key
	Description string
	// UnencryptedAccessKey is the plain text access key (only populated on creation)
	UnencryptedAccessKey string
	// EncryptedAccessKey is the bcrypt hash of the access key (stored in database)
	EncryptedAccessKey string
	// CreateTime is the time the access key was created
	CreateTime time.Time
	// UpdateTime is the time the access key was last updated
	UpdateTime time.Time
}

type AccessKeyParent struct {
	UserId int64 `aip_pattern:"key=user"`
}

type AccessKeyId struct {
	AccessKeyId int64 `aip_pattern:"key=access_key"`
}
