// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// Assembly kernels implemented in memchr_amd64.s. Each consumes four
// 32-byte AVX2 vectors per iteration (OR-ing the per-vector match masks so
// a 128-byte block needs one extraction and one branch), drops to a
// 32-byte loop for the remainder, and finishes the last 1-31 positions
// with an overlapping rescan: memchr2 and memchr3 reload one 32-byte
// vector ending at the buffer end, while memchrPair reloads both of its
// windows — the byte1 side ending at len-offset, the byte2 side at the
// buffer end.
//
//go:noescape
func memchr2AVX2(haystack []byte, needle1, needle2 byte) int

//go:noescape
func memchr3AVX2(haystack []byte, needle1, needle2, needle3 byte) int

//go:noescape
func memchrPairAVX2(haystack []byte, byte1, byte2 byte, offset int) int

// memchr2Impl dispatches to the AVX2 kernel for inputs large enough to amortize
// the vector setup. The caller guarantees a non-empty haystack.
func memchr2Impl(haystack []byte, needle1, needle2 byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchr2AVX2(haystack, needle1, needle2)
	}
	return memchr2Generic(haystack, needle1, needle2)
}

// memchr3Impl dispatches like memchr2Impl.
func memchr3Impl(haystack []byte, needle1, needle2, needle3 byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchr3AVX2(haystack, needle1, needle2, needle3)
	}
	return memchr3Generic(haystack, needle1, needle2, needle3)
}

// memchrPairImpl dispatches to the AVX2 kernel when there is room for at least
// one 32-byte vector of pair positions. The caller guarantees offset >= 1
// and len(haystack) > offset.
func memchrPairImpl(haystack []byte, byte1, byte2 byte, offset int) int {
	if hasAVX2 && len(haystack) >= MinLen+offset {
		return memchrPairAVX2(haystack, byte1, byte2, offset)
	}
	return memchrPairGeneric(haystack, byte1, byte2, offset)
}
