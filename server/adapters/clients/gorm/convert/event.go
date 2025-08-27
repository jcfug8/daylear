package convert

import (
	"time"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

func EventFromCoreModel(mEvent cmodel.Event) (gmodel.Event, gmodel.EventData, error) {
	// Convert location to Point type
	var locationPoint *gmodel.Point
	if mEvent.Geo.Latitude != 0 || mEvent.Geo.Longitude != 0 {
		locationPoint = &gmodel.Point{
			Longitude: mEvent.Geo.Longitude,
			Latitude:  mEvent.Geo.Latitude,
		}
	}

	// Create EventData
	eventData := gmodel.EventData{
		CalendarId:  mEvent.Parent.CalendarId,
		Title:       &mEvent.Title,
		Description: &mEvent.Description,
		Location:    &mEvent.Location,
		Geo:         locationPoint,
		URL:         &mEvent.URL,
		CreateTime:  &mEvent.CreateTime,
		UpdateTime:  &mEvent.UpdateTime,
	}

	// Create Event
	event := gmodel.Event{
		EventId:            mEvent.Id.EventId,
		ParentEventId:      mEvent.ParentEventId,
		OverridenStartTime: mEvent.OverridenStartTime,
		RecurrenceRule:     mEvent.RecurrenceRule,
		ExcludedDates:      mEvent.ExcludedDates,
		AdditionalDates:    mEvent.AdditionalDates,
		StartTime:          mEvent.StartTime,
		EndTime:            mEvent.EndTime,
		IsAllDay:           mEvent.IsAllDay,
		RecurrenceEndTime:  mEvent.RecurrenceEndTime,
	}

	// Handle recurring event logic
	if mEvent.ParentEventId != nil {
		// This is an instance of a recurring event
		event.ParentEventId = mEvent.ParentEventId
		if mEvent.OverridenStartTime != nil {
			event.OverridenStartTime = mEvent.OverridenStartTime
		}
	} else if mEvent.RecurrenceRule != nil && *mEvent.RecurrenceRule != "" {
		// This is a parent recurring event
		// ParentEventId remains nil
	} else {
		// This is a single event
		// ParentEventId remains nil
	}

	return event, eventData, nil
}

func EventToCoreModel(mEvent gmodel.Event, mEventData gmodel.EventData) (cmodel.Event, error) {
	// Convert Point to LatLng
	var geo cmodel.LatLng
	if mEventData.Geo != nil {
		geo = cmodel.LatLng{
			Longitude: mEventData.Geo.Longitude,
			Latitude:  mEventData.Geo.Latitude,
		}
	}

	createTime := time.Time{}
	if mEventData.CreateTime != nil {
		createTime = *mEventData.CreateTime
	}

	updateTime := time.Time{}
	if mEventData.UpdateTime != nil {
		updateTime = *mEventData.UpdateTime
	}

	var deleteTime *time.Time
	if mEventData.DeleteTime != nil {
		deleteTime = mEventData.DeleteTime
	}

	title := ""
	if mEventData.Title != nil {
		title = *mEventData.Title
	}

	description := ""
	if mEventData.Description != nil {
		description = *mEventData.Description
	}

	location := ""
	if mEventData.Location != nil {
		location = *mEventData.Location
	}

	url := ""
	if mEventData.URL != nil {
		url = *mEventData.URL
	}

	recurrenceEndTime := time.Time{}
	if mEvent.RecurrenceEndTime != nil {
		recurrenceEndTime = *mEvent.RecurrenceEndTime
	}

	// Create core Event
	event := cmodel.Event{
		Parent: cmodel.EventParent{
			CalendarId: mEventData.CalendarId,
		},
		Id: cmodel.EventId{
			EventId: mEvent.EventId,
		},
		RecurrenceRule:     mEvent.RecurrenceRule,
		ExcludedDates:      mEvent.ExcludedDates,
		AdditionalDates:    mEvent.AdditionalDates,
		ParentEventId:      mEvent.ParentEventId,
		OverridenStartTime: mEvent.OverridenStartTime,
		CreateTime:         createTime,
		UpdateTime:         updateTime,
		StartTime:          mEvent.StartTime,
		EndTime:            mEvent.EndTime,
		IsAllDay:           mEvent.IsAllDay,
		Title:              title,
		Description:        description,
		Location:           location,
		Geo:                geo,
		URL:                url,
		RecurrenceEndTime:  &recurrenceEndTime,
		DeleteTime:         deleteTime,
	}

	return event, nil
}
