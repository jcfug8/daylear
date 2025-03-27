package model

import (
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// RecipeIngredient -
type RecipeIngredient struct {
	RecipeIngredientId   int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	RecipeId             int64
	IngredientId         int64
	MeasurementAmount    float64
	MeasurementType      pb.Recipe_MeasurementType
	IngredientGroupIndex int
}

// TableName -
func (RecipeIngredient) TableName() string {
	return "recipe_ingredient"
}
