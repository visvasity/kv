// Copyright (c) 2025 Visvasity LLC

package kv

import (
	"context"
)

// WithReader executes the provided function within a temporary snapshot,
// ensuring read-only access to the key-value store. The snapshot is discarded
// after the function completes, regardless of the outcome.
//
// Returns an error if creating the snapshot or executing the function
// fails. The context controls cancellation and timeouts.
func WithReader[T Transaction, S Snapshot](ctx context.Context, db Database[T, S], f func(context.Context, Reader) error) error {
	snap, err := db.NewSnapshot(ctx)
	if err != nil {
		return err
	}
	defer snap.Discard(ctx)
	return f(ctx, snap)
}

// WithReadWriter executes the provided function within a temporary
// transaction, providing read-write access to the key-value store. The
// transaction is committed if the function returns nil; otherwise, it is
// rolled back.
//
// Returns an error if creating the transaction, executing the function, or
// committing the transaction fails. The context controls cancellation and
// timeouts.
func WithReadWriter[T Transaction, S Snapshot](ctx context.Context, db Database[T, S], f func(context.Context, ReadWriter) error) error {
	tx, err := db.NewTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := f(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
