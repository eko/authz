// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build !amd64

package simd

// The portable implementations are the only ones off amd64.

func memchrWordImpl(haystack []byte) int {
	return memchrWordGeneric(haystack)
}

func memchrNotWordImpl(haystack []byte) int {
	return memchrNotWordGeneric(haystack)
}
