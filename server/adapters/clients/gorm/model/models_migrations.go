package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		&Recipe{},
		&RecipeAccess{},
		&RecipeFavorite{},
		&User{},
		&UserAccess{},
		&AccessKey{},
		&UserFavorite{},
		&Circle{},
		&CircleAccess{},
		&CircleFavorite{},
		&Calendar{},
		&CalendarAccess{},
		&CalendarFavorite{},
		&Event{},
		&EventData{},
	}
}
