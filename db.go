// Copyright (c) 2025 Visvasity LLC

package kv

import (
	"context"
)

// Reader combines interfaces for reading key-value pairs.
type Reader interface {
	Getter
	Ranger
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
	// changes to the key-value store. Returns nil on success or if the
	// transaction has already committed. Returns non-nil error if transaction
	// commit has failed or was already rolled back.
	//
	// In case of remote databases, this function MUST NOT return a non-nil error
	// ever (eg: timeout) if the transaction was committed, but success response
	// is lost due to network issues/delays. Database client SHOULD internally
	// retry forever (as necessary) to confirm the final status of the
	// transaction.
	Commit(ctx context.Context) error
}

// Database interface type defines methods required on all key-value databases.
type Database interface {
	NewTransaction(context.Context) (Transaction, error)
	NewSnapshot(context.Context) (Snapshot, error)
}

// GenericDatabase interface is similar to the Database interface, but uses
// generic type arguments to represent database specific concrete types for
// NewTransaction and NewSnapshot methods. [DatabaseFrom] function can be used
// to convert a GenericDatabase to a non-generic Database interface.
type GenericDatabase[T Transaction, S Snapshot] interface {
	NewTransaction(context.Context) (T, error)
	NewSnapshot(context.Context) (S, error)
}

// DatabaseFrom is a helper function that can create non-generic Database
// interface object for different database implementations each with their own
// concrete return types for NewTransaction and NewSnapshot methods.
func DatabaseFrom[T Transaction, S Snapshot](db GenericDatabase[T, S]) Database {
	return &anyDatabase[T, S]{
		newTx:   db.NewTransaction,
		newSnap: db.NewSnapshot,
	}
}

type anyDatabase[T Transaction, S Snapshot] struct {
	newTx   func(context.Context) (T, error)
	newSnap func(context.Context) (S, error)
}

func (v *anyDatabase[T, S]) NewTransaction(ctx context.Context) (Transaction, error) {
	return v.newTx(ctx)
}

func (v *anyDatabase[T, S]) NewSnapshot(ctx context.Context) (Snapshot, error) {
	return v.newSnap(ctx)
}
