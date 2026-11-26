// Package simd provides SIMD-accelerated byte searching and validation
// primitives for the Fiber ecosystem: multi-needle scans, character-class
// scans, paired-byte scans, substring search, and ASCII validation.
//
// Acceleration comes in two tiers. The scan kernels — Memchr2, Memchr3,
// MemchrPair, MemchrDigit, MemchrWord, MemchrNotWord, IsASCII,
// FirstNonASCII, and CountNonASCII — dispatch to AVX2 assembly on amd64
// CPUs (for inputs of MinLen+ bytes) and fall back to portable SWAR loops
// built on the swar package everywhere else. The kernels consume four
// 32-byte vectors per iteration. The first-match scans and IsASCII combine
// the per-vector masks, so a 128-byte block costs one mask extraction and
// one branch, and the scans that report a position then re-extract the
// four masks in address order to locate the hit. CountNonASCII is the
// exception: a count cannot be recovered from a combined mask, so its
// block loop accumulates per-lane counts in vector byte lanes and reduces
// them once at the end, and the remaining whole vectors are extracted and
// population-counted (which is why it, unlike the others, gates on POPCNT
// as well as AVX2).
// Memmem prefilters on amd64 — needles of 2-6 bytes containing two
// distinct values through the MemchrPair kernel from 65-byte haystacks,
// longer or single-valued needles through a single rare-byte
// bytes.IndexByte scan from 128 bytes — and delegates to bytes.Index for
// shorter haystacks or without AVX2.
// MemchrInTable/MemchrNotInTable are plain scalar loops on every
// architecture, since an arbitrary 256-entry membership test has no cheap
// vector form. Accelerated reports whether the AVX2 tier is active.
// Single-byte search is deliberately absent: bytes.IndexByte is already
// vector-accelerated by the Go runtime on all major architectures.
//
// Unlike the top-level utils helpers, this package operates on []byte only —
// it is the raw engine underneath the generic helpers. Callers holding a
// string can use the top-level generic functions instead of converting.
//
// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
package simd

// MinLen is the input length from which the AVX2 kernels engage on amd64:
// they consume whole 32-byte vectors, and below one vector their setup cost
// outweighs the win, so shorter inputs take the SWAR fallbacks. Callers
// that route between their own word loops and this package (as the
// top-level utils helpers do) should gate on this constant rather than
// hard-coding 32, so a retune here propagates everywhere.
const MinLen = 32

// Accelerated reports whether the AVX2 assembly kernels are in use on this
// CPU. When false, all functions use portable fallback implementations.
func Accelerated() bool {
	return hasAVX2
}
