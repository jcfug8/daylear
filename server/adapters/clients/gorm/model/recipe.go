package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	RecipeFields_RecipeId             = "recipe.recipe_id"
	RecipeFields_Title                = "recipe.title"
	RecipeFields_Description          = "recipe.description"
	RecipeFields_Directions           = "recipe.directions"
	RecipeFields_ImageURI             = "recipe.image_uri"
	RecipeFields_IngredientGroups     = "recipe.ingredient_groups"
	RecipeFields_VisibilityLevel      = "recipe.visibility_level"
	RecipeFields_Citation             = "recipe.citation"
	RecipeFields_CookDurationSeconds  = "recipe.cook_duration_seconds"
	RecipeFields_PrepDurationSeconds  = "recipe.prep_duration_seconds"
	RecipeFields_TotalDurationSeconds = "recipe.total_duration_seconds"
	RecipeFields_CookingMethod        = "recipe.cooking_method"
	RecipeFields_Categories           = "recipe.categories"
	RecipeFields_YieldAmount          = "recipe.yield_amount"
	RecipeFields_Cuisines             = "recipe.cuisines"
	RecipeFields_CreateTime           = "recipe.create_time"
	RecipeFields_UpdateTime           = "recipe.update_time"
)

var RecipeFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	model.RecipeField_Id:                   {RecipeFields_RecipeId},
	model.RecipeField_Title:                {RecipeFields_Title},
	model.RecipeField_Description:          {RecipeFields_Description},
	model.RecipeField_Directions:           {RecipeFields_Directions},
	model.RecipeField_ImageURI:             {RecipeFields_ImageURI},
	model.RecipeField_IngredientGroups:     {RecipeFields_IngredientGroups},
	model.RecipeField_VisibilityLevel:      {RecipeFields_VisibilityLevel},
	model.RecipeField_Citation:             {RecipeFields_Citation},
	model.RecipeField_CookDurationSeconds:  {RecipeFields_CookDurationSeconds},
	model.RecipeField_PrepDurationSeconds:  {RecipeFields_PrepDurationSeconds},
	model.RecipeField_TotalDurationSeconds: {RecipeFields_TotalDurationSeconds},
	model.RecipeField_CookingMethod:        {RecipeFields_CookingMethod},
	model.RecipeField_Categories:           {RecipeFields_Categories},
	model.RecipeField_YieldAmount:          {RecipeFields_YieldAmount},
	model.RecipeField_Cuisines:             {RecipeFields_Cuisines},
	model.RecipeField_CreateTime:           {RecipeFields_CreateTime},
	model.RecipeField_UpdateTime:           {RecipeFields_UpdateTime},

	model.RecipeField_RecipeAccess: {
		RecipeAccessFields_RecipeAccessId,
		RecipeAccessFields_PermissionLevel,
		RecipeAccessFields_State,
	},
})
var UpdateRecipeFieldMasker = fieldmask.NewFieldMasker(map[string][]string{
	model.RecipeField_Title:                {RecipeFields_Title},
	model.RecipeField_Description:          {RecipeFields_Description},
	model.RecipeField_Directions:           {RecipeFields_Directions},
	model.RecipeField_ImageURI:             {RecipeFields_ImageURI},
	model.RecipeField_IngredientGroups:     {RecipeFields_IngredientGroups},
	model.RecipeField_VisibilityLevel:      {RecipeFields_VisibilityLevel},
	model.RecipeField_Citation:             {RecipeFields_Citation},
	model.RecipeField_CookDurationSeconds:  {RecipeFields_CookDurationSeconds},
	model.RecipeField_PrepDurationSeconds:  {RecipeFields_PrepDurationSeconds},
	model.RecipeField_TotalDurationSeconds: {RecipeFields_TotalDurationSeconds},
	model.RecipeField_CookingMethod:        {RecipeFields_CookingMethod},
	model.RecipeField_Categories:           {RecipeFields_Categories},
	model.RecipeField_YieldAmount:          {RecipeFields_YieldAmount},
	model.RecipeField_Cuisines:             {RecipeFields_Cuisines},
})
var RecipeSQLConverter = filter.NewSQLConverter(map[string]string{
	"visibility": RecipeFields_VisibilityLevel,
	"permission": RecipeAccessFields_PermissionLevel,
	"state":      RecipeAccessFields_State,
}, true)

// Recipe -
type Recipe struct {
	RecipeId             int64 `gorm:"primaryKey;bigint;not null;<-:false"`
	Title                string
	Description          string
	Directions           []byte `gorm:"type:jsonb"`
	ImageURI             string
	IngredientGroups     []byte                `gorm:"type:jsonb"`
	VisibilityLevel      types.VisibilityLevel `gorm:"not null;default:1"`
	Citation             string                `gorm:"type:varchar(512)"`
	CookDuration         int64                 `gorm:"column:cook_duration_nanos;type:bigint"` // nanoseconds
	CookingMethod        string                `gorm:"type:varchar(64)"`
	Categories           []byte                `gorm:"type:jsonb"`
	YieldAmount          string                `gorm:"type:varchar(64)"`
	Cuisines             []byte                `gorm:"type:jsonb"`
	CreateTime           time.Time             `gorm:"column:create_time;autoCreateTime"`
	UpdateTime           time.Time             `gorm:"column:update_time;autoUpdateTime"`
	PrepDurationSeconds  int64                 `gorm:"column:prep_duration_seconds;type:bigint"`
	CookDurationSeconds  int64                 `gorm:"column:cook_duration_seconds;type:bigint"`
	TotalDurationSeconds int64                 `gorm:"column:total_duration_seconds;type:bigint"`

	// RecipeAccess data
	RecipeAccessId  int64                 `gorm:"->;-:migration"` // only used for read from a join
	PermissionLevel types.PermissionLevel `gorm:"->;-:migration"` // only used for read from a join
	State           types.AccessState     `gorm:"->;-:migration"` // only used for read from a join
}

// TableName -
func (Recipe) TableName() string {
	return "recipe"
}
