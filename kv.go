// Copyright (c) 2025 Visvasity LLC

package kv

import (
	"context"
	"io"
	"iter"
)

// Getter defines an interface for retrieving key-value pairs.
type Getter interface {
	// Get retrieves the value for a given key. Returns os.ErrInvalid if the key
	// is empty and os.ErrNotExist if the key is not found or a non-nil error if
	// the operation fails.
	Get(ctx context.Context, key string) (io.Reader, error)
}

// Setter defines an interface for creating or updating key-value pairs.
type Setter interface {
	// Set creates or updates a key-value pair. Returns os.ErrInvalid if the key
	// is empty or value is nil and returns a non-nil error if the operation
	// fails.
	Set(ctx context.Context, key string, value io.Reader) error
}

// Deleter defines an interface for removing key-value pairs.
type Deleter interface {
	// Delete removes a key-value pair. Returns os.ErrInvalid if the key is empty
	// and os.ErrNotExist if the key does not exist or a non-nil error if the
	// operation fails.
	Delete(ctx context.Context, key string) error
}

// Ranger defines an interface for iterating over key-value pairs within a
// specified range.
type Ranger interface {
	// Ascend returns an iterator over key-value pairs in ascending order within
	// the range defined by begin and end.
	//
	//  - If begin and end are both empty, range includes all pairs in the database.
	//  - If begin is empty, range begins at the smallest key.
	//  - If end is empty, range ends after the largest key.
	//  - If both are non-empty, begin must be less than or equal to end, or
	//    os.ErrInvalid is returned.
	//
	// The range includes the begin key and excludes the end key. Errors are
	// stored in errp.
	Ascend(ctx context.Context, beg, end string, errp *error) iter.Seq2[string, io.Reader]

	// Descend is similar to Ascend but iterates in the descending order.
	Descend(ctx context.Context, beg, end string, errp *error) iter.Seq2[string, io.Reader]
}
