// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	casebytes "github.com/gofiber/utils/v2/bytes"
)

// ToLowerBytes converts an ASCII byte slice to lower-case in-place.
//
// Deprecated: use [github.com/gofiber/utils/v2/bytes.UnsafeToLower] instead.
//
//go:fix inline
func ToLowerBytes(b []byte) []byte {
	return casebytes.UnsafeToLower(b)
}

// ToUpperBytes converts an ASCII byte slice to upper-case in-place.
//
// Deprecated: use [github.com/gofiber/utils/v2/bytes.UnsafeToUpper] instead.
//
//go:fix inline
func ToUpperBytes(b []byte) []byte {
	return casebytes.UnsafeToUpper(b)
}

// AddTrailingSlashBytes returns b if it ends with '/', else a new writable copy with '/' appended.
func AddTrailingSlashBytes(b []byte) []byte {
	n := len(b)
	if n > 0 && b[n-1] == '/' {
		return b
	}
	if n == 0 {
		return []byte{'/'} // cheaper than make and copy
	}
	// Not AddTrailingSlashString: its "/" is read-only.
	out := make([]byte, n+1)
	copy(out, b)
	out[n] = '/'
	return out
}
