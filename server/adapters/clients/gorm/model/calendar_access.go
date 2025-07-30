package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// CalendarAccessFields defines the calendarAccess fields.
var CalendarAccessFields = calendarAccessFields{
	CalendarAccessId:  "calendar_access.calendar_access_id",
	CalendarId:        "calendar_access.calendar_id",
	RequesterUserId:   "calendar_access.requester_user_id",
	RequesterCircleId: "calendar_access.requester_circle_id",
	RecipientUserId:   "calendar_access.recipient_user_id",
	RecipientCircleId: "calendar_access.recipient_circle_id",
	PermissionLevel:   "calendar_access.permission_level",
	State:             "calendar_access.state",
}

type calendarAccessFields struct {
	CalendarAccessId  string
	CalendarId        string
	RequesterUserId   string
	RequesterCircleId string
	RecipientUserId   string
	RecipientCircleId string
	PermissionLevel   string
	State             string
}

// Map maps the calendarAccess fields to their corresponding model values.
func (fields calendarAccessFields) Map(m CalendarAccess) map[string]any {
	return map[string]any{
		fields.CalendarAccessId:  m.CalendarAccessId,
		fields.CalendarId:        m.CalendarId,
		fields.RequesterUserId:   m.RequesterUserId,
		fields.RequesterCircleId: m.RequesterCircleId,
		fields.RecipientUserId:   m.RecipientUserId,
		fields.RecipientCircleId: m.RecipientCircleId,
		fields.PermissionLevel:   m.PermissionLevel,
		fields.State:             m.State,
	}
}

// Mask returns a FieldMask for the calendarAccess fields.
func (fields calendarAccessFields) Mask() []string {
	return []string{
		fields.CalendarAccessId,
		fields.CalendarId,
		fields.RequesterUserId,
		fields.RequesterCircleId,
		fields.RecipientUserId,
		fields.RecipientCircleId,
		fields.PermissionLevel,
		fields.State,
	}
}

// CalendarAccess represents access control for a calendar
type CalendarAccess struct {
	CalendarAccessId  int64                 `gorm:"primaryKey;bigint;not null;<-:false"`
	CalendarId        int64                 `gorm:"not null;index;uniqueIndex:idx_calendar_id_recipient_user_id,idx_calendar_id_recipient_circle_id"`
	RequesterUserId   int64                 `gorm:"index"`
	RequesterCircleId int64                 `gorm:"index"`
	RecipientUserId   int64                 `gorm:"not null;uniqueIndex:idx_calendar_id_recipient_user_id,where:recipient_user_id <> 0"`
	RecipientCircleId int64                 `gorm:"not null;uniqueIndex:idx_calendar_id_recipient_circle_id,where:recipient_circle_id <> 0"`
	PermissionLevel   types.PermissionLevel `gorm:"not null"`
	State             types.AccessState     `gorm:"not null"`

	// Read-only fields from joins
	RecipientUsername     string `gorm:"->;-:migration"` // read only from join
	RecipientGivenName    string `gorm:"->;-:migration"` // read only from join
	RecipientFamilyName   string `gorm:"->;-:migration"` // read only from join
	RecipientCircleTitle  string `gorm:"->;-:migration"` // read only from join
	RecipientCircleHandle string `gorm:"->;-:migration"` // read only from join
}

// TableName sets the table name for the CalendarAccess model.
func (CalendarAccess) TableName() string {
	return "calendar_access"
}
