package bytes

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
)

// Sub-word inputs convert here, not in caseconv: at these lengths the call costs
// more than the table lookups. Delegating them measured +47%/+43% on unchanged
// three-byte input and +65%/+60% for the in-place forms (benchstat n=15,
// interleaved, Go 1.26, linux/amd64; the http-get cases reproduce it).

// ToLower converts an ASCII byte slice to lower-case without modifying the input.
// If b has no upper-case letter it is returned itself, not a copy, so a write to
// the result then writes to b.
func ToLower(b []byte) []byte {
	n := len(b)
	if n < caseconv.WordLen {
		table := caseconv.ToLowerTable
		for i := 0; i < n; i++ {
			c := b[i]
			low := table[c]
			if low != c {
				dst := make([]byte, n)
				copy(dst, b[:i])
				dst[i] = low
				for i++; i < n; i++ {
					dst[i] = table[b[i]]
				}
				return dst
			}
		}
		return b
	}

	i := caseconv.FirstUpperIndex(b)
	if i < 0 {
		return b
	}
	return caseconv.ToLowerFrom(b, i)
}

// ToUpper converts an ASCII byte slice to upper-case without modifying the input.
// If b has no lower-case letter it is returned itself, not a copy, so a write to
// the result then writes to b.
func ToUpper(b []byte) []byte {
	n := len(b)
	if n < caseconv.WordLen {
		table := caseconv.ToUpperTable
		for i := 0; i < n; i++ {
			c := b[i]
			up := table[c]
			if up != c {
				dst := make([]byte, n)
				copy(dst, b[:i])
				dst[i] = up
				for i++; i < n; i++ {
					dst[i] = table[b[i]]
				}
				return dst
			}
		}
		return b
	}

	i := caseconv.FirstLowerIndex(b)
	if i < 0 {
		return b
	}
	return caseconv.ToUpperFrom(b, i)
}

// UnsafeToLower converts an ASCII byte slice to lower-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToLower(b []byte) []byte {
	if len(b) < caseconv.WordLen {
		table := caseconv.ToLowerTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToLowerInPlace(b)
	return b
}

// UnsafeToUpper converts an ASCII byte slice to upper-case in-place.
// The passed slice content is modified and the same slice is returned.
func UnsafeToUpper(b []byte) []byte {
	if len(b) < caseconv.WordLen {
		table := caseconv.ToUpperTable
		for i := range b {
			b[i] = table[b[i]]
		}
		return b
	}
	caseconv.ToUpperInPlace(b)
	return b
}
