// Copyright (c) 2025 Visvasity LLC

package kvutil

import (
	"fmt"
)

// PrefixRange returns the begin and end keys that cover all keys with a given prefix.
func PrefixRange(dir string) (begin string, end string) {
	n := len(dir)
	if n == 0 {
		return "", ""
	}
	begin = dir
	end = dir[:n-1] + fmt.Sprintf("%c", dir[n-1]+1)
	return begin, end
}
