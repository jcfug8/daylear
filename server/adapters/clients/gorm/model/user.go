package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	UserTable = "daylear_user"
)

const (
	UserColumn_UserId     = "user_id"
	UserColumn_Username   = "username"
	UserColumn_GivenName  = "given_name"
	UserColumn_FamilyName = "family_name"
	UserColumn_ImageUri   = "image_uri"
	UserColumn_Bio        = "bio"
	UserColumn_Email      = "email"
	UserColumn_AmazonId   = "amazon_id"
	UserColumn_FacebookId = "facebook_id"
	UserColumn_GoogleId   = "google_id"
)

var UserFieldMasker = fieldmask.NewSQLFieldMasker(User{}, map[string][]fieldmask.Field{
	cmodel.UserField_Id:         {{Name: UserColumn_UserId, Table: "daylear_user"}},
	cmodel.UserField_Username:   {{Name: UserColumn_Username, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_GivenName:  {{Name: UserColumn_GivenName, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_FamilyName: {{Name: UserColumn_FamilyName, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_ImageUri:   {{Name: UserColumn_ImageUri, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_Bio:        {{Name: UserColumn_Bio, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_Email:      {{Name: UserColumn_Email, Table: "daylear_user", Updatable: true}},
	cmodel.UserField_AmazonId:   {{Name: UserColumn_AmazonId, Table: "daylear_user"}},
	cmodel.UserField_FacebookId: {{Name: UserColumn_FacebookId, Table: "daylear_user"}},
	cmodel.UserField_GoogleId:   {{Name: UserColumn_GoogleId, Table: "daylear_user"}},

	cmodel.UserField_AccessName: {
		{Name: CircleAccessColumn_CircleAccessId, Table: "circle_access"},
		{Name: UserAccessColumn_UserAccessId, Table: "user_access"},
	},
	cmodel.UserField_AccessPermissionLevel: {
		{Name: CircleAccessColumn_PermissionLevel, Table: "circle_access"},
		{Name: UserAccessColumn_PermissionLevel, Table: "user_access"},
	},
	cmodel.UserField_AccessState: {
		{Name: CircleAccessColumn_State, Table: "circle_access"},
		{Name: UserAccessColumn_State, Table: "user_access"},
	},
})

var UserSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"username":         {Name: UserColumn_Username, Table: "daylear_user"},
	"permission_level": {Name: UserAccessColumn_PermissionLevel, Table: "user_access"},
	"state":            {Name: UserAccessColumn_State, Table: "user_access"},
	"google_id":        {Name: UserColumn_GoogleId, Table: "daylear_user"},
	"facebook_id":      {Name: UserColumn_FacebookId, Table: "daylear_user"},
	"amazon_id":        {Name: UserColumn_AmazonId, Table: "daylear_user"},
}, true)

var UserCircleSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"username":         {Name: UserColumn_Username, Table: "daylear_user"},
	"permission_level": {Name: CircleAccessColumn_PermissionLevel, Table: "circle_access"},
	"state":            {Name: CircleAccessColumn_State, Table: "circle_access"},
	"google_id":        {Name: UserColumn_GoogleId, Table: "daylear_user"},
	"facebook_id":      {Name: UserColumn_FacebookId, Table: "daylear_user"},
	"amazon_id":        {Name: UserColumn_AmazonId, Table: "daylear_user"},
}, true)

// User -
type User struct {
	UserId     int64  `gorm:"primaryKey;bigint;not null;<-:false"`
	Username   string `gorm:"size:50;not null;uniqueIndex"`
	GivenName  string `gorm:"size:100"`
	FamilyName string `gorm:"size:100"`

	Email string `gorm:"size:255;not null;uniqueIndex"`

	AmazonId   *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`
	FacebookId *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`
	GoogleId   *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`

	ImageUri string `gorm:"size:255"`
	Bio      string `gorm:"size:1024"`

	UserAccessId    int64                 `gorm:"->;-:migration"` // only used for read from a join
	CircleAccessId  int64                 `gorm:"->;-:migration"` // only used for read from a join
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
	RequesterUserId int64                 `gorm:"->;-:migration"` // only used for read from a join
}

// TableName -
func (User) TableName() string {
	return UserTable
}
