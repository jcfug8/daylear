package model

import (
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// CalendarAccessFields defines the calendar access fields for filtering.
var CalendarAccessFields = calendarAccessFields{
	Level:           "permission_level",
	State:           "state",
	RecipientUser:   "recipient.user_id",
	RecipientCircle: "recipient.circle_id",
}

type calendarAccessFields struct {
	Level           string
	State           string
	RecipientUser   string
	RecipientCircle string
}

// CalendarAccess represents a user's or circle's access to a calendar
type CalendarAccess struct {
	CalendarAccessParent
	CalendarAccessId
	PermissionLevel types.PermissionLevel
	State           types.AccessState

	Requester CalendarRecipientOrRequester

	Recipient             CalendarRecipientOrRequester
	RecipientUsername     string // username of the recipient (if user)
	RecipientGivenName    string // given name of the recipient (if user)
	RecipientFamilyName   string // family name of the recipient (if user)
	RecipientCircleTitle  string // title of the recipient (if circle)
	RecipientCircleHandle string // handle of the recipient (if circle)
}

type CalendarRecipientOrRequester struct {
	UserId   int64 `aip_pattern:"key=user"`
	CircleId int64 `aip_pattern:"key=circle"`
}

type CalendarAccessParent struct {
	CircleId   int64 `aip_pattern:"key=circle"`
	CalendarId int64 `aip_pattern:"key=calendar"`
}

type CalendarAccessId struct {
	CalendarAccessId int64 `aip_pattern:"key=access"`
}
