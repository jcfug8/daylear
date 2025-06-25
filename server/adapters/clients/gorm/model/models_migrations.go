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
		&RecipeAccess{},
		&RecipeIngredient{},
		&Ingredient{},
		&User{},
		&Circle{},
		&CircleUser{},
	}
}
