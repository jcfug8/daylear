package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
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
}

type circleAccessFields struct {
	CircleAccessId    string
	CircleId          string
	RequesterUserId   string
	RequesterCircleId string
	RecipientUserId   string
	PermissionLevel   string
	State             string
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
	}
}

// CircleAccess -
type CircleAccess struct {
	CircleAccessId    int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	CircleId          int64                 `gorm:"not null;index;uniqueIndex:idx_circle_id_recipient_user_id"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"not null;uniqueIndex:idx_circle_id_recipient_user_id,where:recipient_user_id <> 0"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`

	RecipientUsername   string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName  string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName string `gorm:"->;-:migration"` // read only from join
}

// TableName -
func (CircleAccess) TableName() string {
	return "circle_access"
}
