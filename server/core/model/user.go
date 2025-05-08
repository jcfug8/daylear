package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
)

// User defines the model for a user.
type User struct {
	Id         UserId
	Email      string
	Username   string
	GivenName  string
	FamilyName string

	AmazonId   string
	FacebookId string
	GoogleId   string
}

// UserId defines the name for a user.
type UserId struct {
	UserId int64
}

// ----------------------------------------------------------------------------
// Fields

// UserFields defines the user fields.
var UserFields = userFields{
	Id:         "id",
	Email:      "email",
	Username:   "username",
	GivenName:  "given_name",
	FamilyName: "family_name",
	GoogleId:   "google_id",
	FacebookId: "facebook_id",
	AmazonId:   "amazon_id",
}

type userFields struct {
	Id         string
	Email      string
	Username   string
	GivenName  string
	FamilyName string
	GoogleId   string
	FacebookId string
	AmazonId   string
}

// Mask returns a FieldMask for the user fields.
func (fields userFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Email,
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
		fields.GoogleId,
		fields.FacebookId,
		fields.AmazonId,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields userFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
