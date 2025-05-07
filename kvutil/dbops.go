// Copyright (c) 2025 Visvasity LLC

package kvutil

import (
	"context"

	"github.com/visvasity/kv"
)

type NewSnapshotFunc[S kv.Snapshot] func(context.Context) (S, error)

// WithReader executes the work function within a temporary snapshot, ensuring
// read-only access to the key-value store. The snapshot is discarded after the
// function completes, regardless of the outcome.
//
// Returns an error if creating the snapshot or executing the function
// fails. The context controls cancellation and timeouts.
func WithReader[S kv.Snapshot](ctx context.Context, snapf NewSnapshotFunc[S], work func(context.Context, kv.Reader) error) error {
	snap, err := snapf(ctx)
	if err != nil {
		return err
	}
	defer snap.Discard(ctx)
	return work(ctx, snap)
}

type NewTransactionFunc[T kv.Transaction] func(ctx context.Context) (T, error)

// WithReadWriter executes the work function within a temporary transaction,
// providing read-write access to the key-value store. The transaction is
// committed if the function returns nil; otherwise, it is rolled back.
//
// Returns an error if creating the transaction, executing the function, or
// committing the transaction fails. The context controls cancellation and
// timeouts.
func WithReadWriter[T kv.Transaction](ctx context.Context, txf NewTransactionFunc[T], work func(context.Context, kv.ReadWriter) error) error {
	tx, err := txf(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := work(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
