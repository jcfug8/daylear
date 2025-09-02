package model

import (
	"time"
)

// ListFavorite represents a user's or circle's favorite list
type ListFavorite struct {
	ListFavoriteId ListFavoriteId
	ListId         ListId
	UserId         UserId
	CircleId       CircleId
	CreateTime     time.Time
}

type ListFavoriteId struct {
	ListFavoriteId int64 `aip_pattern:"key=favorite"`
}
