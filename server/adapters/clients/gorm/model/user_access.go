package model

import (
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

// UserAccessFields defines the userAccess fields.
var UserAccessFields = userAccessFields{
	UserAccessId:    "user_access.user_access_id",
	UserId:          "user_access.user_id",
	RequesterUserId: "user_access.requester_user_id",
	RecipientUserId: "user_access.recipient_user_id",
	PermissionLevel: "user_access.permission_level",
	State:           "user_access.state",
	Title:           "user_access.title",
}

type userAccessFields struct {
	UserAccessId    string
	UserId          string
	RequesterUserId string
	RecipientUserId string
	PermissionLevel string
	State           string
	Title           string
}

// Map maps the userAccess fields to their corresponding model values.
func (fields userAccessFields) Map(m UserAccess) map[string]any {
	return map[string]any{
		fields.UserAccessId:    m.UserAccessId,
		fields.UserId:          m.UserId,
		fields.RequesterUserId: m.RequesterUserId,
		fields.RequesterUserId: m.RequesterUserId,
		fields.RecipientUserId: m.RecipientUserId,
		fields.PermissionLevel: m.PermissionLevel,
		fields.State:           m.State,
		fields.Title:           m.Title,
	}
}

// Mask returns a FieldMask for the userAccess fields.
func (fields userAccessFields) Mask() []string {
	return []string{
		fields.UserAccessId,
		fields.UserId,
		fields.RequesterUserId,
		fields.RequesterUserId,
		fields.RecipientUserId,
		fields.PermissionLevel,
		fields.State,
		fields.Title,
	}
}

// UserAccess -
type UserAccess struct {
	UserAccessId    int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	UserId          int64                  `gorm:"not null;index"`
	RequesterUserId int64                  `gorm:"index"`
	RecipientUserId int64                  `gorm:"not null;index"`
	PermissionLevel permPb.PermissionLevel `gorm:"not null"`
	State           pb.Access_State        `gorm:"not null"`
	Title           string                 `gorm:"->"` // read only from join
}

// TableName -
func (UserAccess) TableName() string {
	return "user_access"
}
