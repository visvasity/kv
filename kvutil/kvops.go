// Copyright (c) 2025 Visvasity LLC

package kvutil

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"iter"

	"github.com/visvasity/kv"
)

// GetGob decodes the Gob encoded bytes at the key and returns as an
// object. Input value parameter must be of a pointer type.
func GetGob(ctx context.Context, g kv.Getter, key string, value any) error {
	v, err := g.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("could not Get from %q: %w", key, err)
	}
	if err := gob.NewDecoder(v).Decode(value); err != nil {
		return fmt.Errorf("could not gob-decode value at key %q: %w", key, err)
	}
	return nil
}

// SetGob creates or updates the value at the key to Gob encoded bytes of the
// input value.
func SetGob(ctx context.Context, s kv.Setter, key string, value any) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return err
	}
	return s.Set(ctx, key, &buf)
}

// AscendGob iterates over all values in the input range unmarshaling the Gob
// encoded bytes into a value of type T. Returned value through the iterator is
// overwritten when the next key-value pair is visited.
func AscendGob[T any](ctx context.Context, r kv.Ranger, beg, end string, errp *error) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		var gv, zero T
		for k, v := range r.Ascend(ctx, beg, end, errp) {
			gv = zero
			if err := gob.NewDecoder(v).Decode(&gv); err != nil {
				*errp = err
				return
			}
			if !yield(k, &gv) {
				return
			}
		}
	}
}

// DescendGob iterates over all values in the input range unmarshaling the Gob
// encoded bytes into a value of type T. Returned value through the iterator is
// overwritten when the next key-value pair is visited.
func DescendGob[T any](ctx context.Context, r kv.Ranger, beg, end string, errp *error) iter.Seq2[string, *T] {
	return func(yield func(string, *T) bool) {
		var gv, zero T
		for k, v := range r.Descend(ctx, beg, end, errp) {
			gv = zero
			if err := gob.NewDecoder(v).Decode(&gv); err != nil {
				*errp = err
				return
			}
			if !yield(k, &gv) {
				return
			}
		}
	}
}
