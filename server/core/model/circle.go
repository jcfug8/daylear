package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
)

// Circle defines the model for a circle.
type Circle struct {
	Id     CircleId
	Parent CircleParent
	Title  string
}

// CircleId defines the identifier for a circle.
type CircleId struct {
	CircleId int64
}

// CircleParent defines the parent for a circle.
type CircleParent struct {
	UserId int64
}

// ----------------------------------------------------------------------------
// Fields

// CircleFields defines the circle fields.
var CircleFields = circleFields{
	Id:     "id",
	Title:  "title",
	Parent: "parent",
}

type circleFields struct {
	Id     string
	Title  string
	Parent string
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Title,
		fields.Parent,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields circleFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Title,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
