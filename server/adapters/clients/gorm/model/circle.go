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
		CircleFields.Title)

// CircleFields defines the circle fields in the GORM model.
var CircleFields = circleFields{
	CircleId: "circle_id",
	Title:    "title",
}

type circleFields struct {
	CircleId string
	Title    string
}

// Map maps the circle fields to their corresponding model values.
func (fields circleFields) Map(m Circle) map[string]any {
	return map[string]any{
		fields.CircleId: m.CircleId,
		fields.Title:    m.Title,
	}
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.CircleId,
		fields.Title,
	}
}

// Circle is the GORM model for a circle.
type Circle struct {
	CircleId int64  `gorm:"primaryKey;column:circle_id;autoIncrement"`
	Title    string `gorm:"column:title;not null"`
}

// TableName sets the table name for the Circle model.
func (Circle) TableName() string {
	return "circle"
}
