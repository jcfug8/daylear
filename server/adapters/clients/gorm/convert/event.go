package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

func EventFromCoreModel(mEvent cmodel.Event) (gmodel.Event, gmodel.EventData, error) {
	// Convert location to Point type
	var locationPoint *gmodel.Point
	if mEvent.Location.Latitude != 0 || mEvent.Location.Longitude != 0 {
		locationPoint = &gmodel.Point{
			Longitude: mEvent.Location.Longitude,
			Latitude:  mEvent.Location.Latitude,
		}
	}

	// Create EventData
	eventData := gmodel.EventData{
		CalendarId:  mEvent.Parent.CalendarId,
		Title:       &mEvent.Title,
		Description: &mEvent.Description,
		Location:    locationPoint,
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
	var location cmodel.LatLng
	if mEventData.Location != nil {
		location = cmodel.LatLng{
			Longitude: mEventData.Location.Longitude,
			Latitude:  mEventData.Location.Latitude,
		}
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
		CreateTime:         *mEventData.CreateTime,
		UpdateTime:         *mEventData.UpdateTime,
		StartTime:          mEvent.StartTime,
		EndTime:            mEvent.EndTime,
		IsAllDay:           mEvent.IsAllDay,
		Title:              *mEventData.Title,
		Description:        *mEventData.Description,
		Location:           location,
		URL:                *mEventData.URL,
	}

	return event, nil
}
