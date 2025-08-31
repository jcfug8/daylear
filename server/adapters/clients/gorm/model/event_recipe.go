package model

import (
	"time"

	"github.com/jcfug8/daylear/server/core/fieldmask"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/filter"
)

const (
	EventRecipeTable = "event_recipe"
)

const (
	EventRecipeColumn_EventRecipeId = "event_recipe_id"
	EventRecipeColumn_RecipeId      = "recipe_id"
	EventRecipeColumn_EventId       = "event_id"
	EventRecipeColumn_CreateTime    = "create_time"
)

var EventRecipeFieldMasker = fieldmask.NewSQLFieldMasker(EventRecipe{}, map[string][]fieldmask.Field{
	cmodel.EventRecipeField_Parent:        {{Name: EventRecipeColumn_EventId, Table: EventRecipeTable}},
	cmodel.EventRecipeField_EventRecipeId: {{Name: EventRecipeColumn_EventRecipeId, Table: EventRecipeTable}},
	cmodel.EventRecipeField_RecipeId:      {{Name: EventRecipeColumn_RecipeId, Table: EventRecipeTable}},
	cmodel.EventRecipeField_CreateTime:    {{Name: EventRecipeColumn_CreateTime, Table: EventRecipeTable}},
})

var EventRecipeSQLConverter = filter.NewSQLConverter(map[string]filter.Field{
	"recipe_id": {Name: EventRecipeColumn_RecipeId, Table: EventRecipeTable},
	"event_id":  {Name: EventRecipeColumn_EventId, Table: EventRecipeTable},
}, true)

// EventRecipe is the GORM model for an event recipe.
type EventRecipe struct {
	EventRecipeId int64     `gorm:"primaryKey;column:event_recipe_id;autoIncrement;<-:false"`
	RecipeId      int64     `gorm:"column:recipe_id;not null"`
	EventId       int64     `gorm:"column:event_id;not null"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime"`
}

// TableName sets the table name for the EventRecipe model.
func (EventRecipe) TableName() string {
	return EventRecipeTable
}
