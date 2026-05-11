package repository

import (
	"context"

	"herbst-server/db"
)

type entTransactionRunner struct {
	client *db.Client
}

func NewEntTransactionRunner(client *db.Client) TransactionRunner {
	return &entTransactionRunner{client: client}
}

func (r *entTransactionRunner) WithTx(ctx context.Context, fn func(tx *db.Tx) error) error {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}