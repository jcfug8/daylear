package model

import permPb "github.com/jcfug8/daylear/server/genapi/api/types"

// CircleUserFields defines the circleUser fields.
var CircleUserFields = circleUserFields{
	CircleUserId:    "circle_user_id",
	CircleId:        "circle_id",
	UserId:          "user_id",
	PermissionLevel: "permission_level",
}

type circleUserFields struct {
	CircleUserId    string
	CircleId        string
	UserId          string
	PermissionLevel string
}

// Map maps the circleUser fields to their corresponding model values.
func (fields circleUserFields) Map(m CircleUser) map[string]any {
	return map[string]any{
		fields.CircleUserId:    m.CircleUserId,
		fields.CircleId:        m.CircleId,
		fields.UserId:          m.UserId,
		fields.PermissionLevel: m.PermissionLevel,
	}
}

// Mask returns a FieldMask for the circleUser fields.
func (fields circleUserFields) Mask() []string {
	return []string{
		fields.CircleUserId,
		fields.CircleId,
		fields.UserId,
		fields.PermissionLevel,
	}
}

// CircleUser -
type CircleUser struct {
	CircleUserId    int64                  `gorm:"primaryKey;bigint;not null;<-:false"`
	CircleId        int64                  `gorm:"not null;index"`
	UserId          int64                  `gorm:"not null;index"`
	PermissionLevel permPb.PermissionLevel `gorm:"default:100"`
}

// TableName -
func (CircleUser) TableName() string {
	return "circle_user"
}
