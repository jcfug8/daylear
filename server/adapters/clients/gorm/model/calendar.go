package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CalendarTable = "calendar"
)

const (
	CalendarColumn_CalendarId      = "calendar_id"
	CalendarColumn_Title           = "title"
	CalendarColumn_Description     = "description"
	CalendarColumn_VisibilityLevel = "visibility_level"
	CalendarColumn_CreateTime      = "create_time"
	CalendarColumn_UpdateTime      = "update_time"
)

var CalendarFieldMasker = fieldmask.NewSQLFieldMasker(Calendar{}, map[string][]fieldmask.Field{
	cmodel.CalendarField_Parent:      {{Name: CalendarColumn_CalendarId, Table: CalendarTable}},
	cmodel.CalendarField_CalendarId:  {{Name: CalendarColumn_CalendarId, Table: CalendarTable}},
	cmodel.CalendarField_Title:       {{Name: CalendarColumn_Title, Table: CalendarTable, Updatable: true}},
	cmodel.CalendarField_Description: {{Name: CalendarColumn_Description, Table: CalendarTable, Updatable: true}},
	cmodel.CalendarField_Visibility:  {{Name: CalendarColumn_VisibilityLevel, Table: CalendarTable, Updatable: true}},
	cmodel.CalendarField_CreateTime:  {{Name: CalendarColumn_CreateTime, Table: CalendarTable}},
	cmodel.CalendarField_UpdateTime:  {{Name: CalendarColumn_UpdateTime, Table: CalendarTable}},

	cmodel.CalendarField_Favorited: {{Name: CalendarFavoriteFields_CalendarFavoriteId, Table: CalendarFavoriteTable}},

	cmodel.CalendarField_CalendarAccess: {
		{Name: CalendarAccessColumn_CalendarAccessId, Table: CalendarAccessTable},
		{Name: CalendarAccessColumn_PermissionLevel, Table: CalendarAccessTable},
		{Name: CalendarAccessColumn_State, Table: CalendarAccessTable},
		{Name: CalendarAccessColumn_AcceptTarget, Table: CalendarAccessTable},
	},
})
var CalendarSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"visibility":       {Name: CalendarColumn_VisibilityLevel, Table: CalendarTable},
	"permission_level": {Name: CalendarAccessColumn_PermissionLevel, Table: CalendarAccessTable},
	"state":            {Name: CalendarAccessColumn_State, Table: CalendarAccessTable},
	"favorited":        {Name: CalendarFavoriteFields_CalendarFavoriteId, Table: CalendarFavoriteTable, CustomConverter: favoritedSQLFilterConverter},
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
	AcceptTarget     types.AcceptTarget    `gorm:"->;-:migration"`

	// CalendarFavorite data (only used for read from a join)
	CalendarFavoriteId int64 `gorm:"->;-:migration"`
}

// TableName sets the table name for the Calendar model.
func (Calendar) TableName() string {
	return CalendarTable
}
