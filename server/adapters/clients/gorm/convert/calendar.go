package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// CalendarFromCoreModel converts a core Calendar to a GORM Calendar
func CalendarFromCoreModel(calendar cmodel.Calendar) (gmodel.Calendar, error) {
	return gmodel.Calendar{
		CalendarId:      calendar.CalendarId.CalendarId,
		Title:           calendar.Title,
		Description:     calendar.Description,
		VisibilityLevel: calendar.VisibilityLevel,
		CreateTime:      calendar.CreateTime,
		UpdateTime:      calendar.UpdateTime,
	}, nil
}

// CalendarToCoreModel converts a GORM Calendar to a core Calendar
func CalendarToCoreModel(gormCalendar gmodel.Calendar) (cmodel.Calendar, error) {
	calendar := cmodel.Calendar{
		CalendarId:      cmodel.CalendarId{CalendarId: gormCalendar.CalendarId},
		Title:           gormCalendar.Title,
		Description:     gormCalendar.Description,
		VisibilityLevel: gormCalendar.VisibilityLevel,
		CreateTime:      gormCalendar.CreateTime,
		UpdateTime:      gormCalendar.UpdateTime,
		Favorited:       gormCalendar.CalendarFavoriteId != 0,
		CalendarAccess: cmodel.CalendarAccess{
			CalendarAccessParent: cmodel.CalendarAccessParent{CalendarId: gormCalendar.CalendarId},
			CalendarAccessId:     cmodel.CalendarAccessId{CalendarAccessId: gormCalendar.CalendarAccessId},
			PermissionLevel:      gormCalendar.PermissionLevel,
			State:                gormCalendar.State,
			AcceptTarget:         gormCalendar.AcceptTarget,
		},
	}

	return calendar, nil
}
