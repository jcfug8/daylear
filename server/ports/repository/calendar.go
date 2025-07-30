package repository

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

type calendarClient interface {
	CreateCalendar(ctx context.Context, calendar model.Calendar) (model.Calendar, error)
	DeleteCalendar(ctx context.Context, authAccount model.AuthAccount, id int64) (model.Calendar, error)
	GetCalendar(ctx context.Context, authAccount model.AuthAccount, id int64) (model.Calendar, error)
	ListCalendars(ctx context.Context, authAccount model.AuthAccount, pageSize int32, offset int64, filter string, fieldMask []string) ([]model.Calendar, error)
	UpdateCalendar(ctx context.Context, authAccount model.AuthAccount, calendar model.Calendar, updateMask []string) (model.Calendar, error)

	CreateCalendarAccess(ctx context.Context, access model.CalendarAccess) (model.CalendarAccess, error)
	DeleteCalendarAccess(ctx context.Context, parent model.CalendarAccessParent, id model.CalendarAccessId) error
	BulkDeleteCalendarAccess(ctx context.Context, parent model.CalendarAccessParent) error
	GetCalendarAccess(ctx context.Context, parent model.CalendarAccessParent, id model.CalendarAccessId) (model.CalendarAccess, error)
	ListCalendarAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.CalendarAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.CalendarAccess, error)
	UpdateCalendarAccess(ctx context.Context, access model.CalendarAccess, updateMask []string) (model.CalendarAccess, error)
}
