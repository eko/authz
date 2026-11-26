package utils

import (
	"bytes"
	"iter"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// SplitTrimSeq iterates over the sep-separated elements of s with ASCII
// whitespace trimmed and empty elements skipped, as HTTP list fields are read
// (RFC 9110 Section 5.6.1); the elements alias s. Quoted strings are not
// interpreted, so a sep inside one still separates.
func SplitTrimSeq[S byteSeq](s S, sep byte) iter.Seq[S] {
	return func(yield func(S) bool) {
		rest := s // traversal state stays local, so the sequence can be ranged over again
		for {
			// CutByte's body, inline: it is too large to inline into the iterator.
			i := bytes.IndexByte(unsafeconv.Bytes(rest), sep)
			elem := rest
			if i >= 0 {
				elem = rest[:i]
			}
			if elem = TrimSpace(elem); len(elem) > 0 && !yield(elem) {
				return
			}
			if i < 0 {
				return
			}
			rest = rest[i+1:]
		}
	}
}
