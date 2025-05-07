// Copyright (c) 2025 Visvasity LLC

package kvutil

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/visvasity/kv"
)

// GetGob decodes the Gob encoded bytes at the key and returns as an object.
func GetGob[T any](ctx context.Context, g kv.Getter, key string) (*T, error) {
	value, err := g.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("could not Get from %q: %w", key, err)
	}
	gv := new(T)
	if err := gob.NewDecoder(value).Decode(gv); err != nil {
		return nil, fmt.Errorf("could not gob-decode value at key %q: %w", key, err)
	}
	return gv, nil
}

// SetGob creates or updates the value at the key to Gob encoded bytes of the input value.
func SetGob[T any](ctx context.Context, s kv.Setter, key string, value *T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return s.Set(ctx, key, &buf)
}

// PrefixRange returns the begin and end keys that cover all keys with a given prefix.
func PrefixRange(dir string) (begin string, end string) {
	n := len(dir)
	if n == 0 {
		return "", ""
	}
	begin = dir
	end = dir[:n-1] + fmt.Sprintf("%c", dir[n-1]+1)
	return begin, end
}

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
