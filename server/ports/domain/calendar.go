package domain

import (
	"context"

	model "github.com/jcfug8/daylear/server/core/model"
)

type calendarDomain interface {
	CreateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar) (model.Calendar, error)
	DeleteCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId) (model.Calendar, error)
	GetCalendar(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, id model.CalendarId, fields []string) (model.Calendar, error)
	ListCalendars(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarParent, pageSize int32, offset int64, filter string, fields []string) ([]model.Calendar, error)
	UpdateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, fields []string) (model.Calendar, error)

	CreateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess) (model.CalendarAccess, error)
	DeleteCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) error
	GetCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId, fields []string) (model.CalendarAccess, error)
	ListCalendarAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, pageSize int32, pageOffset int64, filter string, fields []string) ([]model.CalendarAccess, error)
	UpdateCalendarAccess(ctx context.Context, authAccount model.AuthAccount, access model.CalendarAccess, fields []string) (model.CalendarAccess, error)
	AcceptCalendarAccess(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, id model.CalendarAccessId) (model.CalendarAccess, error)
}
