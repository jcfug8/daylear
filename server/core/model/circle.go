package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
)

// Circle defines the model for a circle.
type Circle struct {
	Id       CircleId
	Parent   CircleParent
	Title    string
	IsPublic bool
}

// CircleId defines the identifier for a circle.
type CircleId struct {
	CircleId int64 `aip_pattern:"key=circle,public_circle"`
}

// CircleParent defines the parent for a circle.
type CircleParent struct {
	UserId int64 `aip_pattern:"key=user"`
}

// ----------------------------------------------------------------------------
// Fields

// CircleFields defines the circle fields.
var CircleFields = circleFields{
	Id:       "id",
	Title:    "title",
	Parent:   "parent",
	IsPublic: "is_public",
}

type circleFields struct {
	Id       string
	Title    string
	Parent   string
	IsPublic string
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Title,
		fields.Parent,
		fields.IsPublic,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields circleFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Title,
		fields.IsPublic,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
