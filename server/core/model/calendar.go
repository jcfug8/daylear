package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/genapi/api/types"
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

// ----------------------------------------------------------------------------
// Fields

// CalendarFields defines the calendar fields.
var CalendarFields = calendarFields{
	CalendarId:  "calendar_id",
	Title:       "title",
	Description: "description",
	Visibility:  "visibility",
	AccessId:    "access_id",
	Permission:  "permission",
	State:       "state",
}

type calendarFields struct {
	CalendarId  string
	Title       string
	Description string
	Visibility  string
	AccessId    string
	Permission  string
	State       string
}

// Mask returns a FieldMask for the calendar fields.
func (fields calendarFields) Mask() []string {
	return []string{
		fields.CalendarId,
		fields.Title,
		fields.Description,
		fields.Visibility,
		fields.AccessId,
		fields.Permission,
		fields.State,
	}
}

// UpdateMask returns the subset of provided fields that can be updated.
func (fields calendarFields) UpdateMask(mask []string) []string {
	updatable := []string{
		fields.Title,
		fields.Description,
		fields.Visibility,
	}

	if len(mask) == 0 {
		return updatable
	}

	return masks.Intersection(updatable, mask)
}
