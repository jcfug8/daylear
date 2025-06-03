package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		// &Event{},
		// &Family{},
		// &Meal{},
		&Recipe{},
		&RecipeCircle{},
		&RecipeUser{},
		&RecipeIngredient{},
		&Ingredient{},
		&User{},
		&Circle{},
		&CircleUser{},
	}
}
