package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		&Recipe{},
		&RecipeAccess{},
		&RecipeFavorite{},
		&User{},
		&UserAccess{},
		&Circle{},
		&CircleAccess{},
		&Calendar{},
		&CalendarAccess{},
		&Event{},
		&EventData{},
	}
}
