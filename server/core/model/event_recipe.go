package model

import (
	"time"
)

// EventRecipeFields defines the event recipe fields.
const (
	EventRecipeField_Parent        = "parent"
	EventRecipeField_EventRecipeId = "id"
	EventRecipeField_RecipeId      = "recipe_id"
	EventRecipeField_CreateTime    = "create_time"
)

// EventRecipe represents a connection between a recipe and an event.
type EventRecipe struct {
	// Parent is the parent of the event recipe
	Parent EventRecipeParent
	// EventRecipeId is the unique identifier for the event recipe
	EventRecipeId EventRecipeId
	// RecipeId is the ID of the recipe
	RecipeId RecipeId
	// CreateTime is the time the event recipe was created
	CreateTime time.Time
}

type EventRecipeId struct {
	EventRecipeId int64 `aip_pattern:"key=event_recipe"`
}

type EventRecipeParent struct {
	CalendarId int64 `aip_pattern:"key=calendar"`
	EventId    int64 `aip_pattern:"key=event"`
}
