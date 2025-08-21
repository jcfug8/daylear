package model

import "time"

const (
	CalendarFavoriteTable = "calendar_favorite"
)

const (
	CalendarFavoriteFields_CalendarFavoriteId = "calendar_favorite_id"
	CalendarFavoriteFields_CalendarId         = "calendar_id"
	CalendarFavoriteFields_UserId             = "user_id"
	CalendarFavoriteFields_CircleId           = "circle_id"
	CalendarFavoriteFields_CreateTime         = "create_time"
)

type CalendarFavorite struct {
	CalendarFavoriteId int64     `gorm:"primaryKey;bigint;not null;<-:false"`
	CalendarId         int64     `gorm:"not null;index;uniqueIndex:idx_calendar_id_user_id_circle_id"`
	UserId             int64     `gorm:"index;uniqueIndex:idx_calendar_id_user_id_circle_id"`
	CircleId           int64     `gorm:"index;uniqueIndex:idx_calendar_id_user_id_circle_id"`
	CreateTime         time.Time `gorm:"not null;default:now()"`
}

func (CalendarFavorite) TableName() string {
	return CalendarFavoriteTable
}
