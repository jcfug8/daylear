package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type calendarClient interface {
	CreateCalendar(ctx context.Context, calendar model.Calendar, fields []string) (model.Calendar, error)
	DeleteCalendar(ctx context.Context, id model.CalendarId) (model.Calendar, error)
	GetCalendar(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId, fields []string) (model.Calendar, error)
	ListCalendars(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fields []string) ([]model.Calendar, error)
	UpdateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, fields []string) (model.Calendar, error)

	CreateCalendarFavorite(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId) error
	DeleteCalendarFavorite(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId) error

	FindStandardUserCalendarAccess(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId) (model.CalendarAccess, error)
	FindDelegatedCircleCalendarAccess(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId) (model.CalendarAccess, model.CircleAccess, error)
	FindDelegatedUserCalendarAccess(ctx context.Context, authAccount model.AuthAccount, id model.CalendarId) (model.CalendarAccess, model.UserAccess, error)

	CreateCalendarAccess(ctx context.Context, access model.CalendarAccess, fields []string) (model.CalendarAccess, error)
	DeleteCalendarAccess(ctx context.Context, parent model.CalendarAccessParent, id model.CalendarAccessId) error
	BulkDeleteCalendarAccess(ctx context.Context, parent model.CalendarAccessParent) error
	GetCalendarAccess(ctx context.Context, parent model.CalendarAccessParent, id model.CalendarAccessId, fields []string) (model.CalendarAccess, error)
	ListCalendarAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.CalendarAccess, error)
	UpdateCalendarAccess(ctx context.Context, access model.CalendarAccess, fields []string) (model.CalendarAccess, error)
}
