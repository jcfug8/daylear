package model

import (
	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	CircleColumn_CircleId    = "circle.circle_id"
	CircleColumn_Title       = "circle.title"
	CircleColumn_Description = "circle.description"
	CircleColumn_Handle      = "circle.handle"
	CircleColumn_ImageURI    = "circle.image_uri"
	CircleColumn_Visibility  = "circle.visibility_level"
)

var CircleFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CircleField_CircleId:    {CircleColumn_CircleId},
	cmodel.CircleField_Title:       {CircleColumn_Title},
	cmodel.CircleField_Description: {CircleColumn_Description},
	cmodel.CircleField_Handle:      {CircleColumn_Handle},
	cmodel.CircleField_ImageURI:    {CircleColumn_ImageURI},
	cmodel.CircleField_Visibility:  {CircleColumn_Visibility},

	cmodel.CircleField_CircleAccess: {
		CircleAccessColumn_CircleAccessId,
		CircleAccessColumn_PermissionLevel,
		CircleAccessColumn_State,
	},
})
var UpdateCircleFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	cmodel.CircleField_Title:       {CircleColumn_Title},
	cmodel.CircleField_Description: {CircleColumn_Description},
	cmodel.CircleField_Visibility:  {CircleColumn_Visibility},
})

var CircleSQLConverter = filter.NewSQLConverter(map[string]string{
	"visibility":       CircleColumn_Visibility,
	"permission_level": CircleAccessColumn_PermissionLevel,
	"state":            CircleAccessColumn_State,
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
}

// TableName sets the table name for the Circle model.
func (Circle) TableName() string {
	return "circle"
}
