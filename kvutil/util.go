// Copyright (c) 2025 Visvasity LLC

package kvutil

// PrefixRange returns the begin and end keys that cover all keys with a given prefix.
func PrefixRange(dir string) (begin string, end string) {
	begin = dir
	end = prefixEnd(dir)
	return begin, end
}

func prefixEnd(prefix string) string {
	b := []byte(prefix)
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] != 0xff {
			b[i]++
			return string(b[:i+1])
		}
	}
	return "" // End of the entire Keyspace.
}
