// Copyright (c) 2025 Visvasity LLC

package kvutil

import "testing"

func TestPrefixRange(t *testing.T) {
	beg, end := PrefixRange("abcd")
	t.Logf("begin=%q end=%q", beg, end)
}

func TestEndOfKeyspacePrefixRange(t *testing.T) {
	beg, end := PrefixRange("")
	if beg != "" {
		t.Errorf(`wanted beg="" got %q`, beg)
	}
	if end != "" {
		t.Errorf(`wanted end="" got %q`, beg)
	}

	beg, end = PrefixRange("\xff")
	if beg != "\xff" {
		t.Errorf(`wanted beg="\xff" got %q`, beg)
	}
	if end != "" {
		t.Errorf(`wanted end="" got %q`, end)
	}

	beg, end = PrefixRange("\x00\xff")
	if beg != "\x00\xff" {
		t.Errorf(`wanted beg="\x00\xff" got %q`, beg)
	}
	if end != "\x01" {
		t.Errorf(`wanted end="\x01" got %q`, end)
	}
}
