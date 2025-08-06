package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CalendarColumn_CalendarId      = "calendar_id"
	CalendarColumn_Title           = "title"
	CalendarColumn_Description     = "description"
	CalendarColumn_VisibilityLevel = "visibility_level"
	CalendarColumn_CreateTime      = "create_time"
	CalendarColumn_UpdateTime      = "update_time"
)

var CalendarFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CalendarField_Parent:      {CalendarColumn_CalendarId},
	cmodel.CalendarField_CalendarId:  {CalendarColumn_CalendarId},
	cmodel.CalendarField_Title:       {CalendarColumn_Title},
	cmodel.CalendarField_Description: {CalendarColumn_Description},
	cmodel.CalendarField_Visibility:  {CalendarColumn_VisibilityLevel},

	cmodel.CalendarField_CalendarAccess: {
		CalendarAccessColumn_CalendarAccessId,
		CalendarAccessColumn_PermissionLevel,
		CalendarAccessColumn_State,
	},
})
var UpdateCalendarFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CalendarField_Title:       {CalendarColumn_Title},
	cmodel.CalendarField_Description: {CalendarColumn_Description},
	cmodel.CalendarField_Visibility:  {CalendarColumn_VisibilityLevel},
})
var CalendarSQLConverter = filter.NewSQLConverter(map[string]string{
	"visibility":       CalendarColumn_VisibilityLevel,
	"permission_level": CalendarAccessColumn_PermissionLevel,
	"state":            CalendarAccessColumn_State,
}, true)

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
