package model

// Ingredient -
type Ingredient struct {
	IngredientId int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title        string
}

// TableName -
func (Ingredient) TableName() string {
	return "ingredient"
}
