// Copyright (c) 2025 Visvasity LLC

// Package kv provides a uniform API for key-value databases supporting
// snapshots and transactions. The API is designed to be minimal, unambiguous,
// and well-defined, enabling multiple database backends to implement adapters
// for this interface.
//
// Keys are represented as strings. The empty string is not a valid key, so
// that it can be used to denote the beginning and/or end of key ranges.
//
// Values are represented as io.Reader objects, so that database
// implementations can stream large values or may provide additional metadata
// with the values.
package kv
