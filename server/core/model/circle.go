package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ ResourceId = CircleId{}

const (
	CircleField_Parent      = "parent"
	CircleField_CircleId    = "id"
	CircleField_Title       = "title"
	CircleField_Description = "description"
	CircleField_Handle      = "handle"
	CircleField_ImageURI    = "image_uri"
	CircleField_Visibility  = "visibility"
	CircleField_Favorited   = "favorited"

	CircleField_CircleAccess = "circle_access"
)

// Circle defines the model for a circle.
type Circle struct {
	Id              CircleId
	Parent          CircleParent
	Title           string
	Description     string
	Handle          string
	ImageURI        string
	VisibilityLevel types.VisibilityLevel
	Favorited       bool
	CircleAccess    CircleAccess
}

// CircleId defines the identifier for a circle.
type CircleId struct {
	CircleId int64 `aip_pattern:"key=circle"`
}

// isResourceId - implements the ResourceId interface.
func (c CircleId) isResourceId() {
}

type CircleParent struct {
	UserId int64 `aip_pattern:"key=user"`
}
