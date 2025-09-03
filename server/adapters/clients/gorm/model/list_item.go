package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
)

const (
	ListItemTable = "list_item"
)

const (
	ListItemFields_ListItemId     = "list_item_id"
	ListItemFields_ListId         = "list_id"
	ListItemFields_Title          = "title"
	ListItemFields_Points         = "points"
	ListItemFields_RecurrenceRule = "recurrence_rule"
	ListItemFields_ListSectionId  = "list_section_id"
	ListItemFields_CreateTime     = "create_time"
	ListItemFields_UpdateTime     = "update_time"
)

var ListItemFieldMasker = fieldmask.NewSQLFieldMasker(ListItem{}, map[string][]fieldmask.Field{
	model.ListItemField_Id:             {{Name: ListItemFields_ListItemId, Table: ListItemTable}},
	model.ListItemField_Title:          {{Name: ListItemFields_Title, Table: ListItemTable}},
	model.ListItemField_Points:         {{Name: ListItemFields_Points, Table: ListItemTable}},
	model.ListItemField_RecurrenceRule: {{Name: ListItemFields_RecurrenceRule, Table: ListItemTable}},
	model.ListItemField_ListSectionId:  {{Name: ListItemFields_ListSectionId, Table: ListItemTable}},
	model.ListItemField_CreateTime:     {{Name: ListItemFields_CreateTime, Table: ListItemTable}},
	model.ListItemField_UpdateTime:     {{Name: ListItemFields_UpdateTime, Table: ListItemTable}},
	model.ListItemField_Parent:         {{Name: ListItemFields_ListId, Table: ListItemTable}},
})

var ListItemSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"points": {Name: ListItemFields_Points, Table: ListItemTable},
	"title":  {Name: ListItemFields_Title, Table: ListItemTable},
}, true)

// ListItem represents a list item in the database
type ListItem struct {
	ListItemId     int64  `gorm:"primaryKey;bigint;not null;<-:false"`
	ListId         int64  `gorm:"bigint;not null;index"`
	Title          string `gorm:"not null"`
	Points         int32  `gorm:"not null;default:0"`
	RecurrenceRule string
	ListSectionId  int64     `gorm:"bigint"`
	CreateTime     time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime     time.Time `gorm:"column:update_time;autoUpdateTime"`
}

// TableName returns the table name for the ListItem model
func (ListItem) TableName() string {
	return ListItemTable
}
