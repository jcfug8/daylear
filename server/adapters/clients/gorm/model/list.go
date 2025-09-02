package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	ListTable = "lists"
)

const (
	ListFields_ListId          = "list_id"
	ListFields_Title           = "title"
	ListFields_Description     = "description"
	ListFields_ShowCompleted   = "show_completed"
	ListFields_VisibilityLevel = "visibility_level"
	ListFields_Sections        = "sections"
	ListFields_CreateTime      = "create_time"
	ListFields_UpdateTime      = "update_time"
)

var ListFieldMasker = fieldmask.NewSQLFieldMasker(List{}, map[string][]fieldmask.Field{
	model.ListField_Id:              {{Name: ListFields_ListId}},
	model.ListField_Title:           {{Name: ListFields_Title}},
	model.ListField_Description:     {{Name: ListFields_Description}},
	model.ListField_ShowCompleted:   {{Name: ListFields_ShowCompleted}},
	model.ListField_VisibilityLevel: {{Name: ListFields_VisibilityLevel}},
	model.ListField_Sections:        {{Name: ListFields_Sections}},
	model.ListField_CreateTime:      {{Name: ListFields_CreateTime}},
	model.ListField_UpdateTime:      {{Name: ListFields_UpdateTime}},
	model.ListField_Favorited:       {{Name: ListFavoriteFields_ListFavoriteId}},
	model.ListField_ListAccess: {
		{Name: ListAccessFields_ListAccessId},
		{Name: ListAccessFields_PermissionLevel},
		{Name: ListAccessFields_State},
		{Name: ListAccessFields_AcceptTarget},
	},
})

var ListSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	model.ListField_Title:           {Name: ListFields_Title, Table: ListTable},
	model.ListField_Description:     {Name: ListFields_Description, Table: ListTable},
	model.ListField_ShowCompleted:   {Name: ListFields_ShowCompleted, Table: ListTable},
	model.ListField_VisibilityLevel: {Name: ListFields_VisibilityLevel, Table: ListTable},
	model.ListField_CreateTime:      {Name: ListFields_CreateTime, Table: ListTable},
	model.ListField_UpdateTime:      {Name: ListFields_UpdateTime, Table: ListTable},
}, true)

// List -
type List struct {
	ListId          int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title           string
	Description     string
	ShowCompleted   bool                  `gorm:"not null;default:false"`
	VisibilityLevel types.VisibilityLevel `gorm:"not null;default:1"`
	Sections        []byte                `gorm:"type:jsonb"`
	CreateTime      time.Time             `gorm:"column:create_time;autoCreateTime"`
	UpdateTime      time.Time             `gorm:"column:update_time;autoUpdateTime"`

	// ListAccess data
	ListAccessId    int64                 `gorm:"->;-:migration"` // only used for read from a join
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
	AcceptTarget    types.AcceptTarget    `gorm:"->;-:migration"` // only used for read from a join

	// ListFavorite data
	ListFavoriteId int64 `gorm:"->;-:migration"` // only used for read from a join
}

// TableName -
func (List) TableName() string {
	return ListTable
}
