package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	UserColumn_UserId     = "daylear_user.user_id"
	UserColumn_Username   = "daylear_user.username"
	UserColumn_GivenName  = "daylear_user.given_name"
	UserColumn_FamilyName = "daylear_user.family_name"
	UserColumn_ImageUri   = "daylear_user.image_uri"
	UserColumn_Bio        = "daylear_user.bio"
	UserColumn_Email      = "daylear_user.email"
	UserColumn_AmazonId   = "daylear_user.amazon_id"
	UserColumn_FacebookId = "daylear_user.facebook_id"
	UserColumn_GoogleId   = "daylear_user.google_id"
)

var UserFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.UserField_Id:         {UserColumn_UserId},
	cmodel.UserField_Username:   {UserColumn_Username},
	cmodel.UserField_GivenName:  {UserColumn_GivenName},
	cmodel.UserField_FamilyName: {UserColumn_FamilyName},
	cmodel.UserField_ImageUri:   {UserColumn_ImageUri},
	cmodel.UserField_Bio:        {UserColumn_Bio},
	cmodel.UserField_Email:      {UserColumn_Email},

	cmodel.UserField_AccessName:            {UserAccessColumn_UserAccessId},
	cmodel.UserField_AccessPermissionLevel: {UserAccessColumn_PermissionLevel},
	cmodel.UserField_AccessState:           {UserAccessColumn_State},
})
var UserCircleFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.UserField_Id:         {UserColumn_UserId},
	cmodel.UserField_Username:   {UserColumn_Username},
	cmodel.UserField_GivenName:  {UserColumn_GivenName},
	cmodel.UserField_FamilyName: {UserColumn_FamilyName},
	cmodel.UserField_ImageUri:   {UserColumn_ImageUri},
	cmodel.UserField_Bio:        {UserColumn_Bio},
	cmodel.UserField_Email:      {UserColumn_Email},

	cmodel.UserField_AccessName:            {CircleAccessColumn_CircleAccessId},
	cmodel.UserField_AccessPermissionLevel: {CircleAccessColumn_PermissionLevel},
	cmodel.UserField_AccessState:           {CircleAccessColumn_State},
})

var UpdateUserFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.UserField_Username:   {UserColumn_Username},
	cmodel.UserField_GivenName:  {UserColumn_GivenName},
	cmodel.UserField_FamilyName: {UserColumn_FamilyName},
	cmodel.UserField_ImageUri:   {UserColumn_ImageUri},
	cmodel.UserField_Bio:        {UserColumn_Bio},
	cmodel.UserField_Email:      {UserColumn_Email},
})

var UserSQLConverter = filter.NewSQLConverter(map[string]string{
	"username":         UserColumn_Username,
	"permission_level": UserAccessColumn_PermissionLevel,
	"state":            UserAccessColumn_State,
	"google_id":        UserColumn_GoogleId,
	"facebook_id":      UserColumn_FacebookId,
	"amazon_id":        UserColumn_AmazonId,
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
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
	RequesterUserId int64                 `gorm:"->;-:migration"` // only used for read from a join
}

// TableName -
func (User) TableName() string {
	return "daylear_user"
}
