package gorm

import (
	"context"
	"errors"
	"strings"

	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/ports/repository"

	"gorm.io/gorm"
)

// Begin begins a transaction
func (c *Client) Begin(context.Context) (repository.TxClient, error) {
	return c.beginTransaction()
}

// Commit commits the transaction
func (c *Client) Commit() error {
	errz := errz.Context("manager.gorm.repository.commit")
	return ErrzError(errz, "", c.db.Commit().Error)
}

// Rollback rolls back the transaction
func (c *Client) Rollback() {
	r := recover()

	err := c.db.Rollback().Error
	if err != nil &&
		!errors.Is(err, gorm.ErrInvalidTransaction) &&
		!strings.Contains(err.Error(),
			"transaction has already been committed or rolled back") {
		c.log.Error().Err(err).Msg("unable to rollback transaction")
	}

	if _, ok := r.(error); ok {
		panic(r)
	}
}

// ----------------------------------------------------------------------------

type txRepository struct {
	*Client
	isNested bool
}

func (c *Client) beginTransaction() (*txRepository, error) {
	errz := errz.Context("manager.gorm.repository.begin_transaction")

	if c.level > 0 {
		return &txRepository{
			Client: &Client{
				db:    c.db,
				level: c.level + 1,
				log:   c.log.With().Int("tx_level", c.level+1).Logger(),
			},
			isNested: true,
		}, nil
	}

	tx := c.db.Begin()
	if err := tx.Error; err != nil {
		return nil, ErrzError(errz, "", err)
	}

	return &txRepository{
		Client: &Client{
			db:    tx,
			level: c.level + 1,
			log:   c.log.With().Int("tx_level", c.level+1).Logger(),
		},
	}, nil
}

// Commit commits the transaction if there is no nested transaction, otherwise
// it's a no-op
func (repo *txRepository) Commit() error {
	if !repo.isNested {
		return repo.Client.Commit()
	}
	return nil
}

// Rollback rolls back the transaction if there is no nested transaction,
// otherwise it's a no-op
func (repo *txRepository) Rollback() {
	if !repo.isNested {
		repo.Client.Rollback()
	}
}
