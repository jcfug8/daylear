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
		UserFields.Visibility)

// UserFields defines the user fields.
var UserFields = userFields{
	UserId:     "user_id",
	Username:   "username",
	GivenName:  "given_name",
	FamilyName: "family_name",
	Visibility: "visibility",

	Email: "email",

	AmazonId:   "amazon_id",
	FacebookId: "facebook_id",
	GoogleId:   "google_id",
}

type userFields struct {
	UserId     string
	Username   string
	GivenName  string
	FamilyName string
	Visibility string

	Email string

	AmazonId   string
	FacebookId string
	GoogleId   string
}

// Map maps the user fields to their corresponding model values.
func (fields userFields) Map(m User) map[string]any {
	return map[string]any{
		fields.UserId:     m.UserId,
		fields.Username:   m.Username,
		fields.GivenName:  m.GivenName,
		fields.FamilyName: m.FamilyName,
		fields.Visibility: m.Visibility,

		fields.Email: m.Email,

		fields.AmazonId:   m.AmazonId,
		fields.FacebookId: m.FacebookId,
		fields.GoogleId:   m.GoogleId,
	}
}

// Mask returns a FieldMask for the user fields.
func (fields userFields) Mask() []string {
	return []string{
		fields.UserId,
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
		fields.Visibility,

		fields.Email,

		fields.AmazonId,
		fields.FacebookId,
		fields.GoogleId,
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
}

// TableName -
func (User) TableName() string {
	return "daylear_user"
}
