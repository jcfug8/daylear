package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CircleTable = "circle"
)

const (
	CircleColumn_CircleId    = "circle_id"
	CircleColumn_Title       = "title"
	CircleColumn_Description = "description"
	CircleColumn_Handle      = "handle"
	CircleColumn_ImageURI    = "image_uri"
	CircleColumn_Visibility  = "visibility_level"
)

var CircleFieldMasker = fieldmask.NewSQLFieldMasker(Circle{}, map[string][]fieldmask.Field{
	cmodel.CircleField_CircleId:    {{Name: CircleColumn_CircleId, Table: CircleTable}},
	cmodel.CircleField_Title:       {{Name: CircleColumn_Title, Table: CircleTable, Updatable: true}},
	cmodel.CircleField_Description: {{Name: CircleColumn_Description, Table: CircleTable, Updatable: true}},
	cmodel.CircleField_Handle:      {{Name: CircleColumn_Handle, Table: CircleTable, Updatable: true}},
	cmodel.CircleField_ImageURI:    {{Name: CircleColumn_ImageURI, Table: CircleTable, Updatable: true}},
	cmodel.CircleField_Visibility:  {{Name: CircleColumn_Visibility, Table: CircleTable, Updatable: true}},

	cmodel.CircleField_Favorited: {{Name: CircleFavoriteFields_CircleFavoriteId, Table: CircleFavoriteTable}},

	cmodel.CircleField_CircleAccess: {
		{Name: CircleAccessColumn_CircleAccessId, Table: CircleAccessTable},
		{Name: CircleAccessColumn_PermissionLevel, Table: CircleAccessTable},
		{Name: CircleAccessColumn_State, Table: CircleAccessTable},
		{Name: CircleAccessColumn_AcceptTarget, Table: CircleAccessTable},
	},
})

var CircleSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"visibility":       {Name: CircleColumn_Visibility, Table: CircleTable},
	"permission_level": {Name: CircleAccessColumn_PermissionLevel, Table: CircleAccessTable},
	"state":            {Name: CircleAccessColumn_State, Table: CircleAccessTable},
	"favorited":        {Name: CircleFavoriteFields_CircleFavoriteId, Table: CircleFavoriteTable, CustomConverter: favoritedSQLFilterConverter},
}, true)

// Circle is the GORM model for a circle.
type Circle struct {
	CircleId        int64                 `gorm:"primaryKey;column:circle_id;autoIncrement;<-:false"`
	Title           string                `gorm:"column:title;not null"`
	Description     string                `gorm:"column:description"`
	Handle          string                `gorm:"column:handle;unique;not null"`
	ImageURI        string                `gorm:"column:image_uri"`
	VisibilityLevel types.VisibilityLevel `gorm:"column:visibility_level;not null;default:1"`

	// CircleAccess data
	CircleAccessId  int64                 `gorm:"->;-:migration"` // only used for read from a join
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
	AcceptTarget    types.AcceptTarget    `gorm:"->;-:migration"` // only used for read from a join

	// CircleFavorite data
	CircleFavoriteId int64 `gorm:"->;-:migration"` // only used for read from a join
}

// TableName sets the table name for the Circle model.
func (Circle) TableName() string {
	return CircleTable
}
