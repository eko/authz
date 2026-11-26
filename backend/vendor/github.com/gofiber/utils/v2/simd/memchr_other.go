// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build !amd64

package simd

// The portable SWAR implementations are the only ones off amd64.

func memchr2Impl(haystack []byte, needle1, needle2 byte) int {
	return memchr2Generic(haystack, needle1, needle2)
}

func memchr3Impl(haystack []byte, needle1, needle2, needle3 byte) int {
	return memchr3Generic(haystack, needle1, needle2, needle3)
}

func memchrPairImpl(haystack []byte, byte1, byte2 byte, offset int) int {
	return memchrPairGeneric(haystack, byte1, byte2, offset)
}
