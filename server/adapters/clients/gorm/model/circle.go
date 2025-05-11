package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	coremodel "github.com/jcfug8/daylear/server/core/model"
)

// CircleMap maps the Circle fields to their corresponding fields in the core model.
var CircleMap = masks.NewFieldMap().
	MapFieldToFields(coremodel.CircleFields.Id,
		CircleFields.CircleId).
	MapFieldToFields(coremodel.CircleFields.Title,
		CircleFields.Title).
	MapFieldToFields(coremodel.CircleFields.IsPublic,
		CircleFields.IsPublic)

// CircleFields defines the circle fields in the GORM model.
var CircleFields = circleFields{
	CircleId: "circle_id",
	Title:    "title",
	IsPublic: "is_public",
}

type circleFields struct {
	CircleId string
	Title    string
	IsPublic string
}

// Map maps the circle fields to their corresponding model values.
func (fields circleFields) Map(m Circle) map[string]any {
	return map[string]any{
		fields.CircleId: m.CircleId,
		fields.Title:    m.Title,
		fields.IsPublic: m.IsPublic,
	}
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.CircleId,
		fields.Title,
		fields.IsPublic,
	}
}

// Circle is the GORM model for a circle.
type Circle struct {
	CircleId int64  `gorm:"primaryKey;column:circle_id;autoIncrement;<-:false"`
	Title    string `gorm:"column:title;not null"`
	IsPublic bool   `gorm:"column:is_public;not null;default:false"`
}

// TableName sets the table name for the Circle model.
func (Circle) TableName() string {
	return "circle"
}
