// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"bytes"

	"github.com/gofiber/utils/v2/swar"
)

// The first-match scans below rely on swar.ZeroLanes having false positives
// only above the first true match, so the first set lane is always exact.

// Memchr2 returns the index of the first occurrence in haystack of either
// needle1 or needle2, or -1 if neither is present. There is no single-needle
// Memchr; use bytes.IndexByte.
func Memchr2(haystack []byte, needle1, needle2 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr2Impl(haystack, needle1, needle2)
}

// Memchr3 returns the index of the first occurrence in haystack of needle1,
// needle2, or needle3, or -1 if none is present.
func Memchr3(haystack []byte, needle1, needle2, needle3 byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchr3Impl(haystack, needle1, needle2, needle3)
}

// MemchrPair returns the first index i such that haystack[i] == byte1 and
// haystack[i+offset] == byte2, or -1 if there is none or offset is negative.
func MemchrPair(haystack []byte, byte1, byte2 byte, offset int) int {
	if offset < 0 || len(haystack) <= offset {
		return -1
	}
	if offset == 0 {
		if byte1 != byte2 {
			return -1
		}
		return bytes.IndexByte(haystack, byte1)
	}
	return memchrPairImpl(haystack, byte1, byte2, offset)
}

func memchr2Generic(haystack []byte, needle1, needle2 byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(needle1)
	bc2 := swar.Broadcast(needle2)
	i := 0
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		// Pinned window: see the note in ascii.go.
		win := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		w0, w1 := swar.Load8(win, 0), swar.Load8(win, swar.WordLen)
		m0 := swar.ZeroLanes(w0^bc1) | swar.ZeroLanes(w0^bc2)
		m1 := swar.ZeroLanes(w1^bc1) | swar.ZeroLanes(w1^bc2)
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		w := swar.Load8(haystack, i)
		if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	// Overlapping final word: scanned lanes are non-matching, so a hit is new.
	w := swar.Load8(haystack, n-swar.WordLen)
	if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

func memchr3Generic(haystack []byte, needle1, needle2, needle3 byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if haystack[i] == needle1 || haystack[i] == needle2 || haystack[i] == needle3 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(needle1)
	bc2 := swar.Broadcast(needle2)
	bc3 := swar.Broadcast(needle3)
	i := 0
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		win := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		w0, w1 := swar.Load8(win, 0), swar.Load8(win, swar.WordLen)
		m0 := swar.ZeroLanes(w0^bc1) | swar.ZeroLanes(w0^bc2) | swar.ZeroLanes(w0^bc3)
		m1 := swar.ZeroLanes(w1^bc1) | swar.ZeroLanes(w1^bc2) | swar.ZeroLanes(w1^bc3)
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		w := swar.Load8(haystack, i)
		if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2) | swar.ZeroLanes(w^bc3); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	w := swar.Load8(haystack, n-swar.WordLen)
	if m := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2) | swar.ZeroLanes(w^bc3); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// memchrPairGeneric: the caller guarantees offset >= 1 and len(haystack) >
// offset. ANDing the two ZeroLanes masks never drops a real pair but loses
// lane exactness (a flagged lane may be a false candidate below the first
// pair), so candidates are re-verified scalar. The final window overlaps
// already-checked positions, whose candidates just fail verification again.
func memchrPairGeneric(haystack []byte, byte1, byte2 byte, offset int) int {
	n := len(haystack)
	if n-offset < swar.WordLen {
		for i := 0; i < n-offset; i++ {
			if haystack[i] == byte1 && haystack[i+offset] == byte2 {
				return i
			}
		}
		return -1
	}

	bc1 := swar.Broadcast(byte1)
	bc2 := swar.Broadcast(byte2)
	last := n - offset // candidate positions are 0..last-1
	i := 0
	for ; i <= last-2*swar.WordLen; i += 2 * swar.WordLen {
		a := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		b := haystack[i+offset : i+offset+2*swar.WordLen : i+offset+2*swar.WordLen]
		c0 := swar.ZeroLanes(swar.Load8(a, 0)^bc1) & swar.ZeroLanes(swar.Load8(b, 0)^bc2)
		c1 := swar.ZeroLanes(swar.Load8(a, swar.WordLen)^bc1) & swar.ZeroLanes(swar.Load8(b, swar.WordLen)^bc2)
		if c0|c1 != 0 {
			if pos := verifyPairCandidates(haystack, byte1, byte2, offset, i, c0); pos >= 0 {
				return pos
			}
			if pos := verifyPairCandidates(haystack, byte1, byte2, offset, i+swar.WordLen, c1); pos >= 0 {
				return pos
			}
		}
	}
	for ; i <= last-swar.WordLen; i += swar.WordLen {
		a := haystack[i : i+swar.WordLen : i+swar.WordLen]
		b := haystack[i+offset : i+offset+swar.WordLen : i+offset+swar.WordLen]
		c := swar.ZeroLanes(swar.Load8(a, 0)^bc1) & swar.ZeroLanes(swar.Load8(b, 0)^bc2)
		if c != 0 {
			if pos := verifyPairCandidates(haystack, byte1, byte2, offset, i, c); pos >= 0 {
				return pos
			}
		}
	}
	if i == last {
		return -1
	}
	i = last - swar.WordLen
	a := haystack[i : i+swar.WordLen : i+swar.WordLen]
	b := haystack[i+offset : i+offset+swar.WordLen : i+offset+swar.WordLen]
	c := swar.ZeroLanes(swar.Load8(a, 0)^bc1) & swar.ZeroLanes(swar.Load8(b, 0)^bc2)
	if c != 0 {
		return verifyPairCandidates(haystack, byte1, byte2, offset, i, c)
	}
	return -1
}

func verifyPairCandidates(haystack []byte, byte1, byte2 byte, offset, base int, cand uint64) int {
	for ; cand != 0; cand &= cand - 1 {
		pos := base + swar.FirstLane(cand)
		if haystack[pos] == byte1 && haystack[pos+offset] == byte2 {
			return pos
		}
	}
	return -1
}
