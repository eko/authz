// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// Assembly kernels implemented in memchr_class_amd64.s. Class membership
// is computed with a pair of VPSHUFB nibble lookups instead of one
// clamp-and-compare per range, and four 32-byte vectors are classified per
// iteration.
//
//go:noescape
func memchrWordAVX2(haystack []byte) int

//go:noescape
func memchrNotWordAVX2(haystack []byte) int

// memchrWordImpl dispatches to the AVX2 kernel for inputs large enough to
// amortize the vector setup. The caller guarantees a non-empty haystack.
func memchrWordImpl(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrWordAVX2(haystack)
	}
	return memchrWordGeneric(haystack)
}

// memchrNotWordImpl dispatches like memchrWordImpl.
func memchrNotWordImpl(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrNotWordAVX2(haystack)
	}
	return memchrNotWordGeneric(haystack)
}
