package model

import "time"

const (
	CircleFavoriteTable = "circle_favorite"
)

const (
	CircleFavoriteFields_CircleFavoriteId = "circle_favorite_id"
	CircleFavoriteFields_CircleId         = "circle_id"
	CircleFavoriteFields_UserId           = "user_id"
	CircleFavoriteFields_CreateTime       = "create_time"
)

type CircleFavorite struct {
	CircleFavoriteId int64     `gorm:"primaryKey;bigint;not null;<-:false"`
	CircleId         int64     `gorm:"not null;index;uniqueIndex:idx_circle_id_user_id"`
	UserId           int64     `gorm:"index;uniqueIndex:idx_circle_id_user_id"`
	CreateTime       time.Time `gorm:"not null;default:now()"`
}

func (CircleFavorite) TableName() string {
	return CircleFavoriteTable
}
