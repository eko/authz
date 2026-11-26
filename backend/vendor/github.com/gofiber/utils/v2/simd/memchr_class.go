// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

import (
	"github.com/gofiber/utils/v2/swar"
)

// MemchrWord returns the index of the first \w byte ([A-Za-z0-9_]), or -1.
func MemchrWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrWordImpl(haystack)
}

// MemchrNotWord returns the index of the first non-word byte, or -1 if none.
func MemchrNotWord(haystack []byte) int {
	if len(haystack) == 0 {
		return -1
	}
	return memchrNotWordImpl(haystack)
}

// MemchrInTable returns the index of the first byte b of haystack with
// table[b] true, or -1 if none matches or table is nil.
func MemchrInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if table[b] {
			return i
		}
	}
	return -1
}

// MemchrNotInTable returns the index of the first byte b of haystack with
// table[b] false, or -1 if every byte is in the table or table is nil.
func MemchrNotInTable(haystack []byte, table *[256]bool) int {
	if table == nil {
		return -1
	}
	for i, b := range haystack {
		if !table[b] {
			return i
		}
	}
	return -1
}

func isWordChar(b byte) bool {
	return (b >= 'A' && b <= 'Z') ||
		(b >= 'a' && b <= 'z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

// wordMask flags the lanes of w holding word characters, exactly per lane;
// bytes >= 0x80 never match. Letters share one 'a'..'z' test on a case-folded
// copy; '_' alone folds to 0x7F besides DEL, which w<<2 drops via bit 5.
// Biases are <= 0x50 on lanes <= 0x7F: no carries.
func wordMask(w uint64) uint64 {
	b := w & swar.LowSeven
	bf := b | 0x20*swar.Ones
	letters := (bf + (0x80-'a')*swar.Ones) &^ (bf + (0x80-'z'-1)*swar.Ones)
	digits := (b + (0x80-'0')*swar.Ones) &^ (b + (0x80-'9'-1)*swar.Ones)
	underscore := (bf + swar.Ones) &^ (w << 2)
	return (letters | digits | underscore) &^ w & swar.HighBits
}

func memchrWordGeneric(haystack []byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if isWordChar(haystack[i]) {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		// Pinned window: see the note in ascii.go.
		w := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := wordMask(swar.Load8(w, 0))
		m1 := wordMask(swar.Load8(w, swar.WordLen))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		if m := wordMask(swar.Load8(haystack, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := wordMask(swar.Load8(haystack, n-swar.WordLen)); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}

// wordMask is exact per lane, so its complement within the marker bits is too.
func memchrNotWordGeneric(haystack []byte) int {
	n := len(haystack)
	if n < swar.WordLen {
		for i := range n {
			if !isWordChar(haystack[i]) {
				return i
			}
		}
		return -1
	}
	i := 0
	for ; i <= n-2*swar.WordLen; i += 2 * swar.WordLen {
		w := haystack[i : i+2*swar.WordLen : i+2*swar.WordLen]
		m0 := swar.HighBits &^ wordMask(swar.Load8(w, 0))
		m1 := swar.HighBits &^ wordMask(swar.Load8(w, swar.WordLen))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + swar.WordLen + swar.FirstLane(m1)
		}
	}
	for ; i <= n-swar.WordLen; i += swar.WordLen {
		if m := swar.HighBits &^ wordMask(swar.Load8(haystack, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if m := swar.HighBits &^ wordMask(swar.Load8(haystack, n-swar.WordLen)); m != 0 {
		return n - swar.WordLen + swar.FirstLane(m)
	}
	return -1
}
