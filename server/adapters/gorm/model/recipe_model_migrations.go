package model

// Recipe -
type Recipe struct {
	RecipeId         int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title            string
	Description      string
	Directions       []byte `gorm:"type:jsonb"`
	ImageURI         string
	IngredientGroups []byte `gorm:"type:jsonb"`
}

// TableName -
func (Recipe) TableName() string {
	return "recipe"
}
