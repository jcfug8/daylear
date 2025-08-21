package model

import "time"

const (
	UserFavoriteTable = "user_favorite"
)

const (
	UserFavoriteFields_UserFavoriteId   = "user_favorite_id"
	UserFavoriteFields_FavoritedUserId  = "favorited_user_id"
	UserFavoriteFields_FavoritingUserId = "favoriting_user_id"
	UserFavoriteFields_CreateTime       = "create_time"
)

type UserFavorite struct {
	UserFavoriteId   int64     `gorm:"primaryKey;bigint;not null;<-:false"`
	FavoritedUserId  int64     `gorm:"not null;index;uniqueIndex:idx_favorited_user_id_favoriting_user_id"`
	FavoritingUserId int64     `gorm:"index;uniqueIndex:idx_favorited_user_id_favoriting_user_id"`
	CreateTime       time.Time `gorm:"not null;default:now()"`
}

func (UserFavorite) TableName() string {
	return UserFavoriteTable
}
