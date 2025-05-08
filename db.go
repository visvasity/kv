// Copyright (c) 2025 Visvasity LLC

package kv

import (
	"context"
)

// Reader combines interfaces for reading key-value pairs.
type Reader interface {
	Getter
	Ranger
	Scanner
}

// Writer combines interfaces for modifying key-value pairs.
type Writer interface {
	Setter
	Deleter
}

// ReadWriter combines Reader and Writer interfaces for full key-value pair
// access.
type ReadWriter interface {
	Reader
	Writer
}

// Snapshot represents a read-only view of the database at a specific point in
// time.
type Snapshot interface {
	Reader

	// Discard releases resources associated with the snapshot. Returns a non-nil
	// error if the operation fails.
	Discard(ctx context.Context) error
}

// Transaction represents a read-write transaction with atomic operations.
type Transaction interface {
	ReadWriter

	// Rollback cancels the transaction without checking for conflicts. Returns
	// nil on success or os.ErrClosed if the transaction is already committed or
	// rolled back.
	Rollback(ctx context.Context) error

	// Commit validates all reads and writes for conflicts and atomically applies
	// changes to the key-value store. Returns nil on success or os.ErrClosed if
	// the transaction is already committed or rolled back.
	Commit(ctx context.Context) error
}

// Database interface type defines methods required on all key-value databases.
type Database interface {
	NewTransaction(context.Context) (Transaction, error)
	NewSnapshot(context.Context) (Snapshot, error)
}

type NewSnapshotFunc[S Snapshot] func(context.Context) (S, error)
type NewTransactionFunc[T Transaction] func(ctx context.Context) (T, error)

type anyDatabase[T Transaction, S Snapshot] struct {
	newTx   NewTransactionFunc[T]
	newSnap NewSnapshotFunc[S]
}

func (v *anyDatabase[T, S]) NewTransaction(ctx context.Context) (Transaction, error) {
	return v.newTx(ctx)
}

func (v *anyDatabase[T, S]) NewSnapshot(ctx context.Context) (Snapshot, error) {
	return v.newSnap(ctx)
}

// DatabaseFrom is a helper function that can create uniform Database interface
// objects for different database implementations each with their own concrete
// types for Transactions and Snapshots.
func DatabaseFrom[T Transaction, S Snapshot](newTx NewTransactionFunc[T], newSnap NewSnapshotFunc[S]) Database {
	return &anyDatabase[T, S]{
		newTx:   newTx,
		newSnap: newSnap,
	}
}
