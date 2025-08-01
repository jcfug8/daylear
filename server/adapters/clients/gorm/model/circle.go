package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	coremodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// CircleMap maps the Circle fields to their corresponding fields in the core model.
var CircleMap = masks.NewFieldMap().
	MapFieldToFields(coremodel.CircleFields.Id,
		CircleFields.CircleId).
	MapFieldToFields(coremodel.CircleFields.Title,
		CircleFields.Title).
	MapFieldToFields(coremodel.CircleFields.Visibility,
		CircleFields.Visibility).
	MapFieldToFields(coremodel.CircleFields.ImageURI,
		CircleFields.ImageURI).
	MapFieldToFields(coremodel.CircleFields.Handle,
		CircleFields.Handle).
	MapFieldToFields(coremodel.CircleFields.Description,
		CircleFields.Description)

// CircleFields defines the circle fields in the GORM model.
var CircleFields = circleFields{
	CircleId:    "circle.circle_id",
	Title:       "circle.title",
	Description: "circle.description",
	Handle:      "circle.handle",
	ImageURI:    "circle.image_uri",
	Visibility:  "circle.visibility_level",
	Permission:  "circle_access.permission_level",
	State:       "circle_access.state",
}

type circleFields struct {
	CircleId    string
	Title       string
	Description string
	Handle      string
	ImageURI    string
	Visibility  string
	Permission  string
	State       string
}

// Map maps the circle fields to their corresponding model values.
func (fields circleFields) Map(m Circle) map[string]any {
	return map[string]any{
		fields.CircleId:    m.CircleId,
		fields.Title:       m.Title,
		fields.Description: m.Description,
		fields.Handle:      m.Handle,
		fields.ImageURI:    m.ImageURI,
		fields.Visibility:  m.VisibilityLevel,
		fields.Permission:  m.PermissionLevel,
		fields.State:       m.State,
	}
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.CircleId,
		fields.Title,
		fields.Description,
		fields.Handle,
		fields.ImageURI,
		fields.Visibility,
		fields.Permission,
		fields.State,
	}
}

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
