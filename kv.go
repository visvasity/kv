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

// Ranger provides ordered iteration over key-value pairs in lexicographical
// order. Both Ascend and Descend define a half-open range [begin, end).
//
// Empty string argument is handled as follows:
//
//   - beg == "" => iteration starts at the smallest key (no lower bound)
//
//   - end == "" => iteration stops at the last largest key (no upper bound)
//
//   - beg == "" && end == "" => the range covers the entire database
//
// When both begin and end are non-empty and begin > end (byte-wise
// comparison), the iterator yields no items and sets *errp = os.ErrInvalid.
//
// Iteration stops at the first error encountered. The errp must not be nil
// otherwise Ascend/Descend operations will panic.
type Ranger interface {
	// Ascend iterates over keys in ascending order within the range [beg, end).
	Ascend(ctx context.Context, beg, end string, errp *error) iter.Seq2[string, io.Reader]

	// Descend iterates over keys in descending order within the range [beg, end).
	Descend(ctx context.Context, beg, end string, errp *error) iter.Seq2[string, io.Reader]
}

// KeyRanger is similar to the Ranger, but does not load/access the value automatically.
type KeyRanger interface {
	// AscendKeys iterates over keys in ascending order within the range [beg, end).
	AscendKeys(ctx context.Context, beg, end string, errp *error) iter.Seq[string]

	// DescendKeys iterates over keys in descending order within the range [beg, end).
	DescendKeys(ctx context.Context, beg, end string, errp *error) iter.Seq[string]
}
