package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
)

const (
	EventDataTable = "event_data"
)

const (
	EventDataField_EventDataId = "event_data_id"
	EventDataField_CalendarId  = "calendar_id"
	EventDataField_Title       = "title"
	EventDataField_Description = "description"
	EventDataField_Location    = "location"
	EventDataField_URL         = "url"
	EventDataField_CreateTime  = "create_time"
	EventDataField_UpdateTime  = "update_time"
	EventDataField_DeleteTime  = "delete_time"
)

var EventDataFieldMasker = fieldmask.NewSQLFieldMasker(EventData{}, map[string][]fieldmask.Field{
	fieldmask.AlwaysIncludeKey:   {{Name: EventDataField_EventDataId, Table: EventDataTable}},
	model.EventField_Parent:      {{Name: EventDataField_CalendarId, Table: EventDataTable}},
	model.EventField_Title:       {{Name: EventDataField_Title, Table: EventDataTable, Updatable: true}},
	model.EventField_Description: {{Name: EventDataField_Description, Table: EventDataTable, Updatable: true}},
	model.EventField_Location:    {{Name: EventDataField_Location, Table: EventDataTable, Updatable: true}},
	model.EventField_URL:         {{Name: EventDataField_URL, Table: EventDataTable, Updatable: true}},
	model.EventField_CreateTime:  {{Name: EventDataField_CreateTime, Table: EventDataTable}},
	model.EventField_UpdateTime:  {{Name: EventDataField_UpdateTime, Table: EventDataTable}},
	model.EventField_DeleteTime:  {{Name: EventDataField_DeleteTime, Table: EventDataTable}},
})

// Point represents a PostgreSQL point type for storing latitude/longitude coordinates
type Point struct {
	Longitude float64
	Latitude  float64
}

// Value implements the driver.Valuer interface for GORM
func (p Point) Value() (driver.Value, error) {
	if p.Longitude == 0 && p.Latitude == 0 {
		return nil, nil
	}
	return fmt.Sprintf("(%f,%f)", p.Longitude, p.Latitude), nil
}

// Scan implements the sql.Scanner interface for GORM
func (p *Point) Scan(value any) error {
	if value == nil {
		p.Longitude = 0
		p.Latitude = 0
		return nil
	}

	var str string
	switch v := value.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return fmt.Errorf("cannot scan %T into Point", value)
	}

	// Parse PostgreSQL point format "(longitude,latitude)"
	str = strings.Trim(str, "()")
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return fmt.Errorf("invalid point format: %s", str)
	}

	longitude, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return fmt.Errorf("invalid longitude: %w", err)
	}

	latitude, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return fmt.Errorf("invalid latitude: %w", err)
	}

	p.Longitude = longitude
	p.Latitude = latitude
	return nil
}

// EventData is the GORM model for event data.
// This table contains all the common event data shared between parent events and instances.
type EventData struct {
	EventDataId int64 `gorm:"primaryKey;column:event_data_id;autoIncrement;<-:false"`

	CalendarId int64 `gorm:"column:calendar_id;not null;index"`

	// Event content
	Title       *string `gorm:"column:title"`
	Description *string `gorm:"column:description"`
	Location    *string `gorm:"column:location"`
	Geo         *Point  `gorm:"column:geo;type:point"`
	URL         *string `gorm:"column:url"`

	// Timestamps
	CreateTime *time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime *time.Time `gorm:"column:update_time;autoUpdateTime"`
	DeleteTime *time.Time `gorm:"column:delete_time;index"`
}

// TableName sets the table name for the EventData model.
func (EventData) TableName() string {
	return EventDataTable
}
