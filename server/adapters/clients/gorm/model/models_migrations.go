package model

// AllModels -
func AllModels() []interface{} {
	return []interface{}{
		&Recipe{},
		&RecipeAccess{},
		&User{},
		&UserAccess{},
		&Circle{},
		&CircleAccess{},
	}
}
