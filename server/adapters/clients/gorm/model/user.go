package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// UserMap maps the User fields to their corresponding
// fields in the model.
var UserMap = masks.NewFieldMap().
	MapFieldToFields(model.UserFields.Id,
		UserFields.UserId).
	MapFieldToFields(model.UserFields.Email,
		UserFields.Email).
	MapFieldToFields(model.UserFields.Username,
		UserFields.Username).
	MapFieldToFields(model.UserFields.GivenName,
		UserFields.GivenName).
	MapFieldToFields(model.UserFields.FamilyName,
		UserFields.FamilyName).
	MapFieldToFields(model.UserFields.Visibility,
		UserFields.Visibility).
	MapFieldToFields(model.UserFields.ImageUri,
		UserFields.ImageUri).
	MapFieldToFields(model.UserFields.Bio,
		UserFields.Bio)

// UserFields defines the user fields.
var UserFields = userFields{
	UserId:     "daylear_user.user_id",
	Username:   "daylear_user.username",
	GivenName:  "daylear_user.given_name",
	FamilyName: "daylear_user.family_name",
	ImageUri:   "daylear_user.image_uri",
	Bio:        "daylear_user.bio",
	Visibility: "daylear_user.visibility",

	Email: "daylear_user.email",

	AmazonId:   "daylear_user.amazon_id",
	FacebookId: "daylear_user.facebook_id",
	GoogleId:   "daylear_user.google_id",

	Permission: "user_access.permission_level",
	State:      "user_access.state",
}

type userFields struct {
	UserId     string
	Username   string
	GivenName  string
	FamilyName string
	ImageUri   string
	Bio        string
	Visibility string

	Email string

	AmazonId   string
	FacebookId string
	GoogleId   string

	Permission string
	State      string
}

// Map maps the user fields to their corresponding model values.
func (fields userFields) Map(m User) map[string]any {
	return map[string]any{
		fields.UserId:     m.UserId,
		fields.Username:   m.Username,
		fields.GivenName:  m.GivenName,
		fields.FamilyName: m.FamilyName,
		fields.Visibility: m.Visibility,
		fields.ImageUri:   m.ImageUri,
		fields.Bio:        m.Bio,

		fields.Email: m.Email,

		fields.AmazonId:   m.AmazonId,
		fields.FacebookId: m.FacebookId,
		fields.GoogleId:   m.GoogleId,

		fields.Permission: m.PermissionLevel,
		fields.State:      m.State,
	}
}

// Mask returns a FieldMask for the user fields.
func (fields userFields) Mask() []string {
	return []string{
		fields.UserId,
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
		fields.ImageUri,
		fields.Bio,
		fields.Visibility,

		fields.Email,

		fields.AmazonId,
		fields.FacebookId,
		fields.GoogleId,

		fields.Permission,
		fields.State,
	}
}

// User -
type User struct {
	UserId     int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	Username   string                `gorm:"size:50;not null;uniqueIndex"`
	GivenName  string                `gorm:"size:100"`
	FamilyName string                `gorm:"size:100"`
	Visibility types.VisibilityLevel `gorm:"not null;default:1"`

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
