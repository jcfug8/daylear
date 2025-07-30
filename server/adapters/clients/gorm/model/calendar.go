package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/masks"
	coremodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// CalendarMap maps the Calendar fields to their corresponding fields in the core model.
var CalendarMap = masks.NewFieldMap().
	MapFieldToFields(coremodel.CalendarFields.CalendarId,
		CalendarFields.CalendarId).
	MapFieldToFields(coremodel.CalendarFields.Title,
		CalendarFields.Title).
	MapFieldToFields(coremodel.CalendarFields.Description,
		CalendarFields.Description).
	MapFieldToFields(coremodel.CalendarFields.Visibility,
		CalendarFields.Visibility)

// CalendarFields defines the calendar fields in the GORM model.
var CalendarFields = calendarFields{
	CalendarId:  "calendar.calendar_id",
	Title:       "calendar.title",
	Description: "calendar.description",
	Visibility:  "calendar.visibility_level",
	Permission:  "calendar_access.permission_level",
	State:       "calendar_access.state",
}

type calendarFields struct {
	CalendarId  string
	Title       string
	Description string
	Visibility  string
	Permission  string
	State       string
}

// Map maps the calendar fields to their corresponding model values.
func (fields calendarFields) Map(m Calendar) map[string]any {
	return map[string]any{
		fields.CalendarId:  m.CalendarId,
		fields.Title:       m.Title,
		fields.Description: m.Description,
		fields.Visibility:  m.VisibilityLevel,
		fields.Permission:  m.PermissionLevel,
		fields.State:       m.State,
	}
}

// Mask returns a FieldMask for the calendar fields.
func (fields calendarFields) Mask() []string {
	return []string{
		fields.CalendarId,
		fields.Title,
		fields.Description,
		fields.Visibility,
		fields.Permission,
		fields.State,
	}
}

// Calendar is the GORM model for a calendar.
type Calendar struct {
	CalendarId      int64                 `gorm:"primaryKey;column:calendar_id;autoIncrement;<-:false"`
	Title           string                `gorm:"column:title;not null"`
	Description     string                `gorm:"column:description"`
	VisibilityLevel types.VisibilityLevel `gorm:"column:visibility_level;not null;default:1"`
	CreateTime      time.Time             `gorm:"column:create_time;autoCreateTime"`
	UpdateTime      time.Time             `gorm:"column:update_time;autoUpdateTime"`

	// CalendarAccess data (only used for read from a join)
	CalendarAccessId int64                 `gorm:"->;-:migration"`
	PermissionLevel  types.PermissionLevel `gorm:"->;-:migration"`
	State            types.AccessState     `gorm:"->;-:migration"`
}

// TableName sets the table name for the Calendar model.
func (Calendar) TableName() string {
	return "calendar"
}
