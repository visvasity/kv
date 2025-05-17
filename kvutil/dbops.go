// Copyright (c) 2025 Visvasity LLC

package kvutil

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"math"
	"os"

	"github.com/visvasity/kv"
)

// WithReader executes the work function within a temporary snapshot, ensuring
// read-only access to the key-value store. The snapshot is discarded after the
// function completes, regardless of the outcome.
//
// Returns an error if creating the snapshot or executing the function
// fails. The context controls cancellation and timeouts.
func WithReader(ctx context.Context, db kv.Database, work func(context.Context, kv.Reader) error) error {
	snap, err := db.NewSnapshot(ctx)
	if err != nil {
		return err
	}
	defer snap.Discard(ctx)
	return work(ctx, snap)
}

// WithReadWriter executes the work function within a temporary transaction,
// providing read-write access to the key-value store. The transaction is
// committed if the function returns nil; otherwise, it is rolled back.
//
// Returns an error if creating the transaction, executing the function, or
// committing the transaction fails. The context controls cancellation and
// timeouts.
func WithReadWriter(ctx context.Context, db kv.Database, work func(context.Context, kv.ReadWriter) error) error {
	tx, err := db.NewTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := work(ctx, tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Backup saves database content into the writer. Written data will be a
// consistent snapshot of the database.
func Backup(ctx context.Context, db kv.Database, w io.Writer) error {
	enc := gob.NewEncoder(w)

	snap, err := db.NewSnapshot(ctx)
	if err != nil {
		return err
	}
	defer snap.Discard(ctx)

	for k, v := range snap.Ascend(ctx, "", "", &err) {
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, v); err != nil {
			return err
		}
		// Compute crc64 for the key-value pair.
		hash := crc64.New(crc64.MakeTable(crc64.ISO))
		if _, err := io.WriteString(hash, k); err != nil {
			return err
		}
		if _, err := hash.Write(buffer.Bytes()); err != nil {
			return err
		}
		checksum := hash.Sum64()

		if err := enc.Encode(k); err != nil {
			return err
		}
		if err := enc.Encode(buffer.Bytes()); err != nil {
			return err
		}
		if err := enc.Encode(checksum); err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}

func nextBackupItem(ctx context.Context, dec *gob.Decoder) (string, io.Reader, error) {
	var key string
	var value []byte
	var checksum uint64
	if err := dec.Decode(&key); err != nil {
		return "", nil, err
	}
	if err := dec.Decode(&value); err != nil {
		return "", nil, err
	}
	if err := dec.Decode(&checksum); err != nil {
		return "", nil, err
	}
	// Compute crc64 for the key-value pair.
	hash := crc64.New(crc64.MakeTable(crc64.ISO))
	if _, err := io.WriteString(hash, key); err != nil {
		return "", nil, err
	}
	if _, err := hash.Write(value); err != nil {
		return "", nil, err
	}
	if checksum != hash.Sum64() {
		return "", nil, fmt.Errorf("checksum error detected")
	}
	return key, bytes.NewReader(value), nil
}

// ValidateBackup scans the database backup for checksum errors.
func ValidateBackup(ctx context.Context, r io.Reader) error {
	dec := gob.NewDecoder(r)
	for {
		if _, _, err := nextBackupItem(ctx, dec); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}
	return nil
}

// Clear deletes all key-value pairs in the database. When maxPerTx is
// non-zero, operation is split into multiple transactions, with at most
// maxPerTx deletions in one transaction.
func Clear(ctx context.Context, db kv.Database, maxPerTx int64) error {
	if maxPerTx == 0 {
		maxPerTx = math.MaxInt64
	}
	if maxPerTx < 0 {
		return os.ErrInvalid
	}

	for done := false; !done; {
		err := func() error {
			tx, err := db.NewTransaction(ctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)

			n := int64(0)
			for k := range tx.Ascend(ctx, "", "", &err) {
				if err := tx.Delete(ctx, k); err != nil {
					return err
				}
				n++
				if n == maxPerTx {
					break
				}
			}
			if err != nil {
				return err
			}

			if n == 0 {
				done = true
				return nil
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

// Restore updates the database with key-value content from the reader. When
// maxPerTx is non-zero, restore will happen in multiple transactions with at
// most maxPerTx updates in one transaction.
func Restore(ctx context.Context, db kv.Database, r io.Reader, maxPerTx int64) error {
	if maxPerTx == 0 {
		maxPerTx = math.MaxInt64
	}
	if maxPerTx < 0 {
		return os.ErrInvalid
	}

	dec := gob.NewDecoder(r)
	for done := false; !done; {
		err := func() error {
			tx, err := db.NewTransaction(ctx)
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)

			for i := int64(0); i < maxPerTx; i++ {
				key, value, err := nextBackupItem(ctx, dec)
				if err != nil {
					if errors.Is(err, io.EOF) {
						done = true
						break
					}
					return err
				}
				if err := tx.Set(ctx, key, value); err != nil {
					return err
				}
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}
