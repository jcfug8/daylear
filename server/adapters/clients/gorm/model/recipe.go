package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/masks"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeMap maps the Recipe fields to their corresponding
// fields in the model.
var RecipeMap = masks.NewFieldMap().
	MapFieldToFields(model.RecipeFields.Id,
		RecipeFields.RecipeId).
	MapFieldToFields(model.RecipeFields.Title,
		RecipeFields.Title).
	MapFieldToFields(model.RecipeFields.Description,
		RecipeFields.Description).
	MapFieldToFields(model.RecipeFields.Directions,
		RecipeFields.Directions).
	MapFieldToFields(model.RecipeFields.ImageURI,
		RecipeFields.ImageURI).
	MapFieldToFields(model.RecipeFields.IngredientGroups,
		RecipeFields.IngredientGroups).
	MapFieldToFields(model.RecipeFields.VisibilityLevel,
		RecipeFields.VisibilityLevel).
	MapFieldToFields(model.RecipeFields.PermissionLevel,
		RecipeFields.PermissionLevel).
	MapFieldToFields(model.RecipeFields.Citation,
		RecipeFields.Citation).
	MapFieldToFields(model.RecipeFields.CookDurationSeconds,
		RecipeFields.CookDurationSeconds).
	MapFieldToFields(model.RecipeFields.PrepDurationSeconds,
		RecipeFields.PrepDurationSeconds).
	MapFieldToFields(model.RecipeFields.TotalDurationSeconds,
		RecipeFields.TotalDurationSeconds).
	MapFieldToFields(model.RecipeFields.CookingMethod,
		RecipeFields.CookingMethod).
	MapFieldToFields(model.RecipeFields.Categories,
		RecipeFields.Categories).
	MapFieldToFields(model.RecipeFields.YieldAmount,
		RecipeFields.YieldAmount).
	MapFieldToFields(model.RecipeFields.Cuisines,
		RecipeFields.Cuisines).
	MapFieldToFields(model.RecipeFields.CreateTime,
		RecipeFields.CreateTime).
	MapFieldToFields(model.RecipeFields.UpdateTime,
		RecipeFields.UpdateTime)

// RecipeFields defines the recipe fields.
var RecipeFields = recipeFields{
	RecipeId:             "recipe.recipe_id",
	Title:                "recipe.title",
	Description:          "recipe.description",
	Directions:           "recipe.directions",
	ImageURI:             "recipe.image_uri",
	IngredientGroups:     "recipe.ingredient_groups",
	VisibilityLevel:      "recipe.visibility_level",
	PermissionLevel:      "recipe_access.permission_level",
	State:                "recipe_access.state",
	Citation:             "recipe.citation",
	CookDurationSeconds:  "recipe.cook_duration_seconds",
	PrepDurationSeconds:  "recipe.prep_duration_seconds",
	TotalDurationSeconds: "recipe.total_duration_seconds",
	CookingMethod:        "recipe.cooking_method",
	Categories:           "recipe.categories",
	YieldAmount:          "recipe.yield_amount",
	Cuisines:             "recipe.cuisines",
	CreateTime:           "recipe.create_time",
	UpdateTime:           "recipe.update_time",
}

type recipeFields struct {
	RecipeId             string
	Title                string
	Description          string
	Directions           string
	ImageURI             string
	IngredientGroups     string
	VisibilityLevel      string
	PermissionLevel      string
	State                string
	Citation             string
	CookDurationSeconds  string
	PrepDurationSeconds  string
	TotalDurationSeconds string
	CookingMethod        string
	Categories           string
	YieldAmount          string
	Cuisines             string
	CreateTime           string
	UpdateTime           string
}

// Map maps the recipe fields to their corresponding model values.
func (fields recipeFields) Map(m Recipe) map[string]any {
	return map[string]any{
		fields.RecipeId:             m.RecipeId,
		fields.Title:                m.Title,
		fields.Description:          m.Description,
		fields.Directions:           m.Directions,
		fields.ImageURI:             m.ImageURI,
		fields.IngredientGroups:     m.IngredientGroups,
		fields.VisibilityLevel:      m.VisibilityLevel,
		fields.PermissionLevel:      m.PermissionLevel,
		fields.Citation:             m.Citation,
		fields.CookDurationSeconds:  m.CookDurationSeconds,
		fields.PrepDurationSeconds:  m.PrepDurationSeconds,
		fields.TotalDurationSeconds: m.TotalDurationSeconds,
		fields.CookingMethod:        m.CookingMethod,
		fields.Categories:           m.Categories,
		fields.YieldAmount:          m.YieldAmount,
		fields.Cuisines:             m.Cuisines,
		fields.CreateTime:           m.CreateTime,
		fields.UpdateTime:           m.UpdateTime,
	}
}

// Mask returns a FieldMask for the recipe fields.
func (fields recipeFields) Mask() []string {
	return []string{
		fields.RecipeId,
		fields.Title,
		fields.Description,
		fields.Directions,
		fields.ImageURI,
		fields.IngredientGroups,
		fields.VisibilityLevel,
		fields.PermissionLevel,
		fields.Citation,
		fields.CookDurationSeconds,
		fields.PrepDurationSeconds,
		fields.TotalDurationSeconds,
		fields.CookingMethod,
		fields.Categories,
		fields.YieldAmount,
		fields.Cuisines,
		fields.CreateTime,
		fields.UpdateTime,
	}
}

// Recipe -
type Recipe struct {
	RecipeId         int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title            string
	Description      string
	Directions       []byte `gorm:"type:jsonb"`
	ImageURI         string
	IngredientGroups []byte                `gorm:"type:jsonb"`
	VisibilityLevel  types.VisibilityLevel `gorm:"not null;default:300"`

	// New fields
	Citation             string    `gorm:"type:varchar(512)"`
	CookDuration         int64     `gorm:"column:cook_duration_nanos;type:bigint"` // nanoseconds
	CookingMethod        string    `gorm:"type:varchar(64)"`
	Categories           []byte    `gorm:"type:jsonb"`
	YieldAmount          string    `gorm:"type:varchar(64)"`
	Cuisines             []byte    `gorm:"type:jsonb"`
	CreateTime           time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime           time.Time `gorm:"column:update_time;autoUpdateTime"`
	PrepDurationSeconds  int64     `gorm:"column:prep_duration_seconds;type:bigint"`
	CookDurationSeconds  int64     `gorm:"column:cook_duration_seconds;type:bigint"`
	TotalDurationSeconds int64     `gorm:"column:total_duration_seconds;type:bigint"`

	// RecipeAccess data
	RecipeAccessId  int64                 `gorm:"->;-:migration"` // only used for read from a join
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
}

// TableName -
func (Recipe) TableName() string {
	return "recipe"
}
