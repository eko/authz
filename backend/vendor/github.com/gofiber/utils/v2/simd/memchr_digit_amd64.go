// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// memchrDigitAVX2 is implemented in memchr_digit_amd64.s. It finds bytes in
// ['0','9'] with one biased VPADDB and one signed VPCMPGTB per vector, four
// 32-byte vectors per iteration.
//
//go:noescape
func memchrDigitAVX2(haystack []byte) int

// memchrDigitImpl dispatches to the AVX2 kernel for inputs large enough to
// amortize the vector setup. The caller guarantees a non-empty haystack.
func memchrDigitImpl(haystack []byte) int {
	if hasAVX2 && len(haystack) >= MinLen {
		return memchrDigitAVX2(haystack)
	}
	return memchrDigitGeneric(haystack)
}
