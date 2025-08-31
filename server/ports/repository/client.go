package repository

import (
	"context"
)

// Client defines how to interact with a database client.
type Client interface {
	circleClient
	recipeClient
	userClient
	calendarClient
	eventClient
	eventRecipeClient
	accessKeyClient

	Begin(context.Context) (TxClient, error)
	Migrate() error
}

// TxClient defines how to interact with a database transaction.
type TxClient interface {
	circleClient
	recipeClient
	userClient
	calendarClient
	eventClient
	eventRecipeClient
	accessKeyClient

	Commit() error
	Rollback()
}
