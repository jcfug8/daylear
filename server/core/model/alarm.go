package model

import (
	"time"
)

type Alarm struct {
	// UID is the unique identifier for this alarm
	AlarmId string `json:"alarmId"`
	// Trigger specifies when the alarm should fire
	Trigger *Trigger `json:"trigger"`
	// Description is the text to display or include in notification
	// Required for DISPLAY and EMAIL actions
	Description *string `json:"description,omitempty"`
	// Summary is the summary text for the alarm
	// Optional for DISPLAY action
	Summary *string `json:"summary,omitempty"`
	// DTStamp is the creation timestamp of the alarm
	CreateTime *time.Time `json:"createTime,omitempty"`
	// UpdateTime is the last update timestamp of the alarm
	UpdateTime *time.Time `json:"updateTime,omitempty"`
}

type Trigger struct {
	// Duration represents the relative duration (e.g., -PT15M for 15 minutes before)
	// Only used for relative triggers
	Duration *time.Duration `json:"duration,omitempty"`
	// DateTime represents the absolute date/time for the trigger
	// Only used for absolute triggers
	DateTime *time.Time `json:"dateTime,omitempty"`
}
