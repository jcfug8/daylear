package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		&Recipe{},
		&RecipeAccess{},
		&RecipeIngredient{},
		&Ingredient{},
		&User{},
		&Circle{},
		&CircleAccess{},
	}
}
