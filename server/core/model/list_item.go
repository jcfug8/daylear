package model

import (
	"time"
)

// ----------------------------------------------------------------------------
// ListItem Fields

// ListItemFields defines the list item fields.
const (
	ListItemField_Parent         = "parent"
	ListItemField_Id             = "id"
	ListItemField_Title          = "title"
	ListItemField_Points         = "points"
	ListItemField_RecurrenceRule = "recurrence_rule"
	ListItemField_ListSectionId  = "list_section_id"
	ListItemField_CreateTime     = "create_time"
	ListItemField_UpdateTime     = "update_time"
)

// ListItem defines the model for a list item.
type ListItem struct {
	Parent         ListItemParent
	Id             ListItemId
	Title          string
	Points         int32
	RecurrenceRule string
	ListSectionId  int64 `aip_pattern:"key=list_section"`
	CreateTime     time.Time
	UpdateTime     time.Time
}

// ListItemId defines the ID for a list item.
type ListItemId struct {
	ListItemId int64 `aip_pattern:"key=list_item"`
}

type ListItemParent struct {
	ListId ListId
}
