package model

import (
	"time"

	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ ResourceId = ListId{}

// ----------------------------------------------------------------------------
// Fields

// ListFields defines the list fields.
const (
	ListField_Parent          = "parent"
	ListField_Id              = "id"
	ListField_Title           = "title"
	ListField_Description     = "description"
	ListField_ShowCompleted   = "show_completed"
	ListField_VisibilityLevel = "visibility_level"
	ListField_Sections        = "sections"
	ListField_CreateTime      = "create_time"
	ListField_UpdateTime      = "update_time"
	ListField_Favorited       = "favorited"

	ListField_ListAccess = "list_access"
)

// List defines the model for a list.
type List struct {
	Id              ListId
	Parent          ListParent
	Title           string
	Description     string
	ShowCompleted   bool
	VisibilityLevel types.VisibilityLevel
	Sections        []ListSection
	CreateTime      time.Time
	UpdateTime      time.Time
	Favorited       bool

	// The access details for the current user/circle
	ListAccess ListAccess
}

type ListSection struct {
	ListId int64 `aip_pattern:"key=list"`
	Id     int64 `aip_pattern:"key=list_section"`
	Title  string
}

// ListId defines the name for a list.
type ListId struct {
	ListId int64 `aip_pattern:"key=list"`
}

// isResourceId - implements the ResourceId interface.
func (l ListId) isResourceId() {
}

// ListParent defines the name for a list parent.
// Supports both circle and user parents for lists.
type ListParent struct {
	CircleId int64 `aip_pattern:"key=circle"`
	UserId   int64 `aip_pattern:"key=user"`
}
