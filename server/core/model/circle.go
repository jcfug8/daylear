package model

import (
	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// Circle defines the model for a circle.
type Circle struct {
	Id         CircleId
	Title      string
	Visibility types.VisibilityLevel
}

// CircleId defines the identifier for a circle.
type CircleId struct {
	CircleId int64 `aip_pattern:"key=circle,public_circle"`
}

// ----------------------------------------------------------------------------
// Fields

// CircleFields defines the circle fields.
var CircleFields = circleFields{
	Id:       "id",
	Title:    "title",
	IsPublic: "is_public",
}

type circleFields struct {
	Id       string
	Title    string
	IsPublic string
}

// Mask returns a FieldMask for the circle fields.
func (fields circleFields) Mask() []string {
	return []string{
		fields.Id,
		fields.Title,
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
