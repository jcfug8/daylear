package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
)

// CircleAccessFields defines the circleAccess fields.
var CircleAccessFields = circleAccessFields{
	CircleAccessId:    "circle_access.circle_access_id",
	CircleId:          "circle_access.circle_id",
	RequesterUserId:   "circle_access.requester_user_id",
	RequesterCircleId: "circle_access.requester_circle_id",
	RecipientUserId:   "circle_access.recipient_user_id",
	PermissionLevel:   "circle_access.permission_level",
	State:             "circle_access.state",
	Title:             "circle_access.title",
}

type circleAccessFields struct {
	CircleAccessId    string
	CircleId          string
	RequesterUserId   string
	RequesterCircleId string
	RecipientUserId   string
	PermissionLevel   string
	State             string
	Title             string
}

// Map maps the circleAccess fields to their corresponding model values.
func (fields circleAccessFields) Map(m CircleAccess) map[string]any {
	return map[string]any{
		fields.CircleAccessId:    m.CircleAccessId,
		fields.CircleId:          m.CircleId,
		fields.RequesterUserId:   m.RequesterUserId,
		fields.RequesterCircleId: m.RequesterCircleId,
		fields.RecipientUserId:   m.RecipientUserId,
		fields.PermissionLevel:   m.PermissionLevel,
		fields.State:             m.State,
		fields.Title:             m.Title,
	}
}

// Mask returns a FieldMask for the circleAccess fields.
func (fields circleAccessFields) Mask() []string {
	return []string{
		fields.CircleAccessId,
		fields.CircleId,
		fields.RequesterUserId,
		fields.RequesterCircleId,
		fields.RecipientUserId,
		fields.PermissionLevel,
		fields.State,
		fields.Title,
	}
}

// CircleAccess -
type CircleAccess struct {
	CircleAccessId    int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	CircleId          int64                  `gorm:"not null;index"`
	RequesterUserId   int64                  `gorm:"index"`
	RequesterCircleId int64                  `gorm:"index"`
	RecipientUserId   int64                  `gorm:"not null;index"`
	PermissionLevel   permPb.PermissionLevel `gorm:"not null"`
	State             pb.Access_State        `gorm:"not null"`
	Title             string                 `gorm:"->"` // read only from join
}

// TableName -
func (CircleAccess) TableName() string {
	return "circle_access"
}
