package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/core/model"
)

// UserMap maps the User fields to their corresponding
// fields in the model.
var UserMap = masks.NewFieldMap().
	MapFieldToFields(model.UserFields.Id,
		UserFields.UserId).
	MapFieldToFields(model.UserFields.Email,
		UserFields.Email).
	MapFieldToFields(model.UserFields.Username,
		UserFields.Username)

// UserFields defines the user fields.
var UserFields = userFields{
	UserId:   "user_id",
	Email:    "email",
	Username: "username",

	AmazonId:   "amazon_id",
	FacebookId: "facebook_id",
	GoogleId:   "google_id",
}

type userFields struct {
	UserId   string
	Email    string
	Username string

	AmazonId   string
	FacebookId string
	GoogleId   string
}

// Map maps the user fields to their corresponding model values.
func (fields userFields) Map(m User) map[string]any {
	return map[string]any{
		fields.UserId:   m.UserId,
		fields.Email:    m.Email,
		fields.Username: m.Username,

		fields.AmazonId:   m.AmazonId,
		fields.FacebookId: m.FacebookId,
		fields.GoogleId:   m.GoogleId,
	}
}

// Mask returns a FieldMask for the user fields.
func (fields userFields) Mask() []string {
	return []string{
		fields.UserId,
		fields.Email,
		fields.Username,

		fields.AmazonId,
		fields.FacebookId,
		fields.GoogleId,
	}
}
