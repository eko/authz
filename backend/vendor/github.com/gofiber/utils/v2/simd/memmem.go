// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"
)

// memmemMinPairedHaystack is the haystack size from which Memmem routes
// needles eligible for the paired prefilter through it instead of
// bytes.Index. bytes.Index brute-forces haystacks of up to 64 bytes with
// a vector kernel the prefilter cannot beat — Benchmark_Memmem_Prefilter
// shows the two level at 64B before SelectRareBytes is even paid — and
// from 65 bytes switches to a loop anchored on the needle's first byte,
// which is where a scan anchored on the needle's two rarest bytes at
// their exact distance pulls ahead. Benchmark_Memmem shows the resulting
// end-to-end routing; its 96B row is the band this constant opens up.
const memmemMinPairedHaystack = 65

// memmemMinSingleHaystack is the same routing point for needles that take
// the single rare-byte scan. That scan is far less selective (its false
// candidates need only one byte to match) and its verification is
// longer, so it needs more haystack to amortize the rare-byte selection
// and the candidate loop; 128 keeps a conservative margin over the band
// where Benchmark_Memmem_Prefilter's single rows trail bytes.Index.
const memmemMinSingleHaystack = 128

// memmemMaxPairNeedle is the needle length up to which the paired-byte
// prefilter beats a single rare-byte scan: short needles produce frequent
// single-byte false positives, and verifying them costs more than the second
// byte check saves.
const memmemMaxPairNeedle = 6

// memmemMaxMisses caps how many failed candidate verifications the
// prefilter loops tolerate before handing the search to bytes.Index. A
// pathological input (rare byte everywhere, long almost-matching needle)
// would otherwise cost O(len(haystack) * len(needle)); the cap bounds the
// wasted work at memmemMaxMisses needle comparisons, keeping the worst case
// at bytes.Index's O(n+m) plus a constant.
const memmemMaxMisses = 16

// Memmem returns the index of the first occurrence of needle in haystack,
// or -1 if needle is not present. An empty needle matches at index 0,
// mirroring bytes.Index.
//
// On amd64 with AVX2 it prefilters with the needle's rarest bytes (per
// ByteRank) and verifies only the candidate positions, which typically
// beats bytes.Index by a wide margin on large inputs: needles of 2-6
// bytes containing two distinct byte values are scanned, on haystacks of
// 65+ bytes, for their two rarest values at the exact relative distance,
// while longer or single-valued needles are scanned, on haystacks of
// 128+ bytes, for the single rarest byte. In all other cases it delegates
// to bytes.Index directly, and the prefilter itself falls back to
// bytes.Index after a bounded number of failed candidate verifications,
// so the worst case stays within a constant of the stdlib's O(n+m).
func Memmem(haystack, needle []byte) int {
	needleLen := len(needle)
	if needleLen == 0 {
		return 0
	}
	if len(haystack) < needleLen {
		return -1
	}
	if needleLen == 1 {
		return bytes.IndexByte(haystack, needle[0])
	}
	if !hasAVX2 || len(haystack) < memmemMinPairedHaystack {
		return bytes.Index(haystack, needle)
	}

	var info RareByteInfo
	pairable := needleLen <= memmemMaxPairNeedle
	if pairable {
		info = SelectRareBytes(needle)
		// Distinct byte values imply distinct indices, so Byte1 != Byte2
		// alone guarantees the offset MemchrPair needs is >= 1.
		if info.Byte1 != info.Byte2 {
			return memmemPaired(haystack, needle, info)
		}
	}
	if len(haystack) < memmemMinSingleHaystack {
		return bytes.Index(haystack, needle)
	}
	if !pairable {
		info = SelectRareBytes(needle)
	}
	return memmemSingle(haystack, needle, info.Byte1, info.Index1)
}

// memmemPaired scans for the needle's two rarest bytes at their exact
// relative distance with MemchrPair and verifies each candidate. False
// candidates need both bytes at the right distance, so verification is
// rare; if it still misses too often, memmemFallback finishes the search.
func memmemPaired(haystack, needle []byte, info RareByteInfo) int {
	byte1, idx1 := info.Byte1, info.Index1
	byte2, idx2 := info.Byte2, info.Index2
	if idx1 > idx2 {
		byte1, byte2 = byte2, byte1
		idx1, idx2 = idx2, idx1
	}
	offset := idx2 - idx1

	misses := 0
	for start := 0; ; {
		cand := MemchrPair(haystack[start:], byte1, byte2, offset)
		if cand == -1 {
			return -1
		}
		// cand is the position of byte1; the needle would start idx1 earlier.
		pos := start + cand - idx1
		if pos >= 0 && pos+len(needle) <= len(haystack) &&
			bytes.Equal(haystack[pos:pos+len(needle)], needle) {
			return pos
		}
		start += cand + 1
		if misses++; misses >= memmemMaxMisses {
			return memmemFallback(haystack, needle, start, idx1)
		}
	}
}

// memmemSingle scans for the needle's rarest byte and verifies each
// candidate; for longer needles the rare byte alone is selective enough
// that the paired scan's second comparison stops paying for itself. Like
// memmemPaired it hands persistent misses to memmemFallback.
func memmemSingle(haystack, needle []byte, rareByte byte, rareIdx int) int {
	misses := 0
	for start := 0; ; {
		cand := bytes.IndexByte(haystack[start:], rareByte)
		if cand == -1 {
			return -1
		}
		pos := start + cand - rareIdx
		if pos >= 0 && pos+len(needle) <= len(haystack) &&
			bytes.Equal(haystack[pos:pos+len(needle)], needle) {
			return pos
		}
		start += cand + 1
		if misses++; misses >= memmemMaxMisses {
			return memmemFallback(haystack, needle, start, rareIdx)
		}
	}
}

// memmemFallback finishes a prefiltered search with bytes.Index once the
// candidate loops exhaust their miss budget. Every anchor before start has
// been checked, and any remaining occurrence carries its anchor byte at
// anchorIdx, so it starts at or after start-anchorIdx; searching from there
// cannot skip a match or return one out of order. The subtraction is in
// fact defensive: start-1 is always an anchor at handover and anchorIdx is
// the anchor byte's first occurrence in the needle, so no remaining
// occurrence can start below start itself — resuming at start-anchorIdx
// keeps correctness independent of that selection detail.
func memmemFallback(haystack, needle []byte, start, anchorIdx int) int {
	from := max(start-anchorIdx, 0)
	if pos := bytes.Index(haystack[from:], needle); pos >= 0 {
		return from + pos
	}
	return -1
}
