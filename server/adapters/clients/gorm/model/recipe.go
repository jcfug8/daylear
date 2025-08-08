package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

const (
	RecipeTable = "recipe"
)

const (
	RecipeFields_RecipeId             = "recipe_id"
	RecipeFields_Title                = "title"
	RecipeFields_Description          = "description"
	RecipeFields_Directions           = "directions"
	RecipeFields_ImageURI             = "image_uri"
	RecipeFields_IngredientGroups     = "ingredient_groups"
	RecipeFields_VisibilityLevel      = "visibility_level"
	RecipeFields_Citation             = "citation"
	RecipeFields_CookDurationSeconds  = "cook_duration_seconds"
	RecipeFields_PrepDurationSeconds  = "prep_duration_seconds"
	RecipeFields_TotalDurationSeconds = "total_duration_seconds"
	RecipeFields_CookingMethod        = "cooking_method"
	RecipeFields_Categories           = "categories"
	RecipeFields_YieldAmount          = "yield_amount"
	RecipeFields_Cuisines             = "cuisines"
	RecipeFields_CreateTime           = "create_time"
	RecipeFields_UpdateTime           = "update_time"
)

var RecipeFieldMasker = fieldmask.NewSQLFieldMasker(Recipe{}, map[string][]fieldmask.Field{
	model.RecipeField_Id:                   {{Name: RecipeFields_RecipeId, Table: RecipeTable}},
	model.RecipeField_Title:                {{Name: RecipeFields_Title, Table: RecipeTable, Updatable: true}},
	model.RecipeField_Description:          {{Name: RecipeFields_Description, Table: RecipeTable, Updatable: true}},
	model.RecipeField_Directions:           {{Name: RecipeFields_Directions, Table: RecipeTable, Updatable: true}},
	model.RecipeField_ImageURI:             {{Name: RecipeFields_ImageURI, Table: RecipeTable, Updatable: true}},
	model.RecipeField_IngredientGroups:     {{Name: RecipeFields_IngredientGroups, Table: RecipeTable, Updatable: true}},
	model.RecipeField_VisibilityLevel:      {{Name: RecipeFields_VisibilityLevel, Table: RecipeTable, Updatable: true}},
	model.RecipeField_Citation:             {{Name: RecipeFields_Citation, Table: RecipeTable, Updatable: true}},
	model.RecipeField_CookDurationSeconds:  {{Name: RecipeFields_CookDurationSeconds, Table: RecipeTable, Updatable: true}},
	model.RecipeField_PrepDurationSeconds:  {{Name: RecipeFields_PrepDurationSeconds, Table: RecipeTable, Updatable: true}},
	model.RecipeField_TotalDurationSeconds: {{Name: RecipeFields_TotalDurationSeconds, Table: RecipeTable, Updatable: true}},
	model.RecipeField_CookingMethod:        {{Name: RecipeFields_CookingMethod, Table: RecipeTable, Updatable: true}},
	model.RecipeField_Categories:           {{Name: RecipeFields_Categories, Table: RecipeTable, Updatable: true}},
	model.RecipeField_YieldAmount:          {{Name: RecipeFields_YieldAmount, Table: RecipeTable, Updatable: true}},
	model.RecipeField_Cuisines:             {{Name: RecipeFields_Cuisines, Table: RecipeTable, Updatable: true}},
	model.RecipeField_CreateTime:           {{Name: RecipeFields_CreateTime, Table: RecipeTable}},
	model.RecipeField_UpdateTime:           {{Name: RecipeFields_UpdateTime, Table: RecipeTable}},

	model.RecipeField_RecipeAccess: {
		{Name: RecipeAccessFields_RecipeAccessId, Table: RecipeAccessTable},
		{Name: RecipeAccessFields_PermissionLevel, Table: RecipeAccessTable},
		{Name: RecipeAccessFields_State, Table: RecipeAccessTable},
	},
})
var RecipeSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"visibility": {Name: RecipeFields_VisibilityLevel, Table: RecipeTable},
	"permission": {Name: RecipeAccessFields_PermissionLevel, Table: RecipeAccessTable},
	"state":      {Name: RecipeAccessFields_State, Table: RecipeAccessTable},
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
	return RecipeTable
}
