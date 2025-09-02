package model

import "time"

const (
	ListFavoriteTable = "list_favorite"
)

const (
	ListFavoriteFields_ListFavoriteId = "list_favorite_id"
	ListFavoriteFields_ListId         = "list_id"
	ListFavoriteFields_UserId         = "user_id"
	ListFavoriteFields_CircleId       = "circle_id"
	ListFavoriteFields_CreateTime     = "create_time"
)

type ListFavorite struct {
	ListFavoriteId int64     `gorm:"primaryKey;bigint;not null;<-:false"`
	ListId         int64     `gorm:"not null;index;uniqueIndex:idx_list_id_user_id_circle_id"`
	UserId         int64     `gorm:"index;uniqueIndex:idx_list_id_user_id_circle_id"`
	CircleId       int64     `gorm:"index;uniqueIndex:idx_list_id_user_id_circle_id"`
	CreateTime     time.Time `gorm:"not null;default:now()"`
}

func (ListFavorite) TableName() string {
	return ListFavoriteTable
}
