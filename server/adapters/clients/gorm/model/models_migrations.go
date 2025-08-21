package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		&Recipe{},
		&RecipeAccess{},
		&RecipeFavorite{},
		&User{},
		&UserAccess{},
		&UserFavorite{},
		&Circle{},
		&CircleAccess{},
		&Calendar{},
		&CalendarAccess{},
		&Event{},
		&EventData{},
	}
}
