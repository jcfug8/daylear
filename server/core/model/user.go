package model

import "strings"

var _ ResourceId = UserId{}

const (
	UserField_Parent     = "parent"
	UserField_Id         = "id"
	UserField_Username   = "username"
	UserField_GivenName  = "given_name"
	UserField_FamilyName = "family_name"
	UserField_ImageUri   = "image_uri"
	UserField_Bio        = "bio"
	UserField_Email      = "email"
	UserField_GoogleId   = "google_id"
	UserField_FacebookId = "facebook_id"
	UserField_AmazonId   = "amazon_id"

	UserField_AccessName            = "access_id"
	UserField_AccessPermissionLevel = "permission_level"
	UserField_AccessState           = "state"
	UserField_Favorited             = "favorited"
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

	Favorited bool

	UserAccess   UserAccess
	CircleAccess CircleAccess
}

// UserId defines the name for a user.
type UserId struct {
	UserId int64 `aip_pattern:"key=user,public_user"`
}

// isResourceId - implements the ResourceId interface.
func (u UserId) isResourceId() {}

type UserParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	UserId   int64 `aip_pattern:"key=user"`
}

func (u User) GetFullName() string {
	fullName := strings.TrimSpace(u.GivenName + " " + u.FamilyName)
	if fullName == "" {
		return u.Username
	}
	return fullName
}
