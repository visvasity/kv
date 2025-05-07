// Copyright (c) 2025 Visvasity LLC

package kvutil

import "testing"

func TestPrefixRange(t *testing.T) {
	beg, end := PrefixRange("abcd")
	t.Logf("begin=%q end=%q", beg, end)
}
