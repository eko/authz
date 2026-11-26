// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"math/bits"

	"github.com/gofiber/utils/v2/swar"
)

// IsASCII reports whether every byte in data is ASCII (< 0x80). An empty
// slice is trivially ASCII.
func IsASCII(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	return isASCIIImpl(data)
}

// FirstNonASCII returns the index of the first byte >= 0x80 in data, or -1
// if data is all ASCII.
func FirstNonASCII(data []byte) int {
	if len(data) == 0 {
		return -1
	}
	return firstNonASCIIImpl(data)
}

// CountNonASCII returns the number of bytes >= 0x80 in data.
func CountNonASCII(data []byte) int {
	return countNonASCIIImpl(data)
}

// Pinned windows: loops reslice each word group with three indices so the
// loads inside use constant offsets, whose bounds checks fold away; at a
// variable index swar.Load8 checks against capacity, not the loop's length.

// isASCIIGeneric ORs words per group and tests once, since no index is needed.
func isASCIIGeneric(data []byte) bool {
	n := len(data)
	if n < swar.WordLen {
		for i := range n {
			if data[i] >= 0x80 {
				return false
			}
		}
		return true
	}
	if n < 2*swar.WordLen {
		return (swar.Load8(data, 0)|swar.Load8(data, n-swar.WordLen))&swar.HighBits == 0
	}
	i := 0
	for ; i <= n-4*swar.WordLen; i += 4 * swar.WordLen {
		w := data[i : i+4*swar.WordLen : i+4*swar.WordLen]
		acc := swar.Load8(w, 0) |
			swar.Load8(w, swar.WordLen) |
			swar.Load8(w, 2*swar.WordLen) |
			swar.Load8(w, 3*swar.WordLen)
		if acc&swar.HighBits != 0 {
			return false
		}
	}
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		w := data[i : i+2*swar.WordLen : i+2*swar.WordLen]
		if (swar.Load8(w, 0)|swar.Load8(w, swar.WordLen))&swar.HighBits != 0 {
			return false
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		if swar.Load8(data, i)&swar.HighBits != 0 {
			return false
		}
	}
	if i == n {
		return true
	}
	return swar.Load8(data, n-swar.WordLen)&swar.HighBits == 0
}

func firstNonASCIIGeneric(data []byte) int {
	n := len(data)
	if n < swar.WordLen {
		for i := range n {
			if data[i] >= 0x80 {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		w := data[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := swar.Load8(w, 0) & swar.HighBits
		m1 := swar.Load8(w, swar.WordLen) & swar.HighBits
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		if m := swar.Load8(data, i) & swar.HighBits; m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	// Overlap re-reads known-ASCII lanes, so any hit is in the new bytes.
	if m := swar.Load8(data, n-swar.WordLen) & swar.HighBits; m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

func countNonASCIIGeneric(data []byte) int {
	count := 0
	i := 0
	n := len(data)
	for ; i <= n-4*swar.WordLen; i += 4 * swar.WordLen {
		w := data[i : i+4*swar.WordLen : i+4*swar.WordLen]
		count += bits.OnesCount64(swar.Load8(w, 0)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, swar.WordLen)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, 2*swar.WordLen)&swar.HighBits) +
			bits.OnesCount64(swar.Load8(w, 3*swar.WordLen)&swar.HighBits)
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		count += bits.OnesCount64(swar.Load8(data, i) & swar.HighBits)
	}
	// No overlapping word: it would double-count. The tail is under 8 bytes.
	for ; i < n; i++ {
		if data[i] >= 0x80 {
			count++
		}
	}
	return count
}
