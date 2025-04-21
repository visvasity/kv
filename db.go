// Copyright (c) 2025 Visvasity LLC

package kv

import "context"

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

// Database defines an interface for creating transactions and snapshots.
type Database[T Transaction, S Snapshot] interface {
	// NewTransaction creates a new read-write transaction. Returns a non-nil
	// error if the operation fails.
	NewTransaction(ctx context.Context) (T, error)

	// NewSnapshot creates a new read-only snapshot. Returns a non-nil error if
	// the operation fails.
	NewSnapshot(ctx context.Context) (S, error)
}
