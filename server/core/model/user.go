package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
)

// User defines the model for a user.
type User struct {
	Id         UserId
	Parent     UserParent
	Username   string
	GivenName  string
	FamilyName string

	// the image url for the user
	ImageUri string

	// the bio for the user
	Bio string

	Email string

	AmazonId   string
	FacebookId string
	GoogleId   string

	UserAccess   UserAccess
	CircleAccess CircleAccess
}

// UserId defines the name for a user.
type UserId struct {
	UserId int64 `aip_pattern:"key=user,public_user"`
}

type UserParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	UserId   int64 `aip_pattern:"key=user"`
}

// ----------------------------------------------------------------------------
// Fields

// UserFields defines the user fields.
var UserFields = userFields{
	Id:         "id",
	Username:   "username",
	GivenName:  "given_name",
	FamilyName: "family_name",

	ImageUri: "image_uri",
	Bio:      "bio",

	Email: "email",

	GoogleId:   "google_id",
	FacebookId: "facebook_id",
	AmazonId:   "amazon_id",

	AccessName:            "access_id",
	AccessPermissionLevel: "permission_level",
	AccessState:           "state",
}

type userFields struct {
	Id         string
	Username   string
	GivenName  string
	FamilyName string

	ImageUri string
	Bio      string

	Email string

	GoogleId   string
	FacebookId string
	AmazonId   string

	AccessName            string
	AccessPermissionLevel string
	AccessState           string
}

// Mask returns a FieldMask for the user fields.
func (fields userFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
		fields.ImageUri,
		fields.Bio,

		fields.Email,

		fields.GoogleId,
		fields.FacebookId,
		fields.AmazonId,

		fields.AccessName,
		fields.AccessPermissionLevel,
		fields.AccessState,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields userFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Username,
		fields.GivenName,
		fields.FamilyName,
		fields.ImageUri,
		fields.Bio,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
