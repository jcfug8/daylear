package model

import (
	"time"

	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
)

// EventData is the GORM model for event data.
// This table contains all the common event data shared between parent events and instances.
type EventData struct {
	EventDataId int64 `gorm:"primaryKey;column:event_data_id;autoIncrement;<-:false"`

	CalendarId int64 `gorm:"column:calendar_id;not null;index"`

	// Event timing
	StartTime time.Time  `gorm:"column:start_time;not null;index"`
	EndTime   *time.Time `gorm:"column:end_time;index"`
	IsAllDay  bool       `gorm:"column:is_all_day;not null;default:false"`

	// Event content
	Title       *string `gorm:"column:title"`
	Description *string `gorm:"column:description"`
	Location    *string `gorm:"column:location"`
	URL         *string `gorm:"column:url"`

	// Event metadata
	Status pb.Event_State `gorm:"column:status;not null;default:1"`
	Class  pb.Event_Class `gorm:"column:class;not null;default:1"`

	// Timestamps
	CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime"`

	// Relationships
	Events []*Event `gorm:"foreignKey:EventDataId"`
}

// TableName sets the table name for the EventData model.
func (EventData) TableName() string {
	return "event_data"
}
