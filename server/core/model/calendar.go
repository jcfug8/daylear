package model

// Calendar represents a VCALENDAR component in iCalendar.
// This is the top-level container that holds all calendar components.
type Calendar struct {
	// Parent is the parent of the calendar
	Parent CalendarParent
	// CalendarId is the unique identifier for the calendar
	CalendarId string `aip_pattern:"calendar"`
	// Title is the title of the calendar
	Title string
	// Description is the description of the calendar
	Description string
}

type CalendarParent struct {
	UserId   string `aip_pattern:"user"`
	CircleId string `aip_pattern:"circle"`
}
