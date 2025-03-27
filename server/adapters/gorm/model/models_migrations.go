package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		// &Event{},
		// &Family{},
		// &Meal{},
		&Recipe{},
		&RecipeUser{},
		&RecipeIngredient{},
		&Ingredient{},
		&User{},
	}
}
