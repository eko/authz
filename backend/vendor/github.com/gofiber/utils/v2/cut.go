package utils

import (
	"bytes"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// CutByte is strings.Cut/bytes.Cut for a single-byte separator, generic over
// strings and byte slices; a miss returns s, the zero S ("" or nil), and
// false. The parts alias s.
func CutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.IndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}

// LastCutByte is CutByte around the last occurrence of sep.
func LastCutByte[S byteSeq](s S, sep byte) (before, after S, found bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	if i := bytes.LastIndexByte(unsafeconv.Bytes(s), sep); i >= 0 {
		return s[:i], s[i+1:], true
	}
	return s, after, false
}
