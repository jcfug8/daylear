package model

import (
	"time"

	"github.com/jcfug8/daylear/server/genapi/api/types"
)

var _ ResourceId = CalendarId{}

// CalendarFields defines the calendar fields.
const (
	CalendarField_Parent      = "parent"
	CalendarField_CalendarId  = "id"
	CalendarField_Title       = "title"
	CalendarField_Description = "description"
	CalendarField_Visibility  = "visibility"

	CalendarField_CalendarAccess = "calendar_access"
)

// Calendar represents a VCALENDAR component in iCalendar.
// This is the top-level container that holds all calendar components.
type Calendar struct {
	// Parent is the parent of the calendar
	Parent CalendarParent
	// CalendarId is the unique identifier for the calendar
	CalendarId CalendarId
	// Title is the title of the calendar
	Title string
	// Description is the description of the calendar
	Description string
	// VisibilityLevel is the visibility level of the calendar
	VisibilityLevel types.VisibilityLevel
	// CreateTime is the time the calendar was created
	CreateTime time.Time
	// UpdateTime is the time the calendar was last updated
	UpdateTime time.Time

	CalendarAccess CalendarAccess
}

type CalendarParent struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type CalendarId struct {
	CalendarId int64 `aip_pattern:"key=calendar"`
}

// isResourceId - implements the ResourceId interface.
func (c CalendarId) isResourceId() {
}
