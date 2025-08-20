package model

import "time"

const (
	RecipeFavoriteTable = "recipe_favorite"
)

const (
	RecipeFavoriteFields_RecipeFavoriteId = "recipe_favorite_id"
	RecipeFavoriteFields_RecipeId         = "recipe_id"
	RecipeFavoriteFields_UserId           = "user_id"
	RecipeFavoriteFields_CreateTime       = "create_time"
)

type RecipeFavorite struct {
	RecipeFavoriteId int64     `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId         int64     `gorm:"not null;index;uniqueIndex:idx_recipe_id_user_id"`
	UserId           int64     `gorm:"not null;index;uniqueIndex:idx_recipe_id_user_id"`
	CreateTime       time.Time `gorm:"not null;default:now()"`
}

func (RecipeFavorite) TableName() string {
	return RecipeFavoriteTable
}
